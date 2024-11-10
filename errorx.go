package errorx

import (
	"errors"
	"fmt"
	"runtime"
)

type Error struct {
	Code     int64  `json:"code"`     // 错误代码
	Message  string `json:"message"`  // 错误消息
	Func     string `json:"func"`     // 函数名
	Position string `json:"position"` // 错误发生的位置
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Caller(skip ...int) *Error {
	skip = append(skip, 1)
	pc, file, line, ok := runtime.Caller(skip[0])
	if ok {
		e.Func = runtime.FuncForPC(pc).Name()
		e.Position = fmt.Sprintf("%s:%d", file, line)
	}
	return e
}

func New(message string, code ...int64) *Error {
	code = append(code, UnknownErr)
	return &Error{Code: code[0], Message: message}
}

func If(err error) error { // 返回类型约束不能为*Error，会导致返回值变成type为*Error的nil
	if err != nil {
		return New(err.Error()).Caller(2)
	}
	return nil
}

func Must(err error) *Error {
	return New(err.Error()).Caller(2)
}

func Wrap(err error) *Error {
	if err == nil {
		return nil
	}
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return New(err.Error()).Caller(2)
}

func Unknown(message string) *Error {
	return (&Error{Code: UnknownErr, Message: message}).Caller(2)
}

func Fail(message string) *Error {
	return (&Error{Code: OperateFail, Message: message}).Caller(2)
}
