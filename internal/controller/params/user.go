package params

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

type RegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}
