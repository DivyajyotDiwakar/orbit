package signin

type SignInRequest struct {
	EmailId  string `json:"emailId"   binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}

type SignInResponse struct {
	ID        string `json:"_id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName,omitempty"`
	EmailId   string `json:"emailId"`
	Age       int    `json:"age,omitempty"`
	Role      string `json:"role"`
}
