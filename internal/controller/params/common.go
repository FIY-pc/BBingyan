package params

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
