package helper

func FindBy[T any](slice []T, predicate func(T) bool) (*T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return &item, true
		}
	}
	return nil, false
}

func IncludesBy[T any](slice []T, predicate func(T) bool) bool {
	_, found := FindBy(slice, predicate)
	return found
}

func All[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

func Some[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

func Concat[T any](slices ...[]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

func MapTo[TIn any, TOut any](slice []TIn, mapper func(TIn) TOut) []TOut {
	out := make([]TOut, len(slice))
	for i, item := range slice {
		out[i] = mapper(item)
	}
	return out
}
