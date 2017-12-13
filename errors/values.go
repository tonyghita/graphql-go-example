package errors

var (
	UnableToResolve = New("unable to resolve")
)

func WrongType(expected, actual interface{}) error {
	return Errorf("wrong type: wanted %T, got %T", expected, actual)
}
