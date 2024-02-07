package requests

// UserLogin request
type UserLogin struct {
	Username string `json:"username" xml:"username" form:"username" validate:"required,email"`
	Password string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
}

// UserEdit request to create or update a user
type UserEdit struct {
	Username  string `json:"username" xml:"username" form:"username" validate:"required,email"`
	Password  string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
	Lastname  string `json:"lastname" xml:"lastname" form:"lastname" validate:"required"`
	Firstname string `json:"firstname" xml:"firstname" form:"firstname" validate:"required"`
}

// UserByID request
type UserByID struct {
	ID string `json:"id" xml:"id" form:"id" validate:"required,uuid"`
}
