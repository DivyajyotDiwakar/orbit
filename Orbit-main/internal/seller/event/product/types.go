package product

type RegisterProductRequestBody struct {
	Title       string  `form:"title"       binding:"required,min=2,max=100"`
	Description string  `form:"description" binding:"required,min=10,max=500"`
	Price       float64 `form:"price"       binding:"required,gt=0"`
	Frequency   int     `form:"frequency"   binding:"required,gt=0"`
}

type UpdateProductRequestBody struct {
	Title       *string  `form:"title"       binding:"omitempty,min=2,max=100"`
	Description *string  `form:"description" binding:"omitempty,min=10,max=500"`
	Price       *float64 `form:"price"       binding:"omitempty,gt=0"`
	Frequency   *int     `form:"frequency"   binding:"omitempty,gt=0"`
}
