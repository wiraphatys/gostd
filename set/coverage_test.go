package set

import (
	"reflect"
	"testing"
)

func TestForEach(t *testing.T) {
	s := NewFromSlice([]int{3, 1, 2})
	var got []int
	s.ForEach(func(v int) { got = append(got, v) })
	if !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("ForEach = %v, want [1 2 3]", got)
	}
}
