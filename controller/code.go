package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
)

var codeMagMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "param Wrong",
	CodeUserExist:       "user existed",
	CodeUserNotExist:    "user not existed",
	CodeInvalidPassword: "username or password wrong",
	CodeServerBusy:      "server busy",
}

func (c ResCode) Msg() string {
	msg, ok := codeMagMap[c]
	if !ok {
		msg = codeMagMap[CodeServerBusy]
	}
	return msg
}
