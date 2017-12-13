package errors

var (
	UnableToResolve    = New("unable to resolve")
	UnexpectedResponse = New("unexpected response")
)

func WrongKeyType(expected string, actual interface{}) error {
	return Errorf("wrong type: wanted %s, got %T", expected, actual)
}
