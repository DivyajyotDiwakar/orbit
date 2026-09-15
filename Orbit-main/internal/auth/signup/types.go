package signup

type RegisterRequest struct {
	FirstName string `json:"firstName" binding:"required,min=3,max=20"`
	LastName  string `json:"lastName"  binding:"omitempty,min=3,max=20"`
	EmailId   string `json:"emailId"   binding:"required,email"`
	Age       int    `json:"age"       binding:"omitempty,gte=6,lte=80"`
	Password  string `json:"password"  binding:"required"`
	Role      string `json:"role"      binding:"required,oneof=buyer seller"`
}

type RegisterResponse struct {
	ID        string `json:"_id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName,omitempty"`
	EmailId   string `json:"emailId"`
	Age       int    `json:"age,omitempty"`
	Role      string `json:"role"`
}
