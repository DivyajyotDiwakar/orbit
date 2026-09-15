package main

import (
	"fmt"
	"net/http"
)

func AdminFlow() {

	fmt.Println("========== ADMIN FLOW ==========")

	GetPendingVerifications()

	GetVerificationDetail()

	ApproveVerification()

	fmt.Println()
	fmt.Println("ADMIN FLOW COMPLETED")
	fmt.Println()
}

func GetPendingVerifications() {

	body, status, err := adminClient.Get(
		"/api/admin/verifications/pending",
	)

	Assert(err)

	AssertStatus(http.StatusOK, status)

	resp := Decode[PendingVerificationResponse](body)

	AssertTrue(
		len(resp.Verifications) > 0,
		"Pending Verification Exists",
		"No pending verification found",
	)

	VerificationID = resp.Verifications[0].ID

	Pass("Fetch Pending Verifications")
}

func GetVerificationDetail() {

	body, status, err := adminClient.Get(
		"/api/admin/verifications/" + VerificationID,
	)

	Assert(err)

	AssertStatus(http.StatusOK, status)

	resp := Decode[VerificationDetailResponse](body)

	AssertTrue(
		resp.Verification.ID == VerificationID,
		"Verification Detail",
		"Invalid verification returned",
	)
}

func ApproveVerification() {

	body, status, err := adminClient.Post(
		"/api/admin/verifications/"+VerificationID+"/approve",
		nil,
	)

	Assert(err)

	AssertStatus(http.StatusOK, status)

	resp := Decode[MessageResponse](body)

	AssertTrue(
		resp.Message == "seller verification approved",
		"Approve Verification",
		resp.Message,
	)
}

func RejectVerification(reason string) {

	req := RejectVerificationRequest{
		Reason: reason,
	}

	body, status, err := adminClient.Post(
		"/api/admin/verifications/"+VerificationID+"/reject",
		req,
	)

	Assert(err)

	AssertStatus(http.StatusOK, status)

	resp := Decode[MessageResponse](body)

	AssertTrue(
		resp.Message == "seller verification rejected",
		"Reject Verification",
		resp.Message,
	)
}
