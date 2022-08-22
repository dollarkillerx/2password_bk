package request

type UserLogin struct {
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`

	Account string `json:"account" binding:"required"`
	Sign    string `json:"sign"  binding:"required"`
}

type UserRegistry struct {
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`

	Account string `json:"account" binding:"required"`

	PublicKey           string `json:"public_key" binding:"required"`
	EncryptedPrivateKey string `json:"encrypted_private_key" binding:"required"`
}
