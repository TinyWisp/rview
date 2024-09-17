package rview

type Ref[T any] struct {
	value T
}

func (r *Ref[T]) Get() (val T) {
	return r.value
}

func (r *Ref[T]) Set(val T) {
	r.value = val
}