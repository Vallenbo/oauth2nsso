package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt/v5"

	"github.com/llaoj/oauth2nsso/config"
	"github.com/llaoj/oauth2nsso/pkg/session"
)

var srv *server.Server
var mgr *manage.Manager

func main() {
	config.Setup()
	// init db connection
	// configure db in app.yaml then uncomment
	// model.Setup()
	session.Setup()

	// manager config
	mgr = manage.NewDefaultManager()
	mgr.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    time.Hour * time.Duration(config.Get().OAuth2.AccessTokenExp),
		RefreshTokenExp:   time.Hour * 24 * 3,
		IsGenerateRefresh: true})
	// token store
	mgr.MustTokenStorage(store.NewMemoryTokenStore())
	// or use redis token store
	// mgr.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
	//     Addr: config.Get().Redis.Default.Addr,
	//     DB: config.Get().Redis.Default.DB,
	// }))

	// access token generate method: jwt
	mgr.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte(config.Get().OAuth2.JWTSignedKey), jwt.SigningMethodHS512))
	clientStore := store.NewClientStore()
	for _, v := range config.Get().OAuth2.Client {
		err := clientStore.Set(v.ID, &models.Client{
			ID:     v.ID,
			Secret: v.Secret,
			Domain: v.Domain,
		})
		if err != nil {
			log.Fatalln("config Get OAuth2 clientStore Set err : ", err)
			return
		}
	}
	mgr.MapClientStorage(clientStore)
	// config oauth2 server
	srv = server.NewServer(server.NewConfig(), mgr)
	srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	srv.SetAuthorizeScopeHandler(authorizeScopeHandler)
	srv.SetInternalErrorHandler(internalErrorHandler)
	srv.SetResponseErrorHandler(responseErrorHandler)

	// http server
	http.HandleFunc("/authorize", authorizeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/token", tokenHandler)
	http.HandleFunc("/verify", verifyHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", notFoundHandler)

	log.Println("Server is running at 9096 port.")
	log.Fatal(http.ListenAndServe(":9096", nil))
}
