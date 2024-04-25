package dto

func ToArr[T struct{}, F interface{}](m []F, fn func(F) T) []T {
	arr := make([]T, len(m))
	for i, val := range m {
		arr[i] = fn(val)
	}
	return arr
}