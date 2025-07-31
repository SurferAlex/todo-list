package entity

type User struct {
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type RegistrationRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=100,containsany=abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789"`
}
