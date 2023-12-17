# go-skl

go-skl is a straightforward skiplist implementation in golang.

Implementations are:
```golang
Put(k string, v interface{}) bool
Get(k string) (interface{}, bool)
Remove(k string) bool
List() []*Pair
Size() int
First() *Pair
```