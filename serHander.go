package main

import (
	"context"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/llaoj/oauth2nsso/config"
	"github.com/llaoj/oauth2nsso/model"
	"github.com/llaoj/oauth2nsso/pkg/session"
	"log"
	"net/http"
)

func passwordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	var user model.User
	userID, err = user.Authentication(context.Background(), username, password)
	return
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	v, _ := session.Get(r, "LoggedInUserID")
	if v == nil {
		if r.Form == nil {
			r.ParseForm()
		}
		session.Set(w, r, "RequestForm", r.Form)

		// 登录页面
		// 最终会把userId写进session(LoggedInUserID)
		// 再跳回来
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)

		return
	}
	userID = v.(string)

	// 不记住用户
	// store.Delete("LoggedInUserID")
	// store.Save()

	return
}

// 场景:在登录页面勾选所要访问的资源范围
// 根据client注册的scope,过滤表单中非法scope
// HandleAuthorizeRequest中调用
// set scope for the access token
func authorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	if r.Form == nil {
		r.ParseForm()
	}
	s := config.ScopeFilter(r.Form.Get("client_id"), r.Form.Get("scope"))
	if s == nil {
		err = errors.New("无效的权限范围")
		return
	}
	scope = config.ScopeJoin(s)
	return
}

func internalErrorHandler(err error) (re *errors.Response) {
	log.Println("Internal Error:", err.Error())
	return
}

func responseErrorHandler(re *errors.Response) {
	log.Println("Response Error:", re.Error.Error())
}
