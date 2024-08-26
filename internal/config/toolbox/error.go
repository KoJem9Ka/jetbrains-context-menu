package toolbox

type StateReadError struct {
	error
}

func (this StateReadError) Error() string {
	return this.error.Error()
}
