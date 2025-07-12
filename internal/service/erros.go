package service

type ErrDuplicateEmail struct {
	Email string
}

func (e *ErrDuplicateEmail) Error() string {
	return "email already exits: " + e.Email
}
