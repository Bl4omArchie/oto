package oto

import (
	"fmt"
	"errors"

	"github.com/katalvlaran/lvlath/bfs"
	"github.com/katalvlaran/lvlath/core"
	"github.com/katalvlaran/lvlath/dfs"

	"github.com/Bl4omArchie/oto/models"
)

// Schema holds the static constraint model:
//
//	 - a map of the Parameters accepted in this set with their flagID as a key
//   - a directed dependency graph (requires edges) over ArgID values;
//   - a symmetric conflict relation stored as a map of sets.
//
// The graph is used for:
//   - computing the transitive closure of dependencies via BFS;
//   - checking for cycles and building execution order via DFS.
//
// The conflict map is kept separate for simplicity and cheaper lookups.
type Schema struct {
	flags		[]models.FlagID
	g			*core.Graph
	conflicts	map[models.FlagID]map[models.FlagID]struct{}
}

// CombinationOutput holds :
//
//	- a set of flag IDs (which represent the flag field of a Parameter)
//	- a boolean that indicates if either the combination is valid or not
//
// This struct is used by ValidateSchema() function as a return value
type CombinationOutput struct {
	Flags []models.FlagID
	Valid bool
}

// Sentinel errors for explicit semantics.
var (
	ErrSchemaCycle         = errors.New("flagschema: dependency cycle")
	ErrSchemaContradiction = errors.New("flagschema: schema contradiction between dependency and conflict")
	ErrSelectionConflict   = errors.New("flagschema: conflicting arguments in selection")
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
// lvlath/core and an empty conflict relation.
//
// Complexity of operations on this structure is dominated by BFS/DFS:
//   - Schema validation: O(V + E + C * (V + E)) in the worst case;
//   - Selection validation: O(|S| * (V + E) + C) for practical sizes.
func NewSchema(flags []models.FlagID) *Schema {
	return &Schema{
		flags:		flags,
		g:			core.NewGraph(core.WithDirected(true)),
		conflicts:	make(map[models.FlagID]map[models.FlagID]struct{}),
	}
}

// ensureArg makes sure there is a vertex for id in the underlying graph.
// It is intentionally idempotent: calling it multiple times is safe.
func (s *Schema) ensureArg(id models.FlagID) {
	_ = s.g.AddVertex(string(id))
}

// Require declares a dependency A → B meaning “A requires B”.
// In the graph this is represented as a directed edge A -> B.
func (s *Schema) Require(a, b models.FlagID) {
	s.ensureArg(a)
	s.ensureArg(b)
	// Unweighted dependency edge: weight = 0.
	_, _ = s.g.AddEdge(string(a), string(b), 0)
}

// Conflict registers a symmetric conflict between A and B.
// If A and B end up together in the closure of a selection, that selection
// is considered invalid.
func (s *Schema) Conflict(a, b models.FlagID) {
	s.ensureArg(a)
	s.ensureArg(b)

	if a == b {
		// A self-conflict does not make sense; ignore defensively.
		return
	}

	if s.conflicts[a] == nil {
		s.conflicts[a] = make(map[models.FlagID]struct{})
	}
	if s.conflicts[b] == nil {
		s.conflicts[b] = make(map[models.FlagID]struct{})
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
				Detail: "invalid schema: dependency cycle detected in argument graph",
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

// ExecutionOrder computes a deterministic execution order for the subset
// of arguments given in flags, respecting all dependency edges.
//
// Internally it builds an induced subgraph on the subset and runs a
// topological sort via dfs.TopologicalSort.
func (s *Schema) ExecutionOrder(flags []models.FlagID) ([]models.FlagID, error) {
	if len(flags) == 0 {
		return nil, nil
	}

	needed := make(map[string]struct{}, len(flags))
	for _, id := range flags {
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

	out := make([]models.FlagID, len(order))
	for i, v := range order {
		out[i] = models.FlagID(v)
	}
	return out, nil
}

// reachable runs BFS once and checks whether "to" is reachable from "from"
// using the Depth map from BFSResult. This is enough for schema-level checks
// where only reachability matters, not the exact path.
func reachable(g *core.Graph, from, to models.FlagID) bool {
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
func shortestPath(g *core.Graph, from, to models.FlagID) []models.FlagID {
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

	path := make([]models.FlagID, len(strPath))
	for i, v := range strPath {
		path[i] = models.FlagID(v)
	}
	return path
}
