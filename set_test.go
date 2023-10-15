package set

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"
)

func TestNewAndOf(t *testing.T) {
	type test[T comparable] struct {
		args     []T
		expected Set[T]
	}
	testsForInt := []test[int]{
		{
			args:     []int{},
			expected: &SimpleSet[int]{m: map[int]struct{}{}},
		},
		{
			args: []int{1},
			expected: &SimpleSet[int]{
				m: map[int]struct{}{
					1: EXISTENCE,
				},
			},
		},
		{
			args: []int{1, 2, 3, 4, 1, 2, 3, 4, 5, 6},
			expected: &SimpleSet[int]{
				m: map[int]struct{}{
					1: EXISTENCE,
					2: EXISTENCE,
					3: EXISTENCE,
					4: EXISTENCE,
					5: EXISTENCE,
					6: EXISTENCE,
				},
			},
		},
	}

	for _, test := range testsForInt {
		s := New(test.args...)
		if !reflect.DeepEqual(s, test.expected) {
			t.Errorf("New() expected %v but got %v", test.expected, s)
		}
		s = Of(test.args...)
		if !reflect.DeepEqual(s, test.expected) {
			t.Errorf("Of() expected %v but got %v", test.expected, s)
		}
	}

	testsForString := []test[string]{
		{
			args:     []string{},
			expected: &SimpleSet[string]{m: map[string]struct{}{}},
		},
	}

	for _, test := range testsForString {
		s := New(test.args...)
		if !reflect.DeepEqual(s, test.expected) {
			t.Errorf("New() expected %v but got %v", test.expected, s)
		}
		s = Of(test.args...)
		if !reflect.DeepEqual(s, test.expected) {
			t.Errorf("Of() expected %v but got %v", test.expected, s)
		}
	}
}

func TestEqual(t *testing.T) {
	for _, test := range []struct {
		src1     Set[string]
		src2     Set[string]
		expected bool
	}{
		{
			src1:     (*SimpleSet[string])(nil),
			src2:     (*SimpleSet[string])(nil),
			expected: true,
		},
		{
			src1:     (*SimpleSet[string])(nil),
			src2:     New[string](),
			expected: true,
		},
		{
			src1:     New[string](),
			src2:     (*SimpleSet[string])(nil),
			expected: true,
		},
		{
			src1:     New[string](),
			src2:     New[string](),
			expected: true,
		},
		{
			src1:     New[string](),
			src2:     New[string]("a"),
			expected: false,
		},
		{
			src1:     New[string]("a"),
			src2:     New[string](),
			expected: false,
		},
		{
			src1:     New[string]("a"),
			src2:     New[string]("a"),
			expected: true,
		},
		{
			src1:     New[string]("a"),
			src2:     New[string]("a"),
			expected: true,
		},
		{
			src1:     New[string]("b"),
			src2:     New[string]("a"),
			expected: false,
		},
		{
			src1:     New[string]("a"),
			src2:     New[string]("b"),
			expected: false,
		},
		{
			src1:     New[string]("a", "b"),
			src2:     New[string]("a"),
			expected: false,
		},
		{
			src1:     New[string]("a"),
			src2:     New[string]("a", "b"),
			expected: false,
		},
		{
			src1:     New[string]("a", "b"),
			src2:     New[string]("b", "c"),
			expected: false,
		},
		{
			src1:     New[string]("a", "b"),
			src2:     New[string]("b", "a"),
			expected: true,
		},
	} {
		if result := test.src1.Equal(test.src2); result != test.expected {
			t.Errorf("Equal() expected %v but got %v", test.expected, result)
		}
	}
}

