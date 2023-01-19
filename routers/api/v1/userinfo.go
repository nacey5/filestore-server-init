package v1

import (
	_const "filestore-server/const"
	//dblayer "filestore-server/db"
	"filestore-server/global"
	"filestore-server/model"
	"filestore-server/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type UserMeta struct{}

func NewUserMeta() UserMeta {
	return UserMeta{}
}

func (u UserMeta) SignupHandler(c *gin.Context) {
	r := c.Request
	w := c.Writer

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
		enc_password := util.Sha1([]byte(password + _const.PWD_SALT))
		user := model.User{
			UserName: username,
			UserPwd:  enc_password,
		}
		err := user.Signup(global.DBEngine)
		if err == nil {
			w.Write([]byte("SUCCESS"))
		} else {
			w.Write([]byte("FAILED"))
		}
	}
}

func (u UserMeta) SignInHandler(c *gin.Context) {
	r := c.Request
	w := c.Writer

	if r.Method == _const.GET {
		data, err := os.ReadFile("./static/view/signin.html")
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
		encPassword := util.Sha1([]byte(password + _const.PWD_SALT))
		user := model.User{UserName: username, UserPwd: encPassword}
		userToken := model.UserToken{UserName: username}
		//pwdChecked := dblayer.UserSignIn(username, encPassword)
		pwdChecked := false
		err := user.Signin(global.DBEngine)
		if err == nil {
			pwdChecked = true
		}
		if !pwdChecked {
			w.Write([]byte("FAILED"))
			return
		}
		//2.生成访问凭证(token)
		token := util.GenToken(username)
		userToken.UserToken = token
		//upRes := dblayer.UpdateToken(username, token)
		upRes := false
		err = userToken.Update(global.DBEngine)
		if err == nil {
			upRes = true
		}
		if !upRes {
			w.Write([]byte("FAILED"))
			return
		}
		//3.登陆成功后重定向到首页
		//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
		resp := util.NewRespMsg(0, "ok", struct {
			Location string
			UserName string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			UserName: username,
			Token:    token,
		},
		)
		w.Write(resp.JSONBytes())
	}
}
