package goskl

import "fmt"

var (
	ErrConvertFail = fmt.Errorf("value convert failed")
)

const (
	EqualTo     = 0
	LessThan    = -1
	GreaterThan = 1
)

// Comparator is used to compare two values.
type Comparator interface {
	// Compare should return 0 if l == r; return -1 if l < r; 1 if l > r.
	Compare(l, r interface{}) (int, error)
	Default() interface{}
}

type stringCmp struct{}

func (c stringCmp) convert(s interface{}) (string, error) {
	r, ok := s.(string)
	if !ok {
		return "", ErrConvertFail
	}
	return r, nil
}

func (c stringCmp) Compare(l, r interface{}) (int, error) {
	lhs, errL := c.convert(l)
	rhs, errR := c.convert(r)
	if errL != nil || errR != nil {
		return -1, ErrConvertFail
	}

	if lhs == rhs {
		return EqualTo, nil
	} else if lhs < rhs {
		return LessThan, nil
	}
	return GreaterThan, nil
}

func (s stringCmp) Default() interface{} {
	return ""
}

type intCmp struct{}

func (c intCmp) convert(s interface{}) (int, error) {
	r, ok := s.(int)
	if !ok {
		return 0, ErrConvertFail
	}
	return r, nil
}

func (c intCmp) Compare(l, r interface{}) (int, error) {
	lhs, errL := c.convert(l)
	rhs, errR := c.convert(r)
	if errL != nil || errR != nil {
		return -1, ErrConvertFail
	}

	if lhs == rhs {
		return EqualTo, nil
	} else if lhs < rhs {
		return LessThan, nil
	}
	return GreaterThan, nil
}

func (s intCmp) Default() interface{} {
	return 0
}

var (
	StringCmp = stringCmp{}
	IntCmp    = intCmp{}
)
