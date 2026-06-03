package handler

import (
	"net/http"
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// CreditHandler handles user-facing credit endpoints.
type CreditHandler struct {
	creditService *service.CreditService
}

// NewCreditHandler constructs a CreditHandler.
func NewCreditHandler(creditService *service.CreditService) *CreditHandler {
	return &CreditHandler{creditService: creditService}
}

type creditBalanceResponse struct {
	Balance   int64  `json:"balance"`
	ExpiresAt *int64 `json:"expires_at,omitempty"` // Unix timestamp; omitted when nil (no expiry)
	PlanID    *int64 `json:"plan_id,omitempty"`
}

// GetBalance handles GET /api/v1/credits/balance
func (h *CreditHandler) GetBalance(c *gin.Context) {
	uid := mustUserID(c)
	bal, err := h.creditService.GetBalance(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}
	resp := creditBalanceResponse{Balance: bal.Balance, PlanID: bal.PlanID}
	if bal.ExpiresAt != nil {
		t := bal.ExpiresAt.Unix()
		resp.ExpiresAt = &t
	}
	c.JSON(http.StatusOK, resp)
}

type ledgerEntry struct {
	ID           int64   `json:"id"`
	Delta        int64   `json:"delta"`
	Reason       string  `json:"reason"`
	Model        *string `json:"model,omitempty"`
	BalanceAfter int64   `json:"balance_after"`
	CreatedAt    int64   `json:"created_at"` // Unix timestamp
}

// ListLedger handles GET /api/v1/credits/ledger?page=1&limit=20
func (h *CreditHandler) ListLedger(c *gin.Context) {
	uid := mustUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	entries, err := h.creditService.ListLedger(c.Request.Context(), uid, (page-1)*limit, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list ledger"})
		return
	}
	out := make([]ledgerEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, ledgerEntry{
			ID:           e.ID,
			Delta:        e.Delta,
			Reason:       e.Reason,
			Model:        e.Model,
			BalanceAfter: e.BalanceAfter,
			CreatedAt:    e.CreatedAt.Unix(),
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": out, "page": page, "limit": limit})
}

// mustUserID extracts the int64 userID set by JWTAuthMiddleware. Returns 0 if absent.
func mustUserID(c *gin.Context) int64 {
	v, _ := c.Get("userID")
	if uid, ok := v.(int64); ok {
		return uid
	}
	return 0
}
