package tr

const (
	CodeAuthenticationError Code = "authentication-error"
	CodeInvalidDate         Code = "invalid-date"
	CodeInvalidDuration     Code = "invalid-duration"
	CodeInvalidEmail        Code = "invalid-email"
	CodeInvalidGroupType    Code = "invalid-group-type"
	CodeInvalidTime         Code = "invalid-time"
	CodeInvalidWeekday      Code = "invalid-weekday"
	CodeSessionExpired      Code = "session-expired" // TODO: Remove this from the core?
)

type Code string
