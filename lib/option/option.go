package option

type Option[T any] struct {
	set bool
	val T
}

func OptionNil[T any]() Option[T] {
	r := Option[T]{}
	r.set = false
	return r
}

func OptionVal[T any](val T) Option[T] {
	return Option[T]{true, val}
}

func (o Option[T]) GetSet() bool {
	return o.set
}

func (o Option[T]) GetVal() T {
	return o.val
}
