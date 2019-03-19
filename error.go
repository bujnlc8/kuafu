package kuafu

import "fmt"

type BaseError struct {
	Code int
	Msg string
}

func (err *BaseError)Error() string {
	return fmt.Sprintf("%d:%s", err.Code, err.Msg)
}


var RouterNotFound = &BaseError{Code: 4041, Msg:"router not found"}