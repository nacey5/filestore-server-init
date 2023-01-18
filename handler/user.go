package handler

import (
	_const "filestore-server/const"
	dblayer "filestore-server/db"
	"filestore-server/util"
	"net/http"
	"os"
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
