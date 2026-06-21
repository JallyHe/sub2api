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

type ClientAuthHandler struct {
	authService *service.AuthService
}

func NewClientAuthHandler(authService *service.AuthService) *ClientAuthHandler {
	return &ClientAuthHandler{authService: authService}
}

func (h *ClientAuthHandler) ShowLoginPage(c *gin.Context) {
	callback := c.DefaultQuery("callback", "storyclaw://auth")
	serverURL := detectServerURL(c)
	html := buildLoginPageHTML(callback, serverURL, "")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

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
	redirectURI := buildDeepLinkURI(callback, token, serverURL)
	c.Redirect(http.StatusFound, redirectURI)
}

func (h *ClientAuthHandler) OAuthSuccess(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(oauthSuccessHTML))
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
	oauthRedirect := fmt.Sprintf("/client-auth/oauth/success?callback=%s&server=%s",
		url.QueryEscape(callback), url.QueryEscape(serverURL))
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
      background: #0f0f13; color: #e2e2e7;
      min-height: 100vh; display: flex; align-items: center; justify-content: center;
    }
    .card {
      background: #1a1a24; border: 1px solid #2e2e3e; border-radius: 16px;
      padding: 40px; width: 100%%%%; max-width: 400px;
      box-shadow: 0 24px 60px rgba(0,0,0,.5);
    }
    .logo { display: flex; align-items: center; gap: 10px; margin-bottom: 28px; }
    .logo-icon {
      width: 36px; height: 36px;
      background: linear-gradient(135deg, #7c6aff, #a855f7);
      border-radius: 8px; display: flex; align-items: center; justify-content: center;
      font-size: 18px; font-weight: 800; color: #fff;
    }
    .logo-name { font-size: 18px; font-weight: 700; color: #fff; }
    h2 { font-size: 20px; font-weight: 600; margin-bottom: 6px; }
    p.sub { color: #888; font-size: 13px; margin-bottom: 24px; line-height: 1.5; }
    .field { margin-bottom: 16px; }
    label { display: block; font-size: 12px; color: #aaa; margin-bottom: 6px; font-weight: 500; }
    input {
      width: 100%%%%; padding: 11px 14px; background: #12121a; border: 1px solid #2e2e3e;
      border-radius: 8px; color: #e2e2e7; font-size: 14px; outline: none; transition: border-color .15s;
    }
    input:focus { border-color: #7c6aff; }
    .error {
      background: #3d1515; border: 1px solid #7c2626; color: #ff8080;
      border-radius: 8px; padding: 10px 14px; font-size: 13px; margin-bottom: 16px;
    }
    button.primary {
      width: 100%%%%; padding: 12px;
      background: linear-gradient(135deg, #7c6aff, #a855f7);
      border: none; border-radius: 8px; color: #fff; font-size: 14px; font-weight: 600;
      cursor: pointer; transition: opacity .15s;
    }
    button.primary:hover { opacity: .9; }
    .divider {
      display: flex; align-items: center; gap: 12px;
      margin: 24px 0; color: #555; font-size: 12px;
    }
    .divider::before, .divider::after { content: ''; flex: 1; height: 1px; background: #2e2e3e; }
    .oauth-row { display: flex; flex-direction: column; gap: 10px; }
    .oauth-btn {
      display: flex; align-items: center; justify-content: center; gap: 10px;
      width: 100%%%%; padding: 11px; border: 1px solid #2e2e3e; border-radius: 8px;
      background: transparent; color: #e2e2e7; font-size: 13px; font-weight: 500;
      cursor: pointer; transition: all .15s; text-decoration: none;
    }
    .oauth-btn:hover { border-color: #555; background: #1e1e2e; }
    .oauth-btn svg { width: 18px; height: 18px; flex-shrink: 0; }
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
      <button type="submit" class="primary">授权并返回 StoryClaw</button>
    </form>
    <div class="divider">或通过第三方登录</div>
    <div class="oauth-row">
      <a class="oauth-btn" href="/api/v1/auth/oauth/github/start?redirect=%s">
        <svg viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
        GitHub
      </a>
      <a class="oauth-btn" href="/api/v1/auth/oauth/google/start?redirect=%s">
        <svg viewBox="0 0 24 24"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 01-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4"/><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/></svg>
        Google
      </a>
    </div>
    <p class="hint">登录完成后页面会自动跳转回 StoryClaw 应用</p>
  </div>
</body>
</html>`, errHTML, url.QueryEscape(callback), serverURL, url.QueryEscape(oauthRedirect), url.QueryEscape(oauthRedirect))
}

const oauthSuccessHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>StoryClaw — 授权成功</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
      background: #0f0f13; color: #e2e2e7;
      min-height: 100vh; display: flex; align-items: center; justify-content: center;
    }
    .card {
      background: #1a1a24; border: 1px solid #2e2e3e; border-radius: 16px;
      padding: 40px; width: 100%%; max-width: 400px;
      box-shadow: 0 24px 60px rgba(0,0,0,.5); text-align: center;
    }
    .icon { font-size: 48px; margin-bottom: 16px; }
    h2 { font-size: 20px; font-weight: 600; margin-bottom: 8px; color: #fff; }
    p { color: #888; font-size: 13px; line-height: 1.5; }
    .error-msg { color: #ff8080; margin-top: 12px; }
  </style>
</head>
<body>
  <div class="card">
    <div class="icon">&#x2705;</div>
    <h2>授权成功</h2>
    <p id="status">正在返回 StoryClaw…</p>
    <p id="error" class="error-msg" style="display:none"></p>
  </div>
  <script>
    (function() {
      var params = new URLSearchParams(window.location.search);
      var hash = new URLSearchParams(window.location.hash.substring(1));
      var token = params.get('token') || hash.get('access_token') || '';
      var callback = params.get('callback') || 'storyclaw://auth';
      var server = params.get('server') || '';
      var statusEl = document.getElementById('status');
      var errorEl = document.getElementById('error');
      if (!token) {
        statusEl.textContent = '未获取到授权信息';
        errorEl.style.display = 'block';
        errorEl.textContent = '请关闭此页面，在 StoryClaw 中重新登录。';
        return;
      }
      var deepLink = callback + '?token=' + encodeURIComponent(token) + '&server=' + encodeURIComponent(server) + '&ts=' + Math.floor(Date.now() / 1000);
      window.location.href = deepLink;
      setTimeout(function() {
        statusEl.textContent = '如未自动跳转，请手动返回 StoryClaw 应用。';
      }, 2000);
    })();
  </script>
</body>
</html>`
