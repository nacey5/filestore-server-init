package v1

import (
	_const "filestore-server/const"
	"filestore-server/global"
	"filestore-server/model"
	"filestore-server/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

const (
	pwd_salt = "*#890"
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
		enc_password := util.Sha1([]byte(password + pwd_salt))
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
