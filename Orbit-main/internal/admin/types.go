package admin

type RejectRequestBody struct {
	Reason string `json:"reason" binding:"required,min=5,max=500"`
}
