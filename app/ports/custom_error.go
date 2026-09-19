package ports

type Error interface {
	Code() int
	Error() string
}
