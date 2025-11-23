// fme.go
//
// Original script made by Kyrylo (@katalvlaran on github)
//
// Adaptation by Archie (@Bl4omArchie)

package oto

import (
	"fmt"
	"sort"
	"errors"

	"github.com/katalvlaran/lvlath/bfs"
	"github.com/katalvlaran/lvlath/core"
	"github.com/katalvlaran/lvlath/dfs"
)

// Schema holds the static constraint model:
//
//   - a directed dependency graph (requires edges) over string values;
//   - a symmetric interfer relation stored as a map of sets.
//
// The graph is used for:
//   - computing the transitive closure of dependencies via BFS;
//   - checking for cycles and building execution order via DFS.
//
// The interfer map is kept separate for simplicity and cheaper lookups.
type Schema struct {
	g         *core.Graph
	interfers map[string]map[string]struct{}
}

// InterferInstance describes a concrete interfer discovered in a combination:
//   - A and B are the interfering Flags;
//   - PathAB (if non-nil) is the shortest implication chain A ⇒ ... ⇒ B
//     recovered from BFS parent information. It is useful for "why" messages.
type InterferInstance struct {
	A, B   string
	PathAB []string
}

// CombinationResult is the outcome of validating a concrete user combination.
//
//   - Final   is the closure(combination): selected Flags + all dependencies;
//   - Interfer is non-nil if a interfering pair was detected.
type CombinationResult struct {
	Final    []string
	Interfer *InterferInstance
}

// Sentinel errors for explicit semantics.
var (
	ErrSchemaCycle         = errors.New("argschema: dependency cycle")
	ErrSchemaContradiction = errors.New("argschema: schema contradiction between dependency and interference")
	ErrCombinationInterfer   = errors.New("argschema: interfering arguments in combination")
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


// NewSchema constructs a directed, unweighted dependency graph backed by
// lvlath/core and an empty interfer relation.
//
// Complexity of operations on this structure is dominated by BFS/DFS:
//   - Schema validation: O(V + E + C * (V + E)) in the worst case;
//   - Combination validation: O(|S| * (V + E) + C) for practical sizes.
func NewSchema() *Schema {
	return &Schema{
		g:         core.NewGraph(core.WithDirected(true)),
		interfers: make(map[string]map[string]struct{}),
	}
}

// ensureFlag makes sure there is a vertex for id in the underlying graph.
// It is intentionally idempotent: calling it multiple times is safe.
func (s *Schema) ensureFlag(id string) {
	_ = s.g.AddVertex(string(id))
}

// Require declares a dependency A → B meaning “A requires B”.
// In the graph this is represented as a directed edge A -> B.
func (s *Schema) Require(a, b string) {
	s.ensureFlag(a)
	s.ensureFlag(b)
	// Unweighted dependency edge: weight = 0.
	_, _ = s.g.AddEdge(string(a), string(b), 0)
}

// Interfer registers a symmetric interfer between A and B.
// If A and B end up together in the closure of a combination, that combination
// is considered invalid.
func (s *Schema) Interfer(a, b string) {
	s.ensureFlag(a)
	s.ensureFlag(b)

	if a == b {
		// A self-interfer does not make sense; ignore defensively.
		return
	}

	if s.interfers[a] == nil {
		s.interfers[a] = make(map[string]struct{})
	}
	if s.interfers[b] == nil {
		s.interfers[b] = make(map[string]struct{})
	}
	s.interfers[a][b] = struct{}{}
	s.interfers[b][a] = struct{}{}
}

// ValidateSchema performs static validation of the constraint schema.
//
//  1. Ensures that the dependency graph is a DAG via dfs.TopologicalSort.
//  2. Ensures that no interfer pair {A, B} is such that one is reachable
//     from the other through “requires” edges (which would be a contradiction).
//
// This function is intended to be called once at service startup.
func (s *Schema) ValidateSchema() error {
	// 1. Check that the dependency graph is acyclic (DAG).
	if _, err := dfs.TopologicalSort(s.g); err != nil {
		if errors.Is(err, dfs.ErrCycleDetected) {
			return &SchemaValidationError{
				Kind:   ErrSchemaCycle,
				Detail: "invalid schema: dependency cycle detected in Flag graph",
			}
		}
		// Any other error is unexpected and should be surfaced as-is.
		return err
	}

	// 2. Ensure that no interfering pair is forced by dependencies.
	for a, row := range s.interfers {
		for b := range row {
			// Work with each unordered pair only once (a < b).
			if a >= b {
				continue
			}

			if reachable(s.g, a, b) || reachable(s.g, b, a) {
				msg := fmt.Sprintf(
					"invalid schema: %q and %q are declared as interfering, "+
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

// ValidateCombination:
//
//   - expands the initial combination by adding all transitive dependencies;
//   - checks for interfers in the resulting closure;
//   - returns CombinationResult and either nil or ErrCombinationInterfer.
//
// This is the function you would typically call per user request.
func (s *Schema) ValidateCombination(combination []string) (*CombinationResult, error) {
	need := make(map[string]struct{})

	// Make sure all selected Flags exist as vertices.
	for _, id := range combination {
		s.ensureFlag(id)
	}

	// For each selected Flag, run BFS to collect all dependencies.
	for _, id := range combination {
		res, err := bfs.BFS(s.g, string(id))
		if err != nil {
			return nil, fmt.Errorf("combination: BFS from %q: %w", id, err)
		}
		for _, ID := range res.Order {
			need[string(ID)] = struct{}{}
		}
	}

	// Convert to a sorted slice for deterministic output and testing.
	final := make([]string, 0, len(need))
	for id := range need {
		final = append(final, id)
	}
	sort.Slice(final, func(i, j int) bool { return final[i] < final[j] })

	// Now check interfers inside the closure.
	for a, row := range s.interfers {
		for b := range row {
			if a >= b {
				continue
			}
			_, hasA := need[a]
			_, hasB := need[b]
			if hasA && hasB {
				ci := &InterferInstance{
					A:      a,
					B:      b,
					PathAB: shortestPath(s.g, a, b),
				}
				return &CombinationResult{
					Final:    final,
					Interfer: ci,
				}, ErrCombinationInterfer
			}
		}
	}

	return &CombinationResult{Final: final}, nil
}

// ExecutionOrder computes a deterministic execution order for the subset
// of Flags given in Flags, respecting all dependency edges.
//
// Internally it builds an induced subgraph on the subset and runs a
// topological sort via dfs.TopologicalSort.
func (s *Schema) ExecutionOrder(Flags []string) ([]string, error) {
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
				Detail: "cycle detected in induced subgraph for combination",
			}
		}
		return nil, err
	}

	out := make([]string, len(order))
	for i, v := range order {
		out[i] = string(v)
	}
	return out, nil
}

// reachable runs BFS once and checks whether "to" is reachable from "from"
// using the Depth map from BFSResult. This is enough for schema-level checks
// where only reachability matters, not the exact path.
func reachable(g *core.Graph, from, to string) bool {
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
func shortestPath(g *core.Graph, from, to string) []string {
	res, err := bfs.BFS(g, string(from))
	if err != nil {
		// For explanation purposes, failing silently is acceptable:
		// we simply skip attaching a "reason" path to the interfer.
		return nil
	}

	// Use the BFSResult helper, which already relies on Depth/Parent
	// and returns an error when dest has not been reached.
	strPath, err := res.PathTo(string(to))
	if err != nil {
		// No path from "from" to "to".
		return nil
	}

	path := make([]string, len(strPath))
	for i, v := range strPath {
		path[i] = string(v)
	}
	return path
}
