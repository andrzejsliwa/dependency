package dependency

import (
	"fmt"

	"github.com/deckarep/golang-set"
)

type Node interface{}

type Graph interface {
	GraphUpdate
	// Returns the set of immediate dependencies of node.
	ImmediateDependencies(node Node) mapset.Set
	// Returns the set of immediate dependents of node.
	ImmediateDependents(node Node) mapset.Set
	// Returns the set of all things which node depends on, directly or transitively.
	TransitiveDependencies(node Node) mapset.Set
	// Returns the set of all things which any node in node-set depends on, directly or transitively.
	TransitiveDependenciesSet(nodeSet mapset.Set) mapset.Set
	// Returns the set of all things which depend upon node, directly or transitively.
	TransitiveDependents(node Node) mapset.Set
	// Returns the set of all things which depend upon any node in node-set, directly or transitively.
	TransitiveDependentsSet(nodeSet mapset.Set) mapset.Set
	// Returns the set of all nodes in graph.
	Nodes() mapset.Set
	// Returns all nodes in topological order
	TopologicalSort() []interface{}
	// Returns topological comparator based on graph and
	TopologicalComparator(nodes []interface{}) *comparator
}

type GraphUpdate interface {
	// Adds dependency from node to dep ("node depends on dep"). Forbids circular dependencies.
	Depend(node Node, dep Node) error
	// Removes dependency from node to dep removed.
	RemoveEdge(node Node, dep Node) error
	// Removes dependency graph with all references to node removed.
	RemoveAll(node Node) error
	// Removes the node from the dependency graph without removing it as a dependency of other nodes. That is, removes all outgoing edges from node.
	RemoveNode(node Node) error
}

type graph struct {
	dependencies map[Node]mapset.Set
	dependents   map[Node]mapset.Set
}

func NewGraph() Graph {
	return &graph{make(map[Node]mapset.Set), make(map[Node]mapset.Set)}
}
func (g *graph) ImmediateDependencies(node Node) mapset.Set {
	return getOrDefault(g.dependencies, node)
}
func (g *graph) ImmediateDependents(node Node) mapset.Set {
	return getOrDefault(g.dependents, node)
}
func (g *graph) TransitiveDependencies(node Node) mapset.Set {
	return bfs(g.dependencies, mapset.NewSetWith(node))
}
func (g *graph) TransitiveDependenciesSet(nodeSet mapset.Set) mapset.Set {
	return bfs(g.dependencies, nodeSet)
}
func (g *graph) TransitiveDependents(node Node) mapset.Set {
	return bfs(g.dependents, mapset.NewSetWith(node))
}
func (g *graph) TransitiveDependentsSet(nodeSet mapset.Set) mapset.Set {
	return bfs(g.dependents, nodeSet)
}
func (g *graph) Nodes() mapset.Set {
	return keySet(g.dependencies).Union(keySet(g.dependents))
}
func (g *graph) RemoveEdge(node Node, dep Node) error {
	if _, ok := g.dependencies[node]; ok {
		delete(g.dependencies, node)
	} else {
		return fmt.Errorf("Edge node: %v, dep: %v not exist", node, dep)
	}

	if _, ok := g.dependents[dep]; ok {
		delete(g.dependents, dep)
	} else {
		return fmt.Errorf("Edge node: %v, dep: %v not exist", node, dep)
	}
	return nil
}
func (g *graph) RemoveAll(node Node) error {
	if !g.Nodes().Contains(node) {
		return fmt.Errorf("Unknown node: %v", node)
	}
	for k, v := range g.dependencies {
		if k == node {
			delete(g.dependencies, node)
		}
		if v.Contains(node) {
			v.Remove(node)
		}
	}
	for k, v := range g.dependents {
		if k == node {
			delete(g.dependents, node)
		}
		if v.Contains(node) {
			v.Remove(node)
		}
	}
	return nil
}
func (g *graph) RemoveNode(node Node) error {
	if _, ok := g.dependencies[node]; ok {
		delete(g.dependencies, node)
	} else {
		return fmt.Errorf("Edge node: %v not exist", node)
	}
	return nil
}
func (g *graph) Depend(node Node, dep Node) error {
	if node == dep || g.Depends(dep, node) {
		panic(fmt.Sprintf("Circular dependency: dependency '%v' already depends from '%v' via: %v", node, dep, g.TransitiveDependencies(dep)))
	}
	if _, ok := g.dependencies[node]; !ok {
		g.dependencies[node] = mapset.NewSet()
	}
	g.dependencies[node].Add(dep)

	if _, ok := g.dependents[dep]; !ok {
		g.dependents[dep] = mapset.NewSet()
	}
	g.dependents[dep].Add(node)
	return nil
}
func (g *graph) Depends(x, y Node) bool {
	return g.TransitiveDependencies(x).Contains(y)
}
func keySet(m map[Node]mapset.Set) mapset.Set {
	keys := mapset.NewSet()
	for k := range m {
		keys.Add(k)
	}
	return keys
}
func bfs(neighbors map[Node]mapset.Set, nodeSet mapset.Set) mapset.Set {
	frontier := nodeSet.ToSlice()
	visited := mapset.NewSet()
	next := mapset.NewSet().ToSlice()
	for 0 < len(frontier) {
		next = mapset.NewSet().ToSlice()
		for _, node := range frontier {
			visited.Add(node)
			for _, n := range bfs_frontier(node, neighbors, visited) {
				next = append(next, n)
			}
		}
		frontier = next
	}
	return visited.Difference(nodeSet)
}
func bfs_frontier(node Node, nodes map[Node]mapset.Set, visited mapset.Set) []interface{} {
	next := mapset.NewSet().ToSlice()
	iterator := func(n interface{}) bool { return !visited.Contains(n) }
	if nodes[node] != nil {
		for _, n := range nodes[node].ToSlice() {
			if iterator(n) {
				next = append(next, n)
			}
		}
	}
	return next
}
func getOrDefault(m map[Node]mapset.Set, node Node) mapset.Set {
	if value, ok := m[node]; ok {
		return value
	} else {
		return mapset.NewSet()
	}
}

