package model

type WrongEmailOrPassword struct{}

func (m *WrongEmailOrPassword) Error() string {
	return "wrong email or password"
}
