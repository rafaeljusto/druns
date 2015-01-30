package tr

const (
	CodeSessionExpired      Code = "session-expired"
	CodeInvalidEmail        Code = "invalid-email"
	CodeAuthenticationError Code = "authentication-error"
)

type Code string