func (g *graph) TopologicalSort() []interface{} {
	sorted := make([]interface{}, 0)
	inDegree := map[interface{}]int{}

	// 1. Calculate inDegree of all vertices by going through every edge of the graph.
	// Each child gets inDegree++ during breadth-first run.
	for element, children := range g.dependencies {
		if inDegree[element] == 0 {
			inDegree[element] = 0
		}
		for _, child := range children.ToSlice() {
			inDegree[child]++
		}
	}
	// 2. Collect all vertices with inDegree == 0 onto a stack.
	stack := make([]interface{}, 0)
	for rule, value := range inDegree {
		if value == 0 {
			stack = append(stack, rule)
			inDegree[rule] = -1
		}
	}

	// 3. While zero-degree-stack is not empty.
	for len(stack) > 0 {
		var node interface{}
		// 3.1. Pop element from zero-degree-stack and append it to topological order.
		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// 3.2. Find all children of element and decrease inDegree.
		// If inDegree becomes 0, add to zero-degree-stack.
		if g.dependencies[node] != nil {
			for _, child := range g.dependencies[node].ToSlice() {
				inDegree[child]--
				if inDegree[child] == 0 {
					stack = append(stack, child)
					inDegree[child] = -1
				}
			}
		}
		// 3.3. Append to the sorted list.
		sorted = append(sorted, node)
	}
	return sorted
}

type comparator struct {
	pos map[interface{}]int
}

func (s comparator) Len() int {
	return len(s.pos)
}
func (s comparator) Swap(i, j int) {
	s.pos[i], s.pos[j] = s.pos[j], s.pos[i]
}
func (s comparator) Less(i, j int) bool {
	return s.pos[i] < s.pos[j]
}
func (s comparator) Values() []interface{} {
	keys := make([]interface{}, 0)
	for k := range s.pos {
		keys = append(keys, k)
	}
	return keys
}
func (g *graph) TopologicalComparator(nodes []interface{}) *comparator {
	nodes2 := mapset.NewSetFromSlice(nodes)
	pos := map[interface{}]int{}
	for order, element := range g.TopologicalSort() {
		if nodes2.Contains(element) {
			pos[element] = order
		}
	}
	return &comparator{pos}
}
