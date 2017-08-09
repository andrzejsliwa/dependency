package dependency_test

import (
	"testing"

	"sort"

	"github.com/andrzejsliwa/dependency"
	"github.com/deckarep/golang-set"
	. "github.com/onsi/gomega"
)

//
//         a
//        /|
//       / |
//      b  |
//       \ |
//        \|
//         c
//         |
//         |
//         d
//
func graph1() dependency.Graph {
	g := dependency.NewGraph()
	g.Depend("b", "a")
	g.Depend("c", "b")
	g.Depend("c", "a")
	g.Depend("d", "c")
	return g
}

//       one     five
//        |       |
//       two      |
//       / \      |
//      /   \     |
//     /     \   /
//  three    four
//    |      /
//   six    /
//    |    /
//    |   /
//    |  /
//   seven
//
func graph2() dependency.Graph {
	g := dependency.NewGraph()
	g.Depend("two", "one")
	g.Depend("three", "two")
	g.Depend("four", "two")
	g.Depend("four", "five")
	g.Depend("six", "three")
	g.Depend("seven", "six")
	g.Depend("seven", "four")
	return g
}

//                level0
//               / | |  \
//          -----  | |   -----
//         /       | |        \
//  level1a  level1b  level1c  level1d
//         \       | |        /
//          -----  | |   -----
//               \ | |  /
//                level2
//               / | |  \
//          -----  | |   -----
//         /       | |        \
//  level3a  level3b  level3c  level3d
//         \       | |        /
//          -----  | |   -----
//               \ | |  /
//                level4
//
// ... and so on in a repeating pattern like that, up to level26
func graph3() dependency.Graph {
	g := dependency.NewGraph()
	g.Depend("level1a", "level0")
	g.Depend("level1b", "level0")
	g.Depend("level1c", "level0")
	g.Depend("level1d", "level0")
	g.Depend("level2", "level1a")
	g.Depend("level2", "level1b")
	g.Depend("level2", "level1c")
	g.Depend("level2", "level1d")

	g.Depend("level3a", "level2")
	g.Depend("level3b", "level2")
	g.Depend("level3c", "level2")
	g.Depend("level3d", "level2")
	g.Depend("level4", "level3a")
	g.Depend("level4", "level3b")
	g.Depend("level4", "level3c")
	g.Depend("level4", "level3d")

	g.Depend("level5a", "level4")
	g.Depend("level5b", "level4")
	g.Depend("level5c", "level4")
	g.Depend("level5d", "level4")
	g.Depend("level6", "level5a")
	g.Depend("level6", "level5b")
	g.Depend("level6", "level5c")
	g.Depend("level6", "level5d")

	g.Depend("level7a", "level6")
	g.Depend("level7b", "level6")
	g.Depend("level7c", "level6")
	g.Depend("level7d", "level6")
	g.Depend("level8", "level7a")
	g.Depend("level8", "level7b")
	g.Depend("level8", "level7c")
	g.Depend("level8", "level7d")

	g.Depend("level9a", "level8")
	g.Depend("level9b", "level8")
	g.Depend("level9c", "level8")
	g.Depend("level9d", "level8")
	g.Depend("level10", "level9a")
	g.Depend("level10", "level9b")
	g.Depend("level10", "level9c")
	g.Depend("level10", "level9d")

	g.Depend("level11a", "level10")
	g.Depend("level11b", "level10")
	g.Depend("level11c", "level10")
	g.Depend("level11d", "level10")
	g.Depend("level12", "level11a")
	g.Depend("level12", "level11b")
	g.Depend("level12", "level11c")
	g.Depend("level12", "level11d")

	g.Depend("level13a", "level12")
	g.Depend("level13b", "level12")
	g.Depend("level13c", "level12")
	g.Depend("level13d", "level12")
	g.Depend("level14", "level13a")
	g.Depend("level14", "level13b")
	g.Depend("level14", "level13c")
	g.Depend("level14", "level13d")

	g.Depend("level15a", "level14")
	g.Depend("level15b", "level14")
	g.Depend("level15c", "level14")
	g.Depend("level15d", "level14")
	g.Depend("level16", "level15a")
	g.Depend("level16", "level15b")
	g.Depend("level16", "level15c")
	g.Depend("level16", "level15d")

	g.Depend("level17a", "level16")
	g.Depend("level17b", "level16")
	g.Depend("level17c", "level16")
	g.Depend("level17d", "level16")
	g.Depend("level18", "level17a")
	g.Depend("level18", "level17b")
	g.Depend("level18", "level17c")
	g.Depend("level18", "level17d")
	return g
}

func TestGraph_TransitiveDependencies(t *testing.T) {
	RegisterTestingT(t)

	var cases = []struct {
		in       dependency.Graph
		node     string
		expected mapset.Set
	}{
		{graph1(), "d", mapset.NewSet("a", "b", "c")},
		{graph2(), "seven", mapset.NewSet("four", "three", "two", "five", "one", "six")},
	}

	for _, tt := range cases {
		Expect(tt.in.TransitiveDependencies(tt.node)).
			To(Equal(tt.expected))

	}
}

func TestGraph_ImmediateDependencies(t *testing.T) {
	RegisterTestingT(t)

	expected := graph2().ImmediateDependencies("four")
	Expect(expected).To(Equal(mapset.NewSet("two", "five")))
}