func TestModification(t *testing.T) {
	initargs := []string{
		"one", "two", "three", "two", "one", "four", "five", "six", "seven", "nine", "eight",
	}

	tests := []struct {
		addArgs        []string
		removeArgs     []string
		addExpected    Set[string]
		removeExpected Set[string]
	}{
		{
			addArgs:        []string{},
			removeArgs:     []string{},
			addExpected:    New("one", "two", "three", "two", "one", "four", "five", "six", "seven", "nine", "eight"),
			removeExpected: New("one", "two", "three", "two", "one", "four", "five", "six", "seven", "nine", "eight"),
		},
		{
			addArgs:        []string{"nine", "ten", "eleven"},
			removeArgs:     []string{"nine", "ten", "eleven"},
			addExpected:    New("one", "two", "three", "two", "one", "four", "five", "six", "seven", "nine", "eight", "ten", "eleven"),
			removeExpected: New("one", "two", "three", "two", "one", "four", "five", "six", "seven", "eight"),
		},
	}

	for _, test := range tests {
		a := New(initargs...)
		a.Add(test.addArgs...)
		if !reflect.DeepEqual(a, test.addExpected) {
			t.Errorf("Add() expected %v but got %v", test.addExpected, a)
		}
		a = New(initargs...)
		a.AddSet(New(test.addArgs...))
		if !reflect.DeepEqual(a, test.addExpected) {
			t.Errorf("AddSet() expected %v but got %v", test.addExpected, a)
		}
		r := New(initargs...)
		r.Remove(test.removeArgs...)
		if !reflect.DeepEqual(r, test.removeExpected) {
			t.Errorf("Remove() expected %v but got %v", test.removeExpected, r)
		}
		r = New(initargs...)
		r.RemoveSet(New(test.removeArgs...))
		if !reflect.DeepEqual(r, test.removeExpected) {
			t.Errorf("RemoveSet() expected %v but got %v", test.removeExpected, r)
		}
	}
}

func TestPredicates(t *testing.T) {
	s := New("one", "two", "three", "four", "five")

	for _, test := range []struct {
		arg      string
		expected bool
	}{
		{
			arg:      "two",
			expected: true,
		},
		{
			arg:      "six",
			expected: false,
		},
	} {
		if result := s.Contains(test.arg); result != test.expected {
			t.Errorf("Contains() expected %v but got %v", test.expected, result)
		}
	}

	for _, test := range []struct {
		arg         Set[string]
		allExpected bool
		anyExpected bool
	}{
		{
			arg:         New[string](),
			allExpected: true,
			anyExpected: false,
		},
		{
			arg:         New("three"),
			allExpected: true,
			anyExpected: true,
		},
		{
			arg:         New("four", "five"),
			allExpected: true,
			anyExpected: true,
		},
		{
			arg:         New("four", "five", "six"),
			allExpected: false,
			anyExpected: true,
		},
		{
			arg:         New("six"),
			allExpected: false,
			anyExpected: false,
		},
		{
			arg:         New("six", "seven"),
			allExpected: false,
			anyExpected: false,
		},
	} {
		if result := s.ContainsAll(test.arg); result != test.allExpected {
			t.Errorf("ContainsAll() expected %v but got %v", test.allExpected, result)
		}
		if result := s.ContainsAny(test.arg); result != test.anyExpected {
			t.Errorf("ContainsAny() expected %v but got %v", test.anyExpected, result)
		}
	}
}

func TestValues(t *testing.T) {
	for _, test := range []struct {
		src      Set[string]
		expected []string
	}{
		{
			src:      (*SimpleSet[string])(nil),
			expected: nil,
		},
		{
			src:      New[string](),
			expected: []string{},
		},
		{
			src:      New("a"),
			expected: []string{"a"},
		},
		{
			src:      New("a", "b"),
			expected: []string{"a", "b"},
		},
		{
			src:      New("c", "b", "a"),
			expected: []string{"a", "b", "c"},
		},
	} {
		result := test.src.Values()
		if result != nil {
			sort.Strings(result)
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Value() expected %v but got %v", test.expected, result)
		}
	}
}

