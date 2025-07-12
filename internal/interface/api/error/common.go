package error

func GetErrorCode(err error) string {
	if domainErr, ok := err.(interface{ Code() string }); ok {
		return domainErr.Code()
	}
	return "UNKNOWN_ERROR"
}
