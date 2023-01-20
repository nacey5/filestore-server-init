package handler

import (
	_const "filestore-server/const"
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	pwd_salt = "*#890"
)

// SignupHandler 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == _const.GET {
		data, err := os.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	} else if r.Method == _const.POST {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if len(username) < 3 || len(password) < 5 {
			w.Write([]byte("Invalid parameter"))
			return
		}
		//密码加密处理
		enc_password := util.Sha1([]byte(password + pwd_salt))
		suc := dblayer.UserSignup(username, enc_password)
		if suc {
			w.Write([]byte("SUCCESS"))
		} else {
			w.Write([]byte("FAILED"))
		}
	}
}

// SignInHandler 登陆接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPassword := util.Sha1([]byte(password + pwd_salt))
	//1.校验用户名及其密码
	pwdChecked := dblayer.UserSignIn(username, encPassword)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}
	//2.生成访问凭证(token)
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}
	//3.登陆成功后重定向到首页
	w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
}

func GenToken(username string) string {
	//40位:md5(username+timeStamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// UserInfoHandler 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")
	//2.验证token是否有效:::::在拦截器中定义~~
	//isValidToken := util.IsTokenValid(token)
	//if !isValidToken {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	//3.查询用户信息 todo 查询数据库
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//4.组装并且响应数据
	resp := util.NewRespMsg(0, "ok", user)
	w.Write(resp.JSONBytes())
}
