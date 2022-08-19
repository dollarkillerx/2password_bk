package request

type UserLogin struct {
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`

	Account  string `json:"account" binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserRegistry struct {
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`

	Account  string `json:"account" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password"  binding:"required"`
}
