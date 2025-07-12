package error

type DomainError interface {
	error
	Code() string
	Type() string
}
