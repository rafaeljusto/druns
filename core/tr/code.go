package tr

const (
	CodeAuthenticationError Code = "authentication-error"
	CodeSessionExpired      Code = "session-expired" // TODO: Remove this from the core?
	CodeInvalidDate         Code = "invalid-date"
	CodeInvalidEmail        Code = "invalid-email"
)

type Code string
