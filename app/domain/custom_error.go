package domain

type Error interface {
	Code() int
	Error() string
}
