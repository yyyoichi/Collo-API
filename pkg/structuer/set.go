package structuer

type Set[V any] struct {
	m   map[any]V
	key func(V) any
}

// [V]型のSetを作成する。[key]はVからキーとなる値を取り出すための関数。
func NewSet[V any](key func(V) any) *Set[V] {
	return &Set[V]{
		m:   make(map[any]V),
		key: key,
	}
}
func (s *Set[V]) Add(v V) {
	s.m[s.key(v)] = v
}
func (s *Set[V]) ToSlice() []V {
	data := make([]V, len(s.m))
	var i int
	for _, v := range s.m {
		data[i] = v
		i++
	}
	return data
}
