package goskl

type Pair struct {
	Key   string
	Value interface{}
}

// Node represents a node in skiplist
type Node struct {
	*Pair
	nexts []*Node
}

func NewNode(k string, v interface{}, level int) *Node {
	return &Node{
		Pair: &Pair{
			Key:   k,
			Value: v,
		},
		nexts: make([]*Node, level),
	}
}
