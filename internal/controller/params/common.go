package params

import "github.com/labstack/echo/v4"

type CommonErrorResp struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Error string `json:"error"`
}

type Common200Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func CommonErrorGenerate(c echo.Context, code int, msg string, err error) error {
	return c.JSON(code, CommonErrorResp{
		Code:  code,
		Msg:   msg,
		Error: err.Error(),
	})
}
