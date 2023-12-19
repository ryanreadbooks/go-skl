package goskl_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	goskl "github.com/ryanreadbooks/go-skl"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
)

type TestSkiplistSuite struct {
	suite.Suite
	skl *goskl.Skiplist
}

// customized assertion for convey which test if
func ShouldBeSorted(actual any, expect ...any) string {
	if len(expect) != 0 {
		return "ShouldBeSorted expect should be empty"
	}

	// actual should be sorted in a default manner
	skl, ok := actual.(*goskl.Skiplist)
	if !ok {
		return "ShouldBeSorted should only apply to *goskl.Skiplist type"
	}

	nodes := skl.List()
	for i := 0; i < len(nodes)-1; i++ {
		c, err := skl.Comparator().Compare(nodes[i].Key, nodes[i+1].Key)
		if err != nil {
			return fmt.Sprintf("err should be nil, but got %v", err)
		}
		if c != goskl.LessThan { // >=
			return "skiplist is not sorted"
		}
	}

	return ""
}

func TestSkiplist(t *testing.T) {
	suite.Run(t, new(TestSkiplistSuite))
}

func (s *TestSkiplistSuite) SetupTest() {
	s.skl = goskl.New(16, goskl.StringCmp)
}

func (s *TestSkiplistSuite) TearDownTest() {
	printSkl(s.skl)
	convey.Convey("TearDownTest validation sorted skiplist", s.T(), func() {
		convey.So(s.skl, ShouldBeSorted)
	})
	s.skl = nil // restore
}

func printSkl(skl *goskl.Skiplist) {
	var str strings.Builder
	all := skl.List()
	for i, node := range all {
		str.WriteString(fmt.Sprintf("{%v: %v}", node.Key, node.Value))
		if i != len(all)-1 {
			str.WriteString(" -> ")
		}
	}
	fmt.Println(str.String())
}

func (s *TestSkiplistSuite) TestPut() {
	s.skl.Put("a", "a")
	s.skl.Put("c", "c")
	s.skl.Put("b", "b")
	s.skl.Put("w", "w")
	s.skl.Put("e", "e")

	convey.Convey("Put should be passed", s.T(), func() {
		convey.So(s.skl.Size(), convey.ShouldEqual, 5)
	})
}

// Test updating scenario
func (s *TestSkiplistSuite) TestUpdate() {
	s.skl.Put("a", "a")
	s.skl.Put("c", "c")
	s.skl.Put("b", "b")
	s.skl.Put("w", "w")
	s.skl.Put("e", "e")
	convey.Convey("Update should be passed", s.T(), func() {
		convey.So(s.skl.Size(), convey.ShouldEqual, 5)
		// update
		s.skl.Put("w", "W")
		s.skl.Put("c", "C")
		s.skl.Put("b", "B")
		s.skl.Put("e", "E")
		s.skl.Put("a", "A")
		s.skl.Put("q", "Q")
		convey.So(s.skl.Size(), convey.ShouldEqual, 6)
	})
}

// Test get functionality
func (s *TestSkiplistSuite) TestGet() {
	s.skl.Put("a", "a")
	s.skl.Put("c", "c")
	s.skl.Put("b", "b")
	s.skl.Put("w", "w")
	s.skl.Put("e", "e")

	convey.Convey("Get should be passed", s.T(), func() {
		convey.So(s.skl.Size(), convey.ShouldEqual, 5)

		cases := []struct {
			find      string
			expectVal interface{}
			expect    bool
		}{
			{
				find:      "a",
				expectVal: "a",
				expect:    true,
			},
			{
				find:      "c",
				expectVal: "c",
				expect:    true,
			},
			{
				find:      "b",
				expectVal: "b",
				expect:    true,
			},
			{
				find:      "w",
				expectVal: "w",
				expect:    true,
			},
			{
				find:      "e",
				expectVal: "e",
				expect:    true,
			},
			{
				find:      "qwe",
				expectVal: "qwe",
				expect:    false,
			},
			{
				find:      "123",
				expectVal: "123",
				expect:    false,
			},
			{
				find:      "wqqqq",
				expectVal: "wewewe",
				expect:    false,
			},
		}

		for _, c := range cases {
			got, ok := s.skl.Get(c.find)
			convey.So(ok, convey.ShouldEqual, c.expect)
			if c.expect {
				convey.So(got, convey.ShouldEqual, c.expectVal)
			}
		}
	})
}

