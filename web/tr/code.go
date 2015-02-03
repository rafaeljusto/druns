package tr

const (
	CodeAuthenticationError Code = "authentication-error"
	CodeSessionExpired      Code = "session-expired"
	CodeInvalidDate         Code = "invalid-date"
	CodeInvalidEmail        Code = "invalid-email"
)

type Code string
