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

func Concat[T any](slices ...[]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}
