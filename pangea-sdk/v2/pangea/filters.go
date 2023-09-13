package pangea

type Filter map[string]any

type FilterBase struct {
	f Filter
}

func (fb FilterBase) Filter() Filter {
	return fb.f
}

func NewFilterBase(f Filter) *FilterBase {
	return &FilterBase{
		f: f,
	}
}

type FilterCommon struct {
	name string
	m    *Filter
}

type FilterEqual[T any] struct {
	FilterCommon
}

func NewFilterEqual[T any](name string, filter *Filter) *FilterEqual[T] {
	return &FilterEqual[T]{
		FilterCommon{
			name: name,
			m:    filter,
		},
	}
}

func (f *FilterEqual[T]) Set(value *T) {
	if value != nil {
		(*f.m)[f.name] = value
	} else {
		delete(*f.m, f.name)
	}
}

func (f *FilterEqual[T]) Get() *T {
	if v, ok := (*f.m)[f.name].(*T); ok {
		return v
	}
	return nil
}

type FilterMatch[T any] struct {
	FilterEqual[T]
}

func NewFilterMatch[T any](name string, filter *Filter) *FilterMatch[T] {
	return &FilterMatch[T]{
		FilterEqual: *NewFilterEqual[T](name, filter),
	}
}

func (f *FilterMatch[T]) SetIn(value []T) {
	if value != nil {
		(*f.m)[f.name+"__in"] = value
	} else {
		delete(*f.m, f.name+"__in")
	}
}

func (f *FilterMatch[T]) In() []T {
	if v, ok := (*f.m)[f.name+"__in"].([]T); ok {
		return v
	}
	return nil
}

func (f *FilterMatch[T]) SetContains(value []T) {
	if value != nil {
		(*f.m)[f.name+"__contains"] = value
	} else {
		delete(*f.m, f.name+"__contains")
	}
}

func (f *FilterMatch[T]) Contains() []T {
	if v, ok := (*f.m)[f.name+"__contains"].([]T); ok {
		return v
	}
	return nil
}

type FilterRange[T any] struct {
	FilterEqual[T]
}

func NewFilterRange[T any](name string, filter *Filter) *FilterRange[T] {
	return &FilterRange[T]{
		FilterEqual: *NewFilterEqual[T](name, filter),
	}
}

func (f *FilterRange[T]) SetLessThan(value *T) {
	if value != nil {
		(*f.m)[f.name+"__lt"] = value
	} else {
		delete(*f.m, f.name+"__lt")
	}
}

func (f *FilterRange[T]) LessThan() *T {
	if v, ok := (*f.m)[f.name+"__lt"].(*T); ok {
		return v
	}
	return nil
}

func (f *FilterRange[T]) SetLessThanEqual(value *T) {
	if value != nil {
		(*f.m)[f.name+"__lte"] = value
	} else {
		delete(*f.m, f.name+"__lte")
	}
}

func (f *FilterRange[T]) LessThanEqual() *T {
	if v, ok := (*f.m)[f.name+"__lte"].(*T); ok {
		return v
	}
	return nil
}

func (f *FilterRange[T]) SetGreaterThan(value *T) {
	if value != nil {
		(*f.m)[f.name+"__gt"] = value
	} else {
		delete(*f.m, f.name+"__gt")
	}
}

func (f *FilterRange[T]) GreaterThan() *T {
	if v, ok := (*f.m)[f.name+"__gt"].(*T); ok {
		return v
	}
	return nil
}

func (f *FilterRange[T]) SetGreaterThanEqual(value *T) {
	if value != nil {
		(*f.m)[f.name+"__gte"] = value
	} else {
		delete(*f.m, f.name+"__gte")
	}
}

func (f *FilterRange[T]) GreaterThanEqual() *T {
	if v, ok := (*f.m)[f.name+"__gte"].(*T); ok {
		return v
	}
	return nil
}
