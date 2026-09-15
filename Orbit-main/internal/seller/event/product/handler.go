package product

import (
	"Orbit/internal/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var productUploader utils.Uploader = &utils.Product{}

func extractClaim(c *gin.Context) (*utils.Claims, bool) {
	raw, ok := c.Get("UserFields")
	if !ok {
		return nil, false
	}
	claim, ok := raw.(*utils.Claims)
	return claim, ok
}

func mapServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrEventNotFound), errors.Is(err, ErrProductNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, ErrUnauthorized):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNoUpdateFields):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrInvalidSeller), errors.Is(err, ErrInvalidEvent):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func RegisterProductHandler(c *gin.Context) {
	claim, ok := extractClaim(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not permitted"})
		return
	}

	eventIdStr := c.Param("eventId")

	var req RegisterProductRequestBody
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product image is required"})
		return
	}

	imgURL, err := productUploader.UploadPhoto(c.Request.Context(), fileHeader, eventIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := RegisterProductService(c.Request.Context(), claim.ID, eventIdStr, req, imgURL); err != nil {
		_ = productUploader.DeletePhoto(c.Request.Context(), "product", eventIdStr)
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "product registered successfully"})
}

func GetAllEventProductsHandler(c *gin.Context) {
	eventIdStr := c.Param("eventId")

	products, err := GetAllEventProductsService(c.Request.Context(), eventIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func GetAnEventProductHandler(c *gin.Context) {
	productIdStr := c.Param("productId")

	product, err := GetAnEventProductService(c.Request.Context(), productIdStr)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func UpdateAnEventProductHandler(c *gin.Context) {
	claim, ok := extractClaim(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not permitted"})
		return
	}

	eventIdStr := c.Param("eventId")
	productIdStr := c.Param("productId")

	var req UpdateProductRequestBody
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newImageURL string
	fileHeader, err := c.FormFile("image")
	if err == nil {
		newImageURL, err = productUploader.UploadPhoto(c.Request.Context(), fileHeader, productIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if err := UpdateAnEventProductService(
		c.Request.Context(),
		claim.ID,
		eventIdStr,
		productIdStr,
		req,
		newImageURL,
	); err != nil {
		if newImageURL != "" {
			_ = productUploader.DeletePhoto(c.Request.Context(), "product", productIdStr)
		}
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product updated successfully"})
}

func DeleteAnEventProductHandler(c *gin.Context) {
	claim, ok := extractClaim(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not permitted"})
		return
	}

	eventIdStr := c.Param("eventId")
	productIdStr := c.Param("productId")

	if err := DeleteAnEventProductService(c.Request.Context(), claim.ID, eventIdStr, productIdStr); err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}
