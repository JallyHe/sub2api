package admin

import (
	"net/http"
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminCreditHandler handles admin credit management endpoints.
type AdminCreditHandler struct {
	creditSvc *service.CreditService
}

// NewAdminCreditHandler constructs an AdminCreditHandler.
func NewAdminCreditHandler(creditSvc *service.CreditService) *AdminCreditHandler {
	return &AdminCreditHandler{creditSvc: creditSvc}
}

type grantCreditsRequest struct {
	UserID  int64  `json:"user_id" binding:"required"`
	Credits int64  `json:"credits" binding:"required,min=1"`
	Notes   string `json:"notes"`
}

// ListModelRates handles GET /api/v1/admin/credits/model-rates
func (h *AdminCreditHandler) ListModelRates(c *gin.Context) {
	rates, err := h.creditSvc.GetAllModelRates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rates})
}

// GetUserBalance handles GET /api/v1/admin/credits/users/:id/balance
func (h *AdminCreditHandler) GetUserBalance(c *gin.Context) {
	uid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	bal, err := h.creditSvc.GetBalance(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bal)
}

// GrantCredits handles POST /api/v1/admin/credits/grant
func (h *AdminCreditHandler) GrantCredits(c *gin.Context) {
	var req grantCreditsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.creditSvc.AdminGrantCredits(c.Request.Context(), req.UserID, req.Credits, req.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
