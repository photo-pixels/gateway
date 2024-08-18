package model

func (e Error) Error() string {
	return string(e)
}