func TestLen(t *testing.T) {
	for _, test := range []struct {
		src      Set[string]
		expected int
	}{
		{
			src:      (*SimpleSet[string])(nil),
			expected: 0,
		},
		{
			src:      New[string](),
			expected: 0,
		},
		{
			src:      New("a"),
			expected: 1,
		},
		{
			src:      New("a", "b"),
			expected: 2,
		},
		{
			src:      New("a", "b", "a"),
			expected: 2,
		},
	} {
		if l := test.src.Len(); l != test.expected {
			t.Errorf("Len() expected %v but got %v", test.expected, l)
		}
	}
}

func TestCloneClear(t *testing.T) {
	for _, test := range []struct {
		src      Set[string]
		expected Set[string]
		cleared  Set[string]
	}{
		{
			src:      (*SimpleSet[string])(nil),
			expected: (Set[string])(nil),
			cleared:  (Set[string])(nil),
		},
		{
			src:      New[string](),
			expected: New[string](),
			cleared:  New[string](),
		},
		{
			src:      New("a"),
			expected: New("a"),
			cleared:  New[string](),
		},
		{
			src:      New("a", "b", "c"),
			expected: New("a", "b", "c"),
			cleared:  New[string](),
		},
	} {
		if result := test.src.Clone(); !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Clone() expected %v but got %v", test.expected, result)
		} else {
			if result != nil {
				result.(*SimpleSet[string]).Clear()
				if !reflect.DeepEqual(result, test.cleared) {
					t.Errorf("Clear() expected %v but got %v from %v", test.cleared, result, test.src)
				}
			}
		}
	}
}

func TestFilter(t *testing.T) {
	isEven := func(v int) bool {
		return v%2 == 0
	}
	for _, test := range []struct {
		src      Set[int]
		expected Set[int]
		pred     func(v int) bool
	}{
		{
			src:      (*SimpleSet[int])(nil),
			expected: (Set[int])(nil),
			pred:     isEven,
		},
		{
			src:      New[int](),
			expected: New[int](),
			pred:     isEven,
		},
		{
			src:      New(1),
			expected: New[int](),
			pred:     isEven,
		},
		{
			src:      New(1, 2, 3),
			expected: New(2),
			pred:     isEven,
		},
	} {
		if result := test.src.Filter(test.pred); !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Filter() expected %v but got %v from %v", test.expected, result, test.src)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	for _, test := range []struct {
		src      *SimpleSet[string]
		expected string
	}{
		{
			src:      New[string](),
			expected: `[]`,
		},
		{
			src:      New("a"),
			expected: `["a"]`,
		},
		{
			src:      New("a", "b"),
			expected: `["a","b"]`,
		},
		{
			src:      New("a", "b", "c"),
			expected: `["a","b","c"]`,
		},
	} {
		if result, err := test.src.MarshalJSON(); err != nil {
			t.Errorf("MarshalJSON() unexpected error: %s", err)
		} else {
			var vs []string
			if err := json.Unmarshal(result, &vs); err != nil {
				t.Errorf("MarshalJSON() result cannot be unmarshal to []string but got error: %s", err)
			}
			if result := New(vs...); !reflect.DeepEqual(result, test.src) {
				t.Errorf("MarshalJSON() result unmarshal expected %v but got %v", result, test.src)
			}
		}
	}
}

func TestUnmarshalJson(t *testing.T) {
	for _, test := range []struct {
		src      string
		expected Set[string]
	}{
		{
			src:      `[]`,
			expected: New[string](),
		},
		{
			src:      `["a"]`,
			expected: New("a"),
		},
		{
			src:      `["a","b"]`,
			expected: New("a", "b"),
		},
		{
			src:      `["c","b","a"]`,
			expected: New("a", "b", "c"),
		},
	} {
		var result SimpleSet[string]
		if err := json.Unmarshal([]byte(test.src), &result); err != nil {
			t.Errorf("Unmarshal() unexpected error: %s", err)
		} else if !reflect.DeepEqual(&result, test.expected) {
			t.Errorf("Unmarshal() expected %v but got %v", test.expected, &result)
		}
	}
}
