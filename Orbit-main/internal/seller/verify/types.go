package verify

type RequestBody struct {
	CompanyName       string `json:"companyName" binding:"required"`
	GSTIN             string `json:"gstin" binding:"required"`
	AccountHolderName string `json:"accountHolderName" binding:"required"`
	AccountNumber     string `json:"accountNumber" binding:"required"`
	IFSCCode          string `json:"ifscCode" binding:"required"`
	BankName          string `json:"bankName" binding:"required"`
}
