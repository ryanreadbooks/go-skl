package goskl

import "math/rand"

const (
	DefaultMaxLevel = 64
)

// Skiplist represents a skiplist structure
type Skiplist struct {
	head       *Node
	size       int
	maxLevel   int
	comparator Comparator
}

// New returns a new skiplist instance with whose max level is maxLevel.
func New(maxLevel int, comparator Comparator) *Skiplist {
	if maxLevel < 1 {
		panic("maxLevel should be at least 1")
	}

	l := &Skiplist{
		head:       NewNode(comparator.Default(), nil, 1),
		size:       0,
		maxLevel:   maxLevel,
		comparator: comparator,
	}

	return l
}

// Level returns the level of the skiplist, which is the level of the head node.
func (s *Skiplist) Level() int {
	return len(s.head.nexts)
}

// Put puts a k-v pair into skiplist. If key already exists, value is updated.
func (s *Skiplist) Put(k interface{}, v interface{}) (res bool) {
	pre := s.findLastLeThan(k)
	if pre == nil {
		res = false // theoretically this will not happen
		return
	}

	// update existing key-value
	if s.isEqual(pre.Key, k) {
		pre.Value = v
		return true
	}

	// this is an add operation
	nodeLevel := s.determineLevel()
	newNode := NewNode(k, v, nodeLevel)

	s.adjustLevel(nodeLevel)

	// find the gap and perform the insert operation
	cur := s.head
	for l := s.Level() - 1; l >= 0; l-- {
		for cur.nexts[l] != nil && s.isLess(cur.nexts[l].Key, k) {
			cur = cur.nexts[l]
		}

		if l < nodeLevel {
			// adjust the pointers to insert new node at each level
			newNode.nexts[l] = cur.nexts[l]
			cur.nexts[l] = newNode
		}
	}
	s.size++

	return true
}

func (s *Skiplist) adjustLevel(wanted int) {
	for s.Level() < wanted {
		s.head.nexts = append(s.head.nexts, nil)
	}
}

// Get returns the value of the given key if it has present.
func (s *Skiplist) Get(k interface{}) (interface{}, bool) {
	pre := s.findLastLeThan(k)
	if pre != nil && s.isEqual(pre.Key, k) {
		return pre.Value, true
	}

	return nil, false
}

// Has checks whether the given key is existing in skiplist
func (s *Skiplist) Has(k interface{}) bool {
	_, ok := s.Get(k)
	return ok
}

// Remove removes the given key in skiplist.
func (s *Skiplist) Remove(k interface{}) bool {
	pre := s.findLastLeThan(k)
	if pre == nil {
		return false
	}

	// key not exist
	if !s.isEqual(pre.Key, k) {
		return false
	}

	cur := s.head
	for level := s.Level() - 1; level >= 0; level-- {
		for cur.nexts[level] != nil && s.isLess(cur.nexts[level].Key, k) {
			cur = cur.nexts[level]
		}

		if cur.nexts[level] == nil || s.isGreater(cur.nexts[level].Key, k) {
			continue
		}

		// cur == k, we adjust the pointers of node
		cur.nexts[level] = cur.nexts[level].nexts[level]
	}

	// shrink skiplist level
	height := 0
	for i := s.Level() - 1; i > 0 && s.head.nexts[i] == nil; i-- {
		height++
	}
	s.head.nexts = s.head.nexts[:s.Level()-height]
	s.size--

	return true
}

// determineLevel determines the level of a new node in a random manner.
func (s *Skiplist) determineLevel() int {
	var level int = 1
	for rand.Intn(2) > 0 {
		level++
	}
	if level >= s.maxLevel {
		level = s.maxLevel
	}
	return level
}

// find the last node whose key is smaller than or equal to the given key.
func (s *Skiplist) findLastLeThan(k interface{}) *Node {
	cur := s.head
	// from top to down
	for level := s.Level() - 1; level >= 0; level-- {
		// from left to right
		for cur.nexts[level] != nil && s.isLessEqual(cur.nexts[level].Key, k) {
			cur = cur.nexts[level]
		}
	}

	return cur
}

// List lists all the elements in the skiplist
func (s *Skiplist) List() []*Pair {
	res := make([]*Pair, 0)
	cur := s.head.nexts[0]
	for cur != nil {
		res = append(res, &Pair{Key: cur.Key, Value: cur.Value})
		cur = cur.nexts[0]
	}
	return res
}

// Size returns the length of the skiplist
func (s *Skiplist) Size() int {
	return s.size
}

// First returns the first element in the skiplist
func (s *Skiplist) First() *Pair {
	if s.head.nexts[0] == nil {
		return nil
	}

	// return copy of it
	return &Pair{
		Key:   s.head.nexts[0].Key,
		Value: s.head.nexts[0].Value,
	}
}

func (s *Skiplist) isLessEqual(lhs, rhs interface{}) bool {
	res, err := s.comparator.Compare(lhs, rhs)
	if err != nil {
		return false
	}

	return res != GreaterThan
}

func (s *Skiplist) isLess(lhs, rhs interface{}) bool {
	res, err := s.comparator.Compare(lhs, rhs)
	if err != nil {
		return false
	}

	return res == LessThan
}

func (s *Skiplist) isEqual(lhs, rhs interface{}) bool {
	res, err := s.comparator.Compare(lhs, rhs)
	if err != nil {
		return false
	}

	return res == EqualTo
}

func (s *Skiplist) isGreater(lhs, rhs interface{}) bool {
	res, err := s.comparator.Compare(lhs, rhs)
	if err != nil {
		return false
	}

	return res == GreaterThan
}

func (s *Skiplist) Comparator() Comparator {
	return s.comparator
}