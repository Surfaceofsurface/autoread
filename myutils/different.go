package myutils

type MarkAble interface {
	Mark() string
}

// type Set[T MarkAble] struct {
// 	container map[string]T
// }

func Diff[T MarkAble](a []T, b []T) []T {
	c := make([]T, 0)
	m := make(map[string][0]byte, len(b))
	for _, v := range b {
		m[v.Mark()] = [0]byte{}
	}

	for _, v := range a {
		if _, ok := m[v.Mark()]; !ok {
			c = append(c, v)
		}
	}
	return c
}
