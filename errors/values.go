package errors

var (
	UnableToResolve    = New("unable to resolve")
	UnexpectedResponse = New("unexpected response")
)

func WrongKeyType(expected, actual interface{}) error {
	return Errorf("wrong type: wanted %T, got %T", expected, actual)
}
