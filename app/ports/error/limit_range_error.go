package custom_error

type LimitRangeError struct {
	Limit int
}

func (e *LimitRangeError) Error() string {
	return "Limit must be between 1 and 100"
}
