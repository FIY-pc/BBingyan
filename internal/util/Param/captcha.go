package Param

// CaptchaKey 统一redis中验证码key的样式
func CaptchaKey(email string) string {
	return "captcha:" + email
}
