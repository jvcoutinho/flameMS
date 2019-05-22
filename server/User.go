package server

type User struct {
	identifier  string
	handler     *RequestHandler
	isConnected bool
}

func newUser(identifier string, handler *RequestHandler) *User {
	return &User{identifier, handler, true}
}

func (user *User) setHandler(handler *RequestHandler) {
	user.handler = handler
}

func (user *User) setConnected(isConnected bool) {
	user.isConnected = isConnected
}
