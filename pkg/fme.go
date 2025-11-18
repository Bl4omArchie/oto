// fme.go
//
// High-level pseudocode (how lvlath is used to solve OTO-style constraints):
//
//  1. Build a directed dependency graph G over Flagument IDs using core.Graph.
//     For each rule "A requires B", add a directed edge A -> B with weight 0.
//  2. Maintain a separate symmetric conflict relation C ⊆ V×V in a Go map.
//  3. At service startup:
//     a) Run dfs.TopologicalSort(G) to ensure that dependencies form a DAG.
//        If a cycle is found, reject the schema (misconfigured rules).
//     b) For each conflict {a, b} in C, run bfs.BFS(G, a) and bfs.BFS(G, b).
//        If a and b are mutually reachable via "requires" edges, reject
//        the schema as contradictory.
//  4. For each user selection S:
//     a) Compute closure(S) by running BFS from every a ∈ S over G,
//        adding all reachable Flaguments to the set "need".
//     b) If "need" contains any conflicting pair {a, b} ∈ C, reject the
//        selection, optionally explaining the shortest implication chain.
//     c) Otherwise, build an induced subgraph on "need" and run
//        dfs.TopologicalSort to get a safe execution order for tasks.
//
// This example is a skeleton for an "Flagument-constraints engine" that can
// sit underneath a scheduler like OTO while staying small, explicit, and
// easy to explain to users.

package oto

import (
	"errors"
	"fmt"
	"sort"

	"github.com/katalvlaran/lvlath/bfs"
	"github.com/katalvlaran/lvlath/core"
	"github.com/katalvlaran/lvlath/dfs"
)

// FlagID is a domain-level identifier for an Flagument/flag.
type FlagID string

// Schema holds the static constraint model:
//
//   - a directed dependency graph (requires edges) over FlagID values;
//   - a symmetric conflict relation stored as a map of sets.
//
// The graph is used for:
//   - computing the transitive closure of dependencies via BFS;
//   - checking for cycles and building execution order via DFS.
//
// The conflict map is kept separate for simplicity and cheaper lookups.
type Schema struct {
	g         *core.Graph
	conflicts map[FlagID]map[FlagID]struct{}
}

// Sentinel errors for explicit semantics.
var (
	ErrSchemaCycle         = errors.New("Schema: dependency cycle")
	ErrSchemaContradiction = errors.New("Schema: schema contradiction between dependency and conflict")
	ErrSelectionConflict   = errors.New("Schema: conflicting Flaguments in selection")
)

// SchemaValidationError wraps a schema-level validation failure with:
//   - a well-known Kind (one of ErrSchemaCycle / ErrSchemaContradiction);
//   - a human-readable Detail string that can go to logs or user-facing errors.
type SchemaValidationError struct {
	Kind   error  // ErrSchemaCycle or ErrSchemaContradiction
	Detail string // human-readable explanation
}

func (e *SchemaValidationError) Error() string { return e.Detail }
func (e *SchemaValidationError) Unwrap() error { return e.Kind }

// ConflictInstance describes a concrete conflict discovered in a selection:
//   - A and B are the conflicting Flaguments;
//   - PathAB (if non-nil) is the shortest implication chain A ⇒ ... ⇒ B
//     recovered from BFS parent information. It is useful for "why" messages.
type ConflictInstance struct {
	A, B   FlagID
	PathAB []FlagID
}

// SelectionResult is the outcome of validating a concrete user selection.
//
//   - Final   is the closure(selection): selected Flaguments + all dependencies;
//   - Conflict is non-nil if a conflicting pair was detected.
type SelectionResult struct {
	Final    []FlagID
	Conflict *ConflictInstance
}

// NewSchema constructs a directed, unweighted dependency graph backed by
// lvlath/core and an empty conflict relation.
//
// Complexity of operations on this structure is dominated by BFS/DFS:
//   - Schema validation: O(V + E + C * (V + E)) in the worst case;
//   - Selection validation: O(|S| * (V + E) + C) for practical sizes.
func NewSchema() *Schema {
	return &Schema{
		g:         core.NewGraph(core.WithDirected(true)),
		conflicts: make(map[FlagID]map[FlagID]struct{}),
	}
}

// ensureFlag makes sure there is a vertex for id in the underlying graph.
// It is intentionally idempotent: calling it multiple times is safe.
func (s *Schema) ensureFlag(id FlagID) {
	_ = s.g.AddVertex(string(id))
}

// Require declares a dependency A → B meaning “A requires B”.
// In the graph this is represented as a directed edge A -> B.
func (s *Schema) Require(a, b FlagID) {
	s.ensureFlag(a)
	s.ensureFlag(b)
	// Unweighted dependency edge: weight = 0.
	_, _ = s.g.AddEdge(string(a), string(b), 0)
}

// Conflict registers a symmetric conflict between A and B.
// If A and B end up together in the closure of a selection, that selection
// is considered invalid.
func (s *Schema) Conflict(a, b FlagID) {
	s.ensureFlag(a)
	s.ensureFlag(b)

	if a == b {
		// A self-conflict does not make sense; ignore defensively.
		return
	}

	if s.conflicts[a] == nil {
		s.conflicts[a] = make(map[FlagID]struct{})
	}
	if s.conflicts[b] == nil {
		s.conflicts[b] = make(map[FlagID]struct{})
	}
	s.conflicts[a][b] = struct{}{}
	s.conflicts[b][a] = struct{}{}
}

