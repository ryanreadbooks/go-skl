# go-skl

go-skl is a straightforward skiplist implementation in golang.

Implementations are:
```golang
Put(k interface{}, v interface{}) bool
Get(k interface{}) (interface{}, bool)
Remove(k interface{}) bool
List() []*Pair
Size() int
First() *Pair
```

**Future works**
- Improve benchmark results.
