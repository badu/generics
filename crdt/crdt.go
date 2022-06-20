package crdt

/**

From : https://github.com/bjornaer/CRDT

Conflict-Free Replicated Data Types (CRDTs) are data structures that power real-time collaborative applications in distributed systems.
Conflict-Free Replicated Data Types can be replicated across systems, they can be updated independently and concurrently without coordination between the replicas,
and it is always mathematically possible to resolve inconsistencies that might result.

In other (more practical) words: Conflict-Free Replicated Data Types are a certain form of data types that when replicated across several nodes over a network achieve eventual consistency without the need for a consensus round

Last Write Wins Element Set
---
Last Write Wins Element Set is similar to 2P-Set in that it consists of an "add set" and a "remove set", with a timestamp for each element.
Records are added to an Last Write Wins Element Set by inserting the element into the add set, with a timestamp.
Records are removed from the Last Write Wins Element Set by being added to the remove set, again with a timestamp.
An element is a member of the Last Write Wins Element Set if it is in the add set, and either not in the remove set, or in the remove set but with an earlier timestamp than the latest timestamp in the add set.
Merging two replicas of the Last Write Wins Element Set consists of taking the union of the add sets and the union of the remove sets.
When timestamps are equal, the "bias" of the Last Write Wins Element Set comes into play.
A Last Write Wins Element Set can be biased towards adds or deletions.
An advantage of Last Write Wins Element Set is that it allows an element to be reinserted after having been removed.

Notes
---
Wikipedia page on CRDT : https://en.wikipedia.org/wiki/Conflict-free_replicated_data_type

Consistency without consensus in production systems by Peter Bourgon : https://www.youtube.com/watch?v=em9zLzM8O7c

A comprehensive study of Convergent and Commutative Replicated Data Types : https://hal.inria.fr/file/index/docid/555588/filename/techreport.pdf

Roshi: a CRDT system for timestamped events : https://developers.soundcloud.com/blog/roshi-a-crdt-system-for-timestamped-events

CRDT notes by Paul Frazee : https://github.com/pfrazee/crdt_notes

*/
import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type ITimeSet[T comparable] interface {
	Store(T, time.Time) error
	When(T) (time.Time, bool)
	Each(func(T, time.Time) error) error
	Size() int
}

type TimeSet[T comparable] struct {
	mu      sync.RWMutex
	Records map[T]time.Time `json:"records"`
}

type ILastWriterWinsSet[T comparable] interface {
	Add(T, time.Time) error
	Remove(T, time.Time) error
	Has(T) bool
	Get() ([]T, error)
	Merge(ILastWriterWinsSet[T]) error
	Added() ITimeSet[T]
	Deleted() ITimeSet[T]
}

type LastWriterWinsSet[T comparable] struct {
	additions ITimeSet[T]
	deletions ITimeSet[T]
}

type ILastWriterWinsGraph[T comparable] interface {
	AddVertex(T) error
	Vertices() ([]T, error)
	DeleteVertex(T) error
	HasVertex(T) bool
	AddEdge(v1, v2 T) error
	DeleteEdge(v1, v2 T) error
	HasEdge(v1, v2 T) bool
	VertexEdges(v T) ([]T, error)
	FindPath(v1, v2 T) ([]T, error)
	Merge(ILastWriterWinsGraph[T]) error
	getVertices() ILastWriterWinsSet[T]
	getEdges() map[T]ILastWriterWinsSet[T]
}

// LastWriterWinsGraph is a structure for a graph with vertices and edges based on LWW sets
type LastWriterWinsGraph[T comparable] struct {
	mu       sync.RWMutex
	vertices ILastWriterWinsSet[T]
	edges    map[T]ILastWriterWinsSet[T]
}

// Store an element in the set if one of the following condition is met: we haven't encountered that element or the timestamp of the existing element is lesser than the fresh one
func (s *TimeSet[T]) Store(value T, t time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	when, ok := s.Records[value]
	if !ok || (ok && t.After(when)) {
		s.Records[value] = t
	}

	return nil
}

