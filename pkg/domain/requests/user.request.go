package requests

// UserLogin request
type UserLogin struct {
	Username string `json:"username" xml:"username" form:"username" validate:"required,email"`
	Password string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
}

// UserByID request
type UserByID struct {
	ID string `json:"id" xml:"id" form:"id" validate:"required,uuid"`
}

// UserCreation request to create a user
type UserCreation struct {
	Username  string `json:"username" xml:"username" form:"username" validate:"required,email"`
	Password  string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
	Lastname  string `json:"lastname" xml:"lastname" form:"lastname" validate:"required"`
	Firstname string `json:"firstname" xml:"firstname" form:"firstname" validate:"required"`
}

// UserUpdate request to update a user
type UserUpdate struct {
	ID        string `json:"id" xml:"id" form:"id" validate:"required,uuid"`
	Username  string `json:"username" xml:"username" form:"username" validate:"required,email"`
	Password  string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
	Lastname  string `json:"lastname" xml:"lastname" form:"lastname" validate:"required"`
	Firstname string `json:"firstname" xml:"firstname" form:"firstname" validate:"required"`
}

// UserPasswordUpdate request to update a user password
type UserPasswordUpdate struct {
	Token    string `json:"token" xml:"token" form:"token" validate:"required"`
	Password string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
}

// UserForgotPassword request to reset user password
type UserForgotPassword struct {
	Email string `json:"email" xml:"email" form:"email" validate:"required,email"`
}
