package set

type Set[T comparable] interface {
	// From https://github.com/golang/go/discussions/47331#discussion-3471930
	Add(vs ...T)
	AddSet(vs Set[T])
	Remove(vs ...T)
	RemoveSet(vs Set[T])
	Contains(v T) bool
	ContainsAny(vs Set[T]) bool
	ContainsAll(vs Set[T]) bool
	Values() []T
	Equal(vs Set[T]) bool
	Clear()
	Filter(keep func(T) bool) Set[T]
	Len() int
	Clone() Set[T]
	Do(f func(v T) bool)
}

var EXISTENCE struct{}

type SimpleSet[T comparable] map[T]struct{}

func New[T comparable](vs ...T) SimpleSet[T] {
	s := make(SimpleSet[T], len(vs))
	s.Add(vs...)
	return s
}

// From https://github.com/golang/go/discussions/47331#discussion-3471930
func Of[T comparable](vs ...T) SimpleSet[T] {
	return New[T](vs...)
}

func (s SimpleSet[T]) Add(vs ...T) {
	for _, v := range vs {
		s[v] = EXISTENCE
	}
}

func (s SimpleSet[T]) AddSet(vs Set[T]) {
	vs.Do(func(v T) bool {
		s.Add(v)
		return true
	})
}

func (s SimpleSet[T]) Remove(vs ...T) {
	for _, v := range vs {
		delete(s, v)
	}
}

func (s SimpleSet[T]) RemoveSet(vs Set[T]) {
	vs.Do(func(v T) bool {
		s.Remove(v)
		return true
	})
}

func (s SimpleSet[T]) Contains(v T) bool {
	if s != nil {
		_, ok := s[v]
		return ok
	}
	return false
}

func (s SimpleSet[T]) ContainsAny(vs Set[T]) bool {
	ok := false
	vs.Do(func(v T) bool {
		if s.Contains(v) {
			ok = true
			return false
		}
		return true
	})
	return ok
}

func (s SimpleSet[T]) ContainsAll(vs Set[T]) bool {
	ok := true
	vs.Do(func(v T) bool {
		if !s.Contains(v) {
			ok = false
			return false
		}
		return true
	})
	return ok
}

func (s SimpleSet[T]) Values() []T {
	if s == nil {
		return nil
	}
	vs := make([]T, 0, s.Len())
	s.Do(func(v T) bool {
		vs = append(vs, v)
		return true
	})
	return vs
}

func (s SimpleSet[T]) Equal(vs Set[T]) bool {
	if s.Len() != vs.Len() {
		return false
	}
	return s.ContainsAll(vs)
}

func (s SimpleSet[T]) Clear() {
	s.Do(func(v T) bool {
		s.Remove(v)
		return true
	})
}

func (s SimpleSet[T]) Filter(f func(v T) bool) Set[T] {
	if s == nil {
		return (SimpleSet[T])(nil)
	}
	vs := New[T]()
	s.Do(func(v T) bool {
		if f(v) {
			vs.Add(v)
		}
		return true
	})
	return vs
}

func (s SimpleSet[T]) Len() int {
	return len(s)
}

func (s SimpleSet[T]) Clone() Set[T] {
	return s.Filter(func(_ T) bool { return true })
}

func (s SimpleSet[T]) Do(f func(v T) bool) {
	if s != nil && f != nil {
		for v := range s {
			if !f(v) {
				break
			}
		}
	}
}
