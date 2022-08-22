package request

type AuthJWT struct {
	Account string `json:"account"`
	Name    string `json:"name"`
}

type LoginExpiration struct {
	Expiration int64 `json:"expiration"`
}
