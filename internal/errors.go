package internal

import "fmt"

const (
	ErrUnknownTable   = ConstError("unknown table")
	ErrRecordNotFound = ConstError("record not found")
)

type ConstError string

func (err ConstError) Error() string {
	return string(err)
}

type ErrFieldInvalidType struct {
	Item string
}

func (e ErrFieldInvalidType) Error() string {
	return fmt.Sprintf("field %s have invalid type", e.Item)
}
