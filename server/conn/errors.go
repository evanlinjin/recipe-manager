package conn

import "fmt"

type ErrCloseMessage struct {
	code int
}

func (e *ErrCloseMessage) Error() string {
	return fmt.Sprintf("close message %v", e.code)
}

type ErrUnexpectedMessage struct {
	code int
}

func (e *ErrUnexpectedMessage) Error() string {
	return fmt.Sprintf("got unexpected message %v", e.code)
}
