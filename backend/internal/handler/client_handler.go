package handler

import (
	"fmt"
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	authService   *service.AuthService
	creditService *service.CreditService
	apiKeyService *service.APIKeyService
}

func NewClientHandler(
	authService *service.AuthService,
	creditService *service.CreditService,
	apiKeyService *service.APIKeyService,
) *ClientHandler {
	return &ClientHandler{
		authService:   authService,
		creditService: creditService,
		apiKeyService: apiKeyService,
	}
}

type clientLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type clientLoginResponse struct {
	Token string `json:"token"`
}

func (h *ClientHandler) Login(c *gin.Context) {
	var req clientLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, _, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
		return
	}
	c.JSON(http.StatusOK, clientLoginResponse{Token: token})
}

type clientModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type clientCredit struct {
	Balance   int64  `json:"balance"`
	ExpiresAt *int64 `json:"expires_at,omitempty"`
}

type clientModelsResponse struct {
	Endpoint string        `json:"endpoint"`
	APIKey   string        `json:"api_key"`
	Models   []clientModel `json:"models"`
	Credits  *clientCredit `json:"credits,omitempty"`
}

var defaultModels = []clientModel{
	{ID: "claude-opus-4-5", Name: "Claude Opus 4", Type: "claude"},
	{ID: "claude-sonnet-4-5", Name: "Claude Sonnet 4", Type: "claude"},
	{ID: "claude-haiku-4-5", Name: "Claude Haiku 4", Type: "claude"},
	{ID: "gpt-4o", Name: "GPT-4o", Type: "openai"},
	{ID: "gpt-4o-mini", Name: "GPT-4o Mini", Type: "openai"},
	{ID: "gemini-2.0-flash", Name: "Gemini 2.0 Flash", Type: "gemini"},
	{ID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro", Type: "gemini"},
}

func (h *ClientHandler) GetModels(c *gin.Context) {
	uid := mustUserID(c)

	scheme := "https"
	if c.Request.TLS == nil {
		scheme = "http"
	}
	endpoint := fmt.Sprintf("%s://%s/v1", scheme, c.Request.Host)

	keys, _, err := h.apiKeyService.List(c.Request.Context(), uid,
		pagination.PaginationParams{Page: 1, PageSize: 1},
		service.APIKeyListFilters{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list API keys"})
		return
	}

	var apiKey string
	if len(keys) > 0 {
		apiKey = keys[0].Key
	} else {
		newKey, createErr := h.apiKeyService.Create(c.Request.Context(), uid, service.CreateAPIKeyRequest{
			Name: "storyclaw-default",
		})
		if createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create API key"})
			return
		}
		apiKey = newKey.Key
	}

	var credit *clientCredit
	if bal, balErr := h.creditService.GetBalance(c.Request.Context(), uid); balErr == nil {
		credit = &clientCredit{Balance: bal.Balance}
		if bal.ExpiresAt != nil {
			t := bal.ExpiresAt.Unix()
			credit.ExpiresAt = &t
		}
	}

	c.JSON(http.StatusOK, clientModelsResponse{
		Endpoint: endpoint,
		APIKey:   apiKey,
		Models:   defaultModels,
		Credits:  credit,
	})
}

func (h *ClientHandler) GetCredits(c *gin.Context) {
	uid := mustUserID(c)
	bal, err := h.creditService.GetBalance(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get credits"})
		return
	}
	resp := &clientCredit{Balance: bal.Balance}
	if bal.ExpiresAt != nil {
		t := bal.ExpiresAt.Unix()
		resp.ExpiresAt = &t
	}
	c.JSON(http.StatusOK, resp)
}