// When returns the timestamp of a given element if it exists
func (s *TimeSet[T]) When(value T) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	when, ok := s.Records[value]
	return when, ok
}

func (s *TimeSet[T]) Each(fn func(el T, when time.Time) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for record, when := range s.Records {
		if err := fn(record, when); err != nil {
			return err
		}
	}

	return nil
}

func (s *TimeSet[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	size := 0
	for range s.Records {
		size++
	}

	return size
}

func NewTimeSet[T comparable]() ITimeSet[T] {
	return &TimeSet[T]{Records: make(map[T]time.Time)}
}

func (s *LastWriterWinsSet[T]) Add(value T, t time.Time) error {
	return s.additions.Store(value, t)
}

func (s *LastWriterWinsSet[T]) Added() ITimeSet[T] {
	return s.additions
}

func (s *LastWriterWinsSet[T]) Remove(value T, t time.Time) error {
	return s.deletions.Store(value, t)
}

func (s *LastWriterWinsSet[T]) Deleted() ITimeSet[T] {
	return s.deletions
}

func (s *LastWriterWinsSet[T]) Has(value T) bool {
	when, wasAdded := s.additions.When(value)
	wasRemoved := s.isRemoved(value, when)
	return wasAdded && !wasRemoved
}

// isRemoved checks if an element is marked for removal
func (s *LastWriterWinsSet[T]) isRemoved(value T, since time.Time) bool {
	when, removed := s.deletions.When(value)
	if !removed {
		return false
	}

	if since.Before(when) {
		return true
	}

	return false
}

// Get returns set content
func (s *LastWriterWinsSet[T]) Get() ([]T, error) {
	var result []T

	err := s.additions.Each(
		func(record T, when time.Time) error {
			removed := s.isRemoved(record, when)
			if !removed {
				result = append(result, record)
			}
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Merge additions and deletions from other LastWriterWinsSet into current set
func (s *LastWriterWinsSet[T]) Merge(other ILastWriterWinsSet[T]) error {
	err := other.Added().Each(
		func(record T, when time.Time) error {
			if err := s.Add(record, when); err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	err = other.Deleted().Each(
		func(record T, when time.Time) error {
			if err := s.Remove(record, when); err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewLastWriterWinsSet[T comparable]() ILastWriterWinsSet[T] {
	return &LastWriterWinsSet[T]{additions: NewTimeSet[T](), deletions: NewTimeSet[T]()}
}

func NewLastWriterWinsGraph[T comparable]() ILastWriterWinsGraph[T] {
	return &LastWriterWinsGraph[T]{vertices: NewLastWriterWinsSet[T]()}
}

func (g *LastWriterWinsGraph[T]) getVertices() ILastWriterWinsSet[T] {
	return g.vertices
}

func (g *LastWriterWinsGraph[T]) getEdges() map[T]ILastWriterWinsSet[T] {
	return g.edges
}

func (g *LastWriterWinsGraph[T]) AddVertex(vertex T) error {
	return g.vertices.Add(vertex, time.Now())
}

func (g *LastWriterWinsGraph[T]) Vertices() ([]T, error) {
	return g.vertices.Get()
}

func (g *LastWriterWinsGraph[T]) DeleteVertex(vertex T) error {
	return g.vertices.Remove(vertex, time.Now())
}

func (g *LastWriterWinsGraph[T]) HasVertex(vertex T) bool {
	return g.vertices.Has(vertex)
}

func (g *LastWriterWinsGraph[T]) AddEdge(vertex1, vertex2 T) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.vertices.Has(vertex1) {
		return fmt.Errorf("cannot add edge, missing node in graph: %v", vertex1)
	} else if !g.vertices.Has(vertex2) {
		return fmt.Errorf("cannot add edge, missing node in graph: %v", vertex2)
	}

	if g.edges == nil {
		g.edges = make(map[T]ILastWriterWinsSet[T])
	}

	if _, ok := g.edges[vertex1]; !ok {
		g.edges[vertex1] = NewLastWriterWinsSet[T]()
	}

	if err := g.edges[vertex1].Add(vertex2, time.Now()); err != nil {
		return err
	}

	if _, ok := g.edges[vertex2]; !ok {
		g.edges[vertex2] = NewLastWriterWinsSet[T]()
	}

	if err := g.edges[vertex2].Add(vertex1, time.Now()); err != nil {
		return err
	}
	return nil
}

func (g *LastWriterWinsGraph[T]) DeleteEdge(vertex1, vertex2 T) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.edges == nil {
		g.edges = make(map[T]ILastWriterWinsSet[T])
	}

	if _, ok := g.edges[vertex1]; !ok {
		g.edges[vertex1] = NewLastWriterWinsSet[T]()
	}

	if err := g.edges[vertex1].Remove(vertex2, time.Now()); err != nil {
		return err
	}

	if _, ok := g.edges[vertex2]; !ok {
		g.edges[vertex2] = NewLastWriterWinsSet[T]()
	}

	if err := g.edges[vertex2].Remove(vertex1, time.Now()); err != nil {
		return err
	}
	return nil
}

func (g *LastWriterWinsGraph[T]) HasEdge(vertex1, vertex2 T) bool {
	return g.edges[vertex1].Has(vertex2) && g.edges[vertex2].Has(vertex1)
}

func (g *LastWriterWinsGraph[T]) VertexEdges(vertex T) ([]T, error) {
	if !g.HasVertex(vertex) {
		return nil, errors.New("cannot query for edges, vertex does not exist")
	}
	return g.edges[vertex].Get()
}

func (g *LastWriterWinsGraph[T]) FindPath(vertex1, vertex2 T) ([]T, error) {
	if !g.vertices.Has(vertex1) {
		return nil, fmt.Errorf("cannot find path, missing node in graph: %v", vertex1)
	} else if !g.vertices.Has(vertex2) {
		return nil, fmt.Errorf("cannot find path, missing node in graph: %v", vertex2)
	}

	seen := NewLastWriterWinsSet[T]()
	var emptyPath []T
	_, path, err := g.findPathRecursive(vertex1, vertex2, seen, emptyPath)
	if err != nil {
		return nil, err
	}

	return path, nil
}

func (g *LastWriterWinsGraph[T]) findPathRecursive(vertex1, vertex2 T, visited ILastWriterWinsSet[T], currentPath []T) (ILastWriterWinsSet[T], []T, error) {
	currentPath = append(currentPath, vertex1)
	if err := visited.Add(vertex1, time.Now()); err != nil {
		return nil, nil, err
	}

	if vertex1 == vertex2 {
		return visited, currentPath, nil
	}

	edges, err := g.edges[vertex1].Get()
	if err != nil {
		return nil, nil, err
	}

	for _, vertex := range edges {
		if visited.Has(vertex) {
			continue
		}

		newVisited, newPath, err := g.findPathRecursive(vertex, vertex2, visited, currentPath)
		if err != nil {
			return nil, nil, err
		}

		if newVisited.Has(vertex2) {
			currentPath = newPath
			visited = newVisited
			break
		}

	}

	return visited, currentPath, nil
}

func (g *LastWriterWinsGraph[T]) Merge(other ILastWriterWinsGraph[T]) error {
	if other == nil {
		return errors.New("cannot merge, other graph is nil")
	}

	if err := g.vertices.Merge(other.getVertices()); err != nil {
		return err
	}

	for otherVertex, otherEdges := range other.getEdges() {
		if currentEdges, ok := g.edges[otherVertex]; ok {
			if err := currentEdges.Merge(otherEdges); err != nil {
				return err
			}
			continue
		}

		g.edges[otherVertex] = otherEdges
	}
	return nil
}
