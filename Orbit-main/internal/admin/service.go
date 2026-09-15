package admin

import (
	"Orbit/internal/models"
	"Orbit/internal/repositories"
	"Orbit/internal/utils"
	"context"
)

func GetPendingVerificationsService(ctx context.Context) ([]models.VerificationRequest, error) {
	return repositories.GetPendingVerifications(ctx)
}

func GetVerificationDetailService(ctx context.Context, id string) (*models.VerificationRequest, error) {
	return repositories.GetVerificationByID(ctx, id)
}

func ApproveVerificationService(ctx context.Context, verificationId string, adminId string) error {
	adminObjId, err := utils.GetObjectFiedIdFromString(adminId)
	if err != nil {
		return err
	}
	return repositories.ApproveVerification(ctx, verificationId, adminObjId)
}

func RejectVerificationService(ctx context.Context, verificationId string, adminId string, reason string) error {
	adminObjId, err := utils.GetObjectFiedIdFromString(adminId)
	if err != nil {
		return err
	}
	return repositories.RejectVerification(ctx, verificationId, adminObjId, reason)
}
