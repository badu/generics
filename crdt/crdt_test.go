package crdt

import (
	"testing"
)

const (
	vertex1Name = "one"
	vertex2Name = "two"
	vertex3Name = "three"
	vertex4Name = "four"
)

func setup() ILastWriterWinsGraph[string] {
	g := NewLastWriterWinsGraph[string]()

	g.AddVertex(vertex1Name)
	g.AddVertex(vertex2Name)
	g.AddVertex(vertex3Name)

	g.AddEdge(vertex1Name, vertex2Name)
	g.AddEdge(vertex2Name, vertex3Name)

	return g
}

func samePaths[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}

	return true
}

func has[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func sameSets[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for _, v := range s1 {
		if !has(s2, v) {
			return false
		}
	}

	return true
}

func TestNewVertex(t *testing.T) {
	g := setup()

	err := g.AddVertex(vertex4Name)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	vertices, err := g.Vertices()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !g.HasVertex(vertex4Name) {
		t.Errorf("vertex not found, got: %v, expected: %v.", vertices, append(vertices, vertex4Name))
	}
}

func TestListVertices(t *testing.T) {
	g := setup()

	vertices, err := g.Vertices()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := []string{vertex1Name, vertex2Name, vertex3Name}
	if !sameSets(vertices, expected) {
		t.Errorf("Vertices mismatch, got: %v, expected: %v.", vertices, expected)
	}
}

func TestDeleteVertex(t *testing.T) {
	g := setup()

	err := g.DeleteVertex(vertex3Name)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	vertices, err := g.Vertices()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if g.HasVertex(vertex3Name) {
		t.Errorf("vertex found, got: %v, expected: %v.", vertices, []string{"vertex1", vertex2Name})
	}
}

func TestHasVertices(t *testing.T) {
	g := setup()
	tests := []struct {
		vertex   string
		expected bool
	}{
		{vertex1Name, true},
		{"inexistent_vertex", false},
	}
	for _, tt := range tests {
		got := g.HasVertex(tt.vertex)
		if got != tt.expected {
			t.Errorf("Existence check failed, got: %v, expected: %v.", got, tt.expected)
		}
	}
}

func TestAddEdge(t *testing.T) {
	g := setup()
	err := g.AddEdge(vertex1Name, vertex3Name)
	edges1, _ := g.VertexEdges(vertex1Name)
	edges3, _ := g.VertexEdges(vertex3Name)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !g.HasEdge(vertex1Name, vertex3Name) {
		t.Errorf("missing edge, got:[ %v, %v ], expected: [ %v, %v ].", edges1, edges3, []string{vertex2Name, vertex3Name}, []string{"vertex1", vertex2Name})
	}
}

func TestRemoveEdge(t *testing.T) {
	g := setup()

	err := g.DeleteEdge(vertex1Name, vertex2Name)
	edges1, _ := g.VertexEdges(vertex1Name)
	edges2, _ := g.VertexEdges(vertex2Name)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if g.HasEdge(vertex1Name, vertex2Name) {
		t.Errorf("edge found, got:[ %v, %v ], expected: [ %v, %v ].", edges1, edges2, []string{}, []string{vertex3Name})
	}
}

func TestHasEdge(t *testing.T) {
	g := setup()

	tests := []struct {
		vertices []string
		expected bool
	}{
		{[]string{vertex1Name, vertex2Name}, true},
		{[]string{vertex1Name, "inexistent_vertex"}, false},
	}
	for _, tt := range tests {
		got := g.HasEdge(tt.vertices[0], tt.vertices[1])
		if got != tt.expected {
			t.Errorf("check failed, got: %v, expected: %v.", got, tt.expected)
		}
	}
}

func TestHasEdges(t *testing.T) {
	g := setup()

	edges, err := g.VertexEdges(vertex1Name)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if edges[0] != vertex2Name {
		t.Errorf("edge found, got: %v, expected: %v.", edges, []string{vertex2Name})
	}
}

func TestFindPath(t *testing.T) {
	simpleG := setup()
	cycleG := setup()
	cycleG.AddVertex(vertex4Name)
	cycleG.AddEdge(vertex3Name, vertex4Name)
	disconnectedG := setup()
	disconnectedG.DeleteEdge(vertex1Name, vertex2Name)
	tests := []struct {
		name     string
		graph    ILastWriterWinsGraph[string]
		from     string
		target   string
		expected []string
	}{
		{"simple graph", simpleG, vertex1Name, vertex3Name, []string{vertex1Name, vertex2Name, vertex3Name}},
		{"cycle graph", cycleG, vertex1Name, vertex4Name, []string{vertex1Name, vertex2Name, vertex3Name, vertex4Name}},
		{"disconnected graph", disconnectedG, vertex1Name, vertex3Name, []string{vertex1Name}},
	}
	for _, tt := range tests {
		got, err := tt.graph.FindPath(tt.from, tt.target)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !samePaths(got, tt.expected) {
			t.Errorf("Finding path for %s failed, got: %v, expected: %v.", tt.name, got, tt.expected)
		}
	}
}

func TestMerge(t *testing.T) {
	g1 := setup()
	g2 := setup()

	g2.DeleteVertex(vertex2Name)
	g1.AddVertex(vertex2Name)
	g2.AddEdge(vertex1Name, vertex3Name)
	g2.AddVertex(vertex4Name)
	g1.DeleteVertex(vertex4Name)

	err := g1.Merge(g2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !g1.HasVertex(vertex2Name) || g1.HasVertex(vertex4Name) {
		t.Errorf("Vertices merge failed")
	}

	if !g1.HasEdge(vertex1Name, vertex3Name) {
		t.Errorf("merge failed")
	}
}

func TestAssociativity(t *testing.T) {
	g1 := setup()
	g2 := setup()
	g3 := setup()
	g1.AddVertex(vertex4Name)
	g3.DeleteVertex(vertex4Name)
	g2.Merge(g3)
	g1.Merge(g2)

	g4 := setup()
	g5 := setup()
	g6 := setup()
	g4.AddVertex(vertex4Name)
	g6.DeleteVertex(vertex4Name)
	g4.Merge(g5)
	g4.Merge(g6)

	v1, _ := g1.Vertices()
	v4, _ := g4.Vertices()
	if !sameSets(v1, v4) {
		t.Errorf("not associative, g1 v (g2 v g3): %v, (g1 v g2) v g3: %v.", v1, v4)
	}
}

func TestCommutativity(t *testing.T) {
	g1 := setup()
	g2 := setup()
	g1.AddVertex(vertex4Name)
	g2.DeleteVertex(vertex4Name)
	g1.Merge(g2)

	g4 := setup()
	g5 := setup()
	g4.AddVertex(vertex4Name)
	g5.DeleteVertex(vertex4Name)
	g5.Merge(g4)

	v1, _ := g1.Vertices()
	v5, _ := g5.Vertices()
	if !sameSets(v1, v5) {
		t.Errorf("not associative, g1 v g2: %v, g2 v g1: %v.", v1, v5)
	}
}

func TestIdempotence(t *testing.T) {
	g := setup()
	g2 := setup()

	g.AddVertex(vertex4Name)
	g.DeleteVertex(vertex3Name)

	g2.AddVertex(vertex4Name)
	g2.DeleteVertex(vertex3Name)

	beforeVertices, _ := g2.Vertices()
	g.Merge(g2)
	afterVertices, _ := g.Vertices()
	if !sameSets(beforeVertices, afterVertices) {
		t.Errorf("not idempotent, g1: %v, g1 v g1: %v.", beforeVertices, afterVertices)
	}
}
