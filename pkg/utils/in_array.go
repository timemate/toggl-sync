package utils

func InArray[T string | int | int32 | int64 | float32 | float64](
	collection []T,
	value T,
) bool {
	for _, v := range collection {
		if v == value {
			return true
		}
	}

	return false
}