func TestGraph_ImmediateDependents(t *testing.T) {
	RegisterTestingT(t)

	expected := graph2().ImmediateDependents("four")
	Expect(expected).To(Equal(mapset.NewSet("seven")))
}

func TestGraph_RemoveEdgeSucceed(t *testing.T) {
	RegisterTestingT(t)

	err := graph2().RemoveEdge("three", "four")
	Expect(err).To(BeNil())
}

func TestGraph_RemoveEdgeFailed(t *testing.T) {
	RegisterTestingT(t)

	err := graph2().RemoveEdge("unknown", "unknown2")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).ToNot(Equal(""))
}

func TestGraph_RemoveNodeSucceed(t *testing.T) {
	RegisterTestingT(t)
	graph := graph2()
	err := graph.RemoveNode("three")
	Expect(err).To(BeNil())
	Expect(graph.Nodes()).To(Equal(mapset.NewSet("three", "five", "one", "seven", "two", "four", "six")))
}

func TestGraph_RemoveNodeFailed(t *testing.T) {
	RegisterTestingT(t)

	err := graph2().RemoveNode("unknown")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).ToNot(Equal(""))
}

func TestGraph_RemoveAllSucceed(t *testing.T) {
	RegisterTestingT(t)
	graph := graph2()
	err := graph.RemoveAll("three")
	Expect(err).To(BeNil())
	Expect(graph.Nodes()).To(Equal(mapset.NewSet("five", "one", "seven", "two", "four", "six")))
}

func TestGraph_RemoveAllFailed(t *testing.T) {
	RegisterTestingT(t)

	err := graph2().RemoveAll("unknown")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).ToNot(Equal(""))
}

func TestGraph_TransitiveDependenciesDeep(t *testing.T) {
	RegisterTestingT(t)

	expected := mapset.NewSet("level0", "level1a", "level1b", "level1c", "level1d",
		"level2",
		"level3a", "level3b", "level3c", "level3d",
		"level4",
		"level5a", "level5b", "level5c", "level5d",
		"level6",
		"level7a", "level7b", "level7c", "level7d",
		"level8",
		"level9a", "level9b", "level9c", "level9d",
		"level10",
		"level11a", "level11b", "level11c", "level11d",
		"level12",
		"level13a", "level13b", "level13c", "level13d",
		"level14",
		"level15a", "level15b", "level15c", "level15d",
		"level16",
		"level17a", "level17b", "level17c", "level17d")
	actual := graph3().TransitiveDependencies("level18")
	Expect(actual).To(Equal(expected))
}

func TestGraph_TransitiveDependenciesSet(t *testing.T) {
	RegisterTestingT(t)

	var cases = []struct {
		in       dependency.Graph
		nodes    mapset.Set
		expected mapset.Set
	}{
		{graph2(), mapset.NewSet("six", "four"), mapset.NewSet("one", "two", "three", "five")},
		{graph2(), mapset.NewSet("two", "four"), mapset.NewSet("one", "five")},
		{graph2(), mapset.NewSet("three", "six"), mapset.NewSet("two", "one")},
	}

	for _, tt := range cases {
		Expect(tt.in.TransitiveDependenciesSet(tt.nodes)).
			To(Equal(tt.expected))
	}

}

func TestGraph_TransitiveDependents(t *testing.T) {
	RegisterTestingT(t)

	var cases = []struct {
		in       dependency.Graph
		node     string
		expected mapset.Set
	}{
		{graph2(), "three", mapset.NewSet("six", "seven")},
		{graph2(), "five", mapset.NewSet("four", "seven")},
	}

	for _, tt := range cases {
		Expect(tt.in.TransitiveDependents(tt.node)).
			To(Equal(tt.expected))
	}

}

func TestGraph_TransitiveDependentsSet(t *testing.T) {
	RegisterTestingT(t)

	var cases = []struct {
		in       dependency.Graph
		nodes    mapset.Set
		expected mapset.Set
	}{
		{graph2(), mapset.NewSet("four", "three"), mapset.NewSet("six", "seven")},
		{graph2(), mapset.NewSet("four", "six"), mapset.NewSet("seven")},
	}

	for _, tt := range cases {
		Expect(tt.in.TransitiveDependentsSet(tt.nodes)).
			To(Equal(tt.expected))
	}

}

func TestGraph_Nodes(t *testing.T) {
	RegisterTestingT(t)

	expected := graph2().Nodes()
	Expect(expected).To(Equal(mapset.NewSet("one", "two", "three", "four", "six", "seven", "five")))
}

func TestGraph_TopologicalSort(t *testing.T) {
	RegisterTestingT(t)

	actual := graph2().TopologicalSort()
	nodes := []string{"seven", "four", "five", "six", "three", "two", "one"}
	expected := make([]interface{}, len(nodes))
	for i, s := range nodes {
		expected[i] = s
	}
	Expect(actual).To(Equal(expected))
}

func TestGraph_TopologicalComparator(t *testing.T) {
	RegisterTestingT(t)
	n := []string{"two", "five", "three"}
	c := graph2().TopologicalComparator(toInterfaceSlice(n))
	sort.Sort(c)
	Expect(c.Values()).To(Equal(toInterfaceSlice([]string{"five", "three", "two"})))
}

func toInterfaceSlice(input []string) []interface{} {
	nodes := make([]interface{}, len(input))
	for i, s := range input {
		nodes[i] = s
	}
	return nodes
}
