package bean

const Success = "Success"
const UnexpectedError = "UnexpectedError"
const NotSignIn = "NotSignIn"

var CodeMessage = map[string]struct {
	Code    int
	Message string
}{
	Success: {1, "Success"},

	// -x for basic message
	UnexpectedError: {-1, "Unexpected error"},
	NotSignIn:       {-1, "User is not signed in"},
}
