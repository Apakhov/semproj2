package psql

func fixLimit(limit, max int64) int64 {
	if limit <= 0 {
		return 1
	}
	if limit > max {
		return max
	}
	return limit
}
