package chefs

import "fmt"

type ErrChefAlreadyExists struct {
	email string
}

func (e *ErrChefAlreadyExists) Error() string {
	return fmt.Sprintf("another chef is using %s as their email", e.email)
}

type ErrInvalidEmail struct {
	email string
}

func (e *ErrInvalidEmail) Error() string {
	return fmt.Sprintf("%s is an invalid email", e.email)
}

type ErrInternal struct {
	e error
}

func (e *ErrInternal) Error() string {
	return fmt.Sprintf("internal error: %v", e.e)
}

type ErrInvalidActivationMethod struct {
	reason string
}

func (e *ErrInvalidActivationMethod) Error() string {
	return fmt.Sprintf("unable to activate account, reason: %s", e.reason)
}