// ValidateSchema performs static validation of the constraint schema.
//
//  1. Ensures that the dependency graph is a DAG via dfs.TopologicalSort.
//  2. Ensures that no conflict pair {A, B} is such that one is reachable
//     from the other through “requires” edges (which would be a contradiction).
//
// This function is intended to be called once at service startup.
func (s *Schema) ValidateSchema() error {
	// 1. Check that the dependency graph is acyclic (DAG).
	if _, err := dfs.TopologicalSort(s.g); err != nil {
		if errors.Is(err, dfs.ErrCycleDetected) {
			return &SchemaValidationError{
				Kind:   ErrSchemaCycle,
				Detail: "invalid schema: dependency cycle detected in Flagument graph",
			}
		}
		// Any other error is unexpected and should be surfaced as-is.
		return err
	}

	// 2. Ensure that no conflicting pair is forced by dependencies.
	for a, row := range s.conflicts {
		for b := range row {
			// Work with each unordered pair only once (a < b).
			if a >= b {
				continue
			}

			if reachable(s.g, a, b) || reachable(s.g, b, a) {
				msg := fmt.Sprintf(
					"invalid schema: %q and %q are declared as conflicting, "+
						"but one is reachable from the other via requires edges",
					a, b,
				)
				return &SchemaValidationError{
					Kind:   ErrSchemaContradiction,
					Detail: msg,
				}
			}
		}
	}

	return nil
}

// ValidateSelection:
//
//   - expands the initial selection by adding all transitive dependencies;
//   - checks for conflicts in the resulting closure;
//   - returns SelectionResult and either nil or ErrSelectionConflict.
//
// This is the function you would typically call per user request.
func (s *Schema) ValidateSelection(selection []FlagID) (*SelectionResult, error) {
	need := make(map[FlagID]struct{})

	// Make sure all selected Flaguments exist as vertices.
	for _, id := range selection {
		s.ensureFlag(id)
	}

	// For each selected Flagument, run BFS to collect all dependencies.
	for _, id := range selection {
		res, err := bfs.BFS(s.g, string(id))
		if err != nil {
			return nil, fmt.Errorf("selection: BFS from %q: %w", id, err)
		}
		for _, ID := range res.Order {
			need[FlagID(ID)] = struct{}{}
		}
	}

	// Convert to a sorted slice for deterministic output and testing.
	final := make([]FlagID, 0, len(need))
	for id := range need {
		final = append(final, id)
	}
	sort.Slice(final, func(i, j int) bool { return final[i] < final[j] })

	// Now check conflicts inside the closure.
	for a, row := range s.conflicts {
		for b := range row {
			if a >= b {
				continue
			}
			_, hasA := need[a]
			_, hasB := need[b]
			if hasA && hasB {
				ci := &ConflictInstance{
					A:      a,
					B:      b,
					PathAB: shortestPath(s.g, a, b),
				}
				return &SelectionResult{
					Final:    final,
					Conflict: ci,
				}, ErrSelectionConflict
			}
		}
	}

	return &SelectionResult{Final: final}, nil
}

// ExecutionOrder computes a deterministic execution order for the subset
// of Flaguments given in Flags, respecting all dependency edges.
//
// Internally it builds an induced subgraph on the subset and runs a
// topological sort via dfs.TopologicalSort.
func (s *Schema) ExecutionOrder(Flags []FlagID) ([]FlagID, error) {
	if len(Flags) == 0 {
		return nil, nil
	}

	needed := make(map[string]struct{}, len(Flags))
	for _, id := range Flags {
		needed[string(id)] = struct{}{}
	}

	// Build an induced subgraph on the needed vertices.
	sub := s.g.CloneEmpty()
	for id := range needed {
		_ = sub.AddVertex(id)
	}

	for _, e := range s.g.Edges() {
		if _, ok := needed[e.From]; !ok {
			continue
		}
		if _, ok := needed[e.To]; !ok {
			continue
		}
		_, _ = sub.AddEdge(e.From, e.To, e.Weight, core.WithEdgeDirected(e.Directed))
	}

	order, err := dfs.TopologicalSort(sub)
	if err != nil {
		if errors.Is(err, dfs.ErrCycleDetected) {
			// Should not happen if ValidateSchema has already passed.
			return nil, &SchemaValidationError{
				Kind:   ErrSchemaCycle,
				Detail: "cycle detected in induced subgraph for selection",
			}
		}
		return nil, err
	}

	out := make([]FlagID, len(order))
	for i, v := range order {
		out[i] = FlagID(v)
	}
	return out, nil
}

// reachable runs BFS once and checks whether "to" is reachable from "from"
// using the Depth map from BFSResult. This is enough for schema-level checks
// where only reachability matters, not the exact path.
func reachable(g *core.Graph, from, to FlagID) bool {
	res, err := bfs.BFS(g, string(from))
	if err != nil {
		// In a well-formed schema this should not fail; treat as not reachable.
		return false
	}

	_, ok := res.Depth[string(to)]
	return ok
}

// shortestPath reconstructs a shortest path (in edges) between from and to
// in the unweighted dependency graph, using the PathTo helper from BFSResult.
// If there is no path, it returns nil.
func shortestPath(g *core.Graph, from, to FlagID) []FlagID {
	res, err := bfs.BFS(g, string(from))
	if err != nil {
		// For explanation purposes, failing silently is acceptable:
		// we simply skip attaching a "reason" path to the conflict.
		return nil
	}

	// Use the BFSResult helper, which already relies on Depth/Parent
	// and returns an error when dest has not been reached.
	strPath, err := res.PathTo(string(to))
	if err != nil {
		// No path from "from" to "to".
		return nil
	}

	path := make([]FlagID, len(strPath))
	for i, v := range strPath {
		path[i] = FlagID(v)
	}
	return path
}
