package login

type User struct {
	Username string
	Response string
	Nonce    string
}

type AuthInfo struct {
	Nonce string
}
