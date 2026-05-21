package live

type Ring[T any] struct {
	items []T
	next  int
	full  bool
}

func NewRing[T any](size int) *Ring[T] {
	if size < 1 {
		size = 1
	}
	return &Ring[T]{items: make([]T, size)}
}

func (r *Ring[T]) Push(v T) {
	r.items[r.next] = v
	r.next = (r.next + 1) % len(r.items)
	if r.next == 0 {
		r.full = true
	}
}

func (r *Ring[T]) Snapshot() []T {
	var out []T
	if r.full {
		out = append(out, r.items[r.next:]...)
		out = append(out, r.items[:r.next]...)
	} else {
		out = append(out, r.items[:r.next]...)
	}
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}
