package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// ClientAuthHandler handles the browser-based OAuth-style flow for StoryClaw desktop clients.
// Flow:
//  1. StoryClaw opens browser: GET /client-auth?callback=storyclaw://auth&server=https://this-server.com
//  2. User fills login form → POST /client-auth/login
//  3. On success, server redirects to storyclaw://auth?token={jwt}&server={server}
type ClientAuthHandler struct {
	authService *service.AuthService
}

// NewClientAuthHandler constructs a ClientAuthHandler.
func NewClientAuthHandler(authService *service.AuthService) *ClientAuthHandler {
	return &ClientAuthHandler{authService: authService}
}

// ShowLoginPage renders the standalone login page for desktop clients.
// GET /client-auth?callback=storyclaw://auth
func (h *ClientAuthHandler) ShowLoginPage(c *gin.Context) {
	callback := c.DefaultQuery("callback", "storyclaw://auth")
	serverURL := detectServerURL(c)

	html := buildLoginPageHTML(callback, serverURL, "")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// HandleLogin processes the login form submission.
// POST /client-auth/login
func (h *ClientAuthHandler) HandleLogin(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")
	callback := c.DefaultQuery("callback", "storyclaw://auth")
	serverURL := c.PostForm("server_url")
	if serverURL == "" {
		serverURL = detectServerURL(c)
	}

	if email == "" || password == "" {
		html := buildLoginPageHTML(callback, serverURL, "请填写邮箱和密码")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
		return
	}

	token, _, err := h.authService.Login(c.Request.Context(), email, password)
	if err != nil {
		html := buildLoginPageHTML(callback, serverURL, "邮箱或密码错误，请重试")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
		return
	}

	// Build the deep link redirect URI
	redirectURI := buildDeepLinkURI(callback, token, serverURL)
	c.Redirect(http.StatusFound, redirectURI)
}

func detectServerURL(c *gin.Context) string {
	scheme := "https"
	if c.Request.TLS == nil {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}

func buildDeepLinkURI(callback, token, serverURL string) string {
	u, err := url.Parse(callback)
	if err != nil {
		u = &url.URL{Scheme: "storyclaw", Host: "auth"}
	}
	q := url.Values{}
	q.Set("token", token)
	q.Set("server", serverURL)
	q.Set("ts", fmt.Sprintf("%d", time.Now().Unix()))
	u.RawQuery = q.Encode()
	return u.String()
}

func buildLoginPageHTML(callback, serverURL, errMsg string) string {
	errHTML := ""
	if errMsg != "" {
		errHTML = fmt.Sprintf(`<div class="error">%s</div>`, errMsg)
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>StoryClaw — 登录授权</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
      background: #0f0f13;
      color: #e2e2e7;
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .card {
      background: #1a1a24;
      border: 1px solid #2e2e3e;
      border-radius: 16px;
      padding: 40px;
      width: 100%%;
      max-width: 400px;
      box-shadow: 0 24px 60px rgba(0,0,0,.5);
    }
    .logo {
      display: flex;
      align-items: center;
      gap: 10px;
      margin-bottom: 28px;
    }
    .logo-icon {
      width: 36px; height: 36px;
      background: linear-gradient(135deg, #7c6aff, #a855f7);
      border-radius: 8px;
      display: flex; align-items: center; justify-content: center;
      font-size: 18px; font-weight: 800; color: #fff;
    }
    .logo-name { font-size: 18px; font-weight: 700; color: #fff; }
    h2 { font-size: 20px; font-weight: 600; margin-bottom: 6px; }
    p.sub { color: #888; font-size: 13px; margin-bottom: 24px; line-height: 1.5; }
    .field { margin-bottom: 16px; }
    label { display: block; font-size: 12px; color: #aaa; margin-bottom: 6px; font-weight: 500; }
    input {
      width: 100%%; padding: 11px 14px;
      background: #12121a; border: 1px solid #2e2e3e;
      border-radius: 8px; color: #e2e2e7; font-size: 14px;
      outline: none; transition: border-color .15s;
    }
    input:focus { border-color: #7c6aff; }
    .error {
      background: #3d1515; border: 1px solid #7c2626;
      color: #ff8080; border-radius: 8px;
      padding: 10px 14px; font-size: 13px; margin-bottom: 16px;
    }
    button {
      width: 100%%; padding: 12px;
      background: linear-gradient(135deg, #7c6aff, #a855f7);
      border: none; border-radius: 8px;
      color: #fff; font-size: 14px; font-weight: 600;
      cursor: pointer; transition: opacity .15s;
    }
    button:hover { opacity: .9; }
    .hint { text-align: center; color: #555; font-size: 12px; margin-top: 20px; }
  </style>
</head>
<body>
  <div class="card">
    <div class="logo">
      <div class="logo-icon">S</div>
      <span class="logo-name">StoryClaw</span>
    </div>
    <h2>登录授权</h2>
    <p class="sub">登录后，StoryClaw 将自动完成配置，无需手动填写任何参数。</p>
    %s
    <form method="POST" action="/client-auth/login?callback=%s">
      <input type="hidden" name="server_url" value="%s" />
      <div class="field">
        <label>邮箱</label>
        <input type="email" name="email" placeholder="your@email.com" required autofocus />
      </div>
      <div class="field">
        <label>密码</label>
        <input type="password" name="password" placeholder="••••••••" required />
      </div>
      <button type="submit">授权并返回 StoryClaw</button>
    </form>
    <p class="hint">登录完成后页面会自动跳转回 StoryClaw 应用</p>
  </div>
</body>
</html>`, errHTML, url.QueryEscape(callback), serverURL)
}
