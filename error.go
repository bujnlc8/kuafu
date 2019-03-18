package kuafu

import "fmt"

type BaseError struct {
	Code int
	Msg string
}

func (err *BaseError)Error() string {
	return fmt.Sprintf("error code :%d, msg:%s", err.Code, err.Msg)
}


var RouterNotFound = &BaseError{Code: 1404, Msg:"router not found"}