package event

import (
	"Orbit/internal/models"
	"Orbit/internal/repositories"
	"Orbit/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var eventUploader utils.Uploader = &utils.Event{}

func CreateAnEventHandler(c *gin.Context) {
	rawClaim, ok := c.Get("UserFields")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not permitted"})
		return
	}
	sellerClaim, ok := rawClaim.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid session"})
		return
	}

	sellerObjId, err := utils.GetObjectFiedIdFromString(sellerClaim.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req RequestBody
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := c.FormFile("imageBanner")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "imageBanner file is required"})
		return
	}

	newEventId := bson.NewObjectID()

	imgURL, err := eventUploader.UploadPhoto(c.Request.Context(), fileHeader, newEventId.Hex())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := &models.Event{
		ID:          newEventId,
		SellerID:    sellerObjId,
		EventName:   req.Title,
		Description: req.Description,
		ScheduledAt: req.ScheduledAt,
		ImageBanner: imgURL,
		IsLive:      false,
	}

	if err := repositories.CreateAnEvent(c.Request.Context(), event); err != nil {
		_ = eventUploader.DeletePhoto(c.Request.Context(), "event-banner", newEventId.Hex())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "event created successfully",
		"eventId": newEventId.Hex(),
	})
}

func GetAllEventsHandler(c *gin.Context) {
	rawClaim, ok := c.Get("UserFields")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not permitted"})
		return
	}
	sellerClaim, ok := rawClaim.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid session"})
		return
	}

	events, err := repositories.GetAllEvents(c.Request.Context(), sellerClaim.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func GetAnEventHandler(c *gin.Context) {
	eventId := c.Param("id")

	rawClaim, ok := c.Get("UserFields")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not permitted"})
		return
	}
	sellerClaim, ok := rawClaim.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid session"})
		return
	}

	event, err := repositories.GetAnEvent(c.Request.Context(), eventId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	if event.SellerID.Hex() != sellerClaim.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

func UpdateAnEventHandler(c *gin.Context) {
	rawClaim, ok := c.Get("UserFields")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not permitted"})
		return
	}
	sellerClaim, ok := rawClaim.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid session"})
		return
	}

	eventId := c.Param("id")

	event, err := repositories.GetAnEvent(c.Request.Context(), eventId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	if event.SellerID.Hex() != sellerClaim.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var req UpdateEventRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := bson.M{}
	if req.Title != nil {
		updateData["eventName"] = *req.Title
	}
	if req.Description != nil {
		updateData["description"] = *req.Description
	}
	if req.ScheduledAt != nil {
		updateData["scheduledAt"] = *req.ScheduledAt
	}
	if req.ImageBanner != nil {
		updateData["imageBanner"] = *req.ImageBanner
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided for update"})
		return
	}

	if err := repositories.UpdateAnEvent(c.Request.Context(), eventId, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event updated successfully"})
}

func DeleteAnEventHandler(c *gin.Context) {
	rawClaim, ok := c.Get("UserFields")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not permitted"})
		return
	}
	sellerClaim, ok := rawClaim.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid session"})
		return
	}

	eventId := c.Param("id")

	event, err := repositories.GetAnEvent(c.Request.Context(), eventId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	if event.SellerID.Hex() != sellerClaim.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := repositories.DeleteAnEvent(c.Request.Context(), eventId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})
}
