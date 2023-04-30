package main

import (
	"MyOauth2/config"
	"MyOauth2/config/global"
	"MyOauth2/dao"

	"MyOauth2/session"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var srv *server.Server
var mgr *manage.Manager

func main() {
	config.ViperSetup()
	config.MysqlDbSetup()
	session.Setup()
	mgr = manage.NewDefaultManager()
	mgr.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    time.Hour * 24 * 3 * time.Duration(global.Config.Oauth2.AccessTokenExp),
		RefreshTokenExp:   time.Hour * 24 * 3 * time.Duration(global.Config.Oauth2.RefreshTokenExp),
		IsGenerateRefresh: global.Config.Oauth2.IsGenerateRefresh})
	// token store
	mgr.MustTokenStorage(store.NewMemoryTokenStore())
	// access token generate method: jwt
	mgr.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte(global.Config.Oauth2.JWTSignedKey), jwt.SigningMethodHS512))
	clientStore := store.NewClientStore()
	for _, v := range global.Config.Oauth2.Client {
		_ = clientStore.Set(v.ID, &models.Client{
			ID:     v.ID,
			Secret: v.Secret,
			Domain: v.Addr + ":" + v.Port,
		})
	}
	mgr.MapClientStorage(clientStore)

	//err := api.InitRouter()
	//if err != nil {
	//	fmt.Println("初始化路由失败")
	//}

	srv := server.NewServer(server.NewConfig(), mgr)
	srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	srv.SetInternalErrorHandler(internalErrorHandler)
	srv.SetResponseErrorHandler(responseErrorHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {

		var form url.Values
		if v, _ := session.Get(r, "RequestForm"); v != nil {
			form = v.(url.Values)
		}
		r.Form = form

		if err := session.Delete(w, r, "RequestForm"); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {

		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {

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
			"client_id":  token.GetClientID(),
			"user_id":    cli.GetUserID(),
			"domain":     cli.GetDomain(),
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(data)
	})
	log.Printf("Server is running at 9096 port.\n")
	log.Printf("Point your OAuth client Auth endpoint to http://localhost:9096%s", "/oauth/authorize")
	log.Printf("Point your OAuth client Token endpoint to http://localhost:9096%s", "/oauth/token")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":9096"), nil))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			internalErrorHandler(err)
			return
		}
	}

	// 检查redirect_uri参数
	redirectURI := r.Form.Get("redirect_uri")
	if redirectURI == "" {
		internalErrorHandler(errors.New("参数不能为空(redirect_uri)"))
		return
	}
	if _, err := url.Parse(redirectURI); err != nil {
		internalErrorHandler(errors.New("参数无效(redirect_uri)"))
		return
	}

	// 删除公共回话
	if err := session.Delete(w, r, "LoggedInUserID"); err != nil {
		internalErrorHandler(errors.New("delete session failed"))
		return
	}

	w.Header().Set("Location", redirectURI)
	w.WriteHeader(http.StatusFound)
}

func responseErrorHandler(re *errors.Response) {
	log.Println("Response Error:", re.Error.Error())
}

func internalErrorHandler(err error) (re *errors.Response) {
	log.Println("Internal Error:", err.Error())
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

func passwordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	username2, _ := strconv.ParseInt(username, 10, 64)
	userID, err = dao.Authentication(ctx, username2, password)

	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		session.Set(w, r, "LoggedInUserID", r.Form.Get("username"))

		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {

	if v, _ := session.Get(r, "LoggedInUserID"); v == nil {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	outputHTML(w, r, "static/auth.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
