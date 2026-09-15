package verify

import (
	"Orbit/internal/models"
	"Orbit/internal/repositories"
	"Orbit/internal/utils"
)

func VerifySellerService(claims *utils.Claims, req *RequestBody) error {
	ObjectifiedId, err := utils.GetObjectFiedIdFromString(claims.ID)
	if err != nil {
		return err
	}
	var verificationReq = &models.VerificationRequest{
		SellerID:          ObjectifiedId,
		EmailId:           claims.EmailId,
		CompanyName:       req.CompanyName,
		GSTIN:             req.GSTIN,
		AccountHolderName: req.AccountHolderName,
		AccountNumber:     req.AccountNumber,
		IFSCCode:          req.IFSCCode,
		BankName:          req.BankName,
		Status:            "pending",
	}

	err = repositories.RegisterSellerVerificationRequest(verificationReq)
	if err != nil {
		return err
	}

	return nil
}
