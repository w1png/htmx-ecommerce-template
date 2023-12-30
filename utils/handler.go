package utils

func GetOffsetAndLimit(page, limit int) (int, int) {
	return limit * (page - 1), limit
}
