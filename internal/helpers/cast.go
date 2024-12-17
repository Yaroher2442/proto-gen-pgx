package helpers

func ToAnyList[T any](items []T) []any {
	var res []any
	for _, item := range items {
		res = append(res, item)
	}
	return res
}
