package main

import (
	"encoding/json"
	"github.com/llaoj/oauth2nsso/config"
	"github.com/llaoj/oauth2nsso/model"
	"github.com/llaoj/oauth2nsso/pkg/session"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

// 首先进入执行
func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	var form url.Values
	v, _ := session.Get(r, "RequestForm")
	if v != nil {
		if r.Form.Get("client_id") == "" {
			form = v.(url.Values)
		}
	}
	r.Form = form

	if err := session.Delete(w, r, "RequestForm"); err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := srv.HandleAuthorizeRequest(w, r); err != nil {
		errorHandler(w, err.Error(), http.StatusBadRequest)
		return
	}
}

type TplData struct {
	Client config.OAuth2Client
	// 用户申请的合规scope
	Scope []config.Scope
	Error string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	form, _ := session.Get(r, "RequestForm")
	if form == nil {
		errorHandler(w, "无效的请求", http.StatusInternalServerError)
		return
	}

	clientID := form.(url.Values).Get("client_id")
	scope := form.(url.Values).Get("scope")

	// 页面数据
	data := TplData{
		Client: config.GetOAuth2Client(clientID),
		Scope:  config.ScopeFilter(clientID, scope),
	}
	if data.Scope == nil {
		errorHandler(w, "无效的权限范围", http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		var userID string
		var err error

		if r.Form == nil {
			err = r.ParseForm()
			if err != nil {
				errorHandler(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// 方式1:账号密码验证
		if r.Form.Get("type") == "password" {
			var user model.User
			userID, err = user.Authentication(r.Context(), r.Form.Get("username"), r.Form.Get("password"))
			if err != nil {
				data.Error = err.Error()
				t, _ := template.ParseFiles("tpl/login.html")
				t.Execute(w, data)
				return
			}
		}

		// 方式2:扫码验证
		// 方式3:手机验证码验证
		// 方式N:...

		err = session.Set(w, r, "LoggedInUserID", userID)
		if err != nil {
			errorHandler(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", "/authorize")
		w.WriteHeader(http.StatusFound)
		return
	}

	t, _ := template.ParseFiles("tpl/login.html")
	t.Execute(w, data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			errorHandler(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 检查redirect_uri参数
	redirectURI := r.Form.Get("redirect_uri")
	if redirectURI == "" {
		errorHandler(w, "参数不能为空(redirect_uri)", http.StatusBadRequest)
		return
	}
	if _, err := url.Parse(redirectURI); err != nil {
		errorHandler(w, "参数无效(redirect_uri)", http.StatusBadRequest)
		return
	}

	// 删除公共回话
	if err := session.Delete(w, r, "LoggedInUserID"); err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", redirectURI)
	w.WriteHeader(http.StatusFound)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	err := srv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	token, err := srv.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cli, err := mgr.GetClient(r.Context(), token.GetClientID())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"user_id":    token.GetUserID(),
		"client_id":  token.GetClientID(),
		"scope":      token.GetScope(),
		"domain":     cli.GetDomain(),
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(data)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	errorHandler(w, "无效的地址", http.StatusNotFound)
	return
}

// 错误显示页面
// 以网页的形式展示大于400的错误
func errorHandler(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	if status >= 400 {
		t, _ := template.ParseFiles("tpl/error.html")
		body := struct {
			Status  int
			Message string
		}{Status: status, Message: message}
		t.Execute(w, body)
	}
}