func (s *TestSkiplistSuite) TestRemove() {
	s.skl.Put("a", "a")
	s.skl.Put("c", "c")
	s.skl.Put("b", "b")
	s.skl.Put("w", "w")
	s.skl.Put("e", "e")

	convey.Convey("Remove should be passed", s.T(), func() {
		cases := []struct {
			key    string
			expect bool
		}{
			{key: "a", expect: true},
			{key: "w", expect: true},
			{key: "c", expect: true},
			{key: "b", expect: true},
			{key: "e", expect: true},
			{key: "q", expect: false},
			{key: "o", expect: false},
		}

		for _, c := range cases {
			ok := s.skl.Remove(c.key)
			convey.So(ok, convey.ShouldEqual, c.expect)
		}

		for _, c := range cases {
			ok := s.skl.Remove(c.key)
			convey.So(ok, convey.ShouldBeFalse)
		}
	})
}

// test the cases where putting is after remove
func (s *TestSkiplistSuite) TestRemoveThenPut() {
	s.skl.Put("a", "a")
	s.skl.Put("c", "c")
	s.skl.Put("b", "b")
	s.skl.Put("w", "w")
	s.skl.Put("e", "e")

	convey.Convey("Remove then put should be passed", s.T(), func() {
		convey.So(s.skl.Size(), convey.ShouldEqual, 5)

		// remove then put back
		cases := []struct {
			key    string
			expect bool
		}{
			{key: "a", expect: true},
			{key: "w", expect: true},
			{key: "c", expect: true},
			{key: "b", expect: true},
			{key: "e", expect: true},
			{key: "q", expect: false},
			{key: "o", expect: false},
		}

		for _, c := range cases {
			ok := s.skl.Remove(c.key)
			convey.So(ok, convey.ShouldEqual, c.expect)
			// put back
			ok = s.skl.Put(c.key, c.key)
			convey.So(ok, convey.ShouldBeTrue)
		}

		// should be the same the length
		convey.So(s.skl.Size(), convey.ShouldEqual, 7)
	})
}

func (s *TestSkiplistSuite) TestFirst() {
	s.skl.Put("a", "a")
	s.skl.Put("c", "c")
	s.skl.Put("b", "b")
	s.skl.Put("w", "w")
	s.skl.Put("e", "e")

	convey.Convey("First should be passed", s.T(), func() {
		cc := []struct {
			key string
			val interface{}
		}{
			{key: "a", val: "a"},
			{key: "b", val: "b"},
			{key: "c", val: "c"},
			{key: "e", val: "e"},
			{key: "w", val: "w"},
		}

		for _, c := range cc {
			first := s.skl.First()
			convey.So(first.Key, convey.ShouldEqual, c.key)
			convey.So(first.Value, convey.ShouldEqual, c.val)
			s.skl.Remove(c.key)
		}

		// should be empty
		first := s.skl.First()
		convey.So(first, convey.ShouldBeNil)
		convey.So(s.skl.Size(), convey.ShouldEqual, 0)
		convey.So(len(s.skl.List()), convey.ShouldEqual, 0)
	})
}

// Benchmarks
func BenchmarkPut(b *testing.B) {
	skl := goskl.New(goskl.DefaultMaxLevel, goskl.StringCmp)
	keys := make([]string, 0, b.N)
	values := make([]interface{}, 0, b.N)
	for i := 0; i < b.N; i++ {
		keys = append(keys, strconv.Itoa(rand.Int()))
		values = append(values, rand.Int())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		skl.Put(keys[i], values[i])
	}
}

func BenchmarkGet(b *testing.B) {
	skl := goskl.New(goskl.DefaultMaxLevel, goskl.StringCmp)
	keys := make([]string, 0, b.N)
	values := make([]interface{}, 0, b.N)
	for i := 0; i < b.N; i++ {
		keys = append(keys, strconv.Itoa(rand.Int()))
		values = append(values, rand.Int())
	}

	// insert random k-v
	for i := 0; i < b.N; i++ {
		skl.Put(keys[i], values[i])
	}

	// shuffle keys for random get
	rand.Shuffle(b.N, func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		skl.Get(keys[i])
	}
}

func BenchmarkRemove(b *testing.B) {
	skl := goskl.New(goskl.DefaultMaxLevel, goskl.StringCmp)
	keys := make([]string, 0, b.N)
	values := make([]interface{}, 0, b.N)
	for i := 0; i < b.N; i++ {
		keys = append(keys, strconv.Itoa(rand.Int()))
		values = append(values, rand.Int())
	}

	// insert random k-v
	for i := 0; i < b.N; i++ {
		skl.Put(keys[i], values[i])
	}
	// shuffle keys for random remove
	rand.Shuffle(b.N, func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		skl.Remove(keys[i])
	}
}
