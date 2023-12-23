package helper

func FindBy[T any](slice []T, predicate func(T) bool) (*T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return &item, true
		}
	}
	return nil, false
}

func Unshift[T any](slice []T, item T) []T {
	return append([]T{item}, slice...)
}
