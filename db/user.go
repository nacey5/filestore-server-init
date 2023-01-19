package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
)

// UserSignup 通过用户名及密码完成user表的注册
func UserSignup(username, password string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user(`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert,err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Failed to insert.err:%s", err.Error())
		return false
	}
	//重复注册校验--只有是初次注册的时候
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

// 判断密码是否一致
func UserSignIn(username, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name= ? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found:" + username)
		return false
	}
	pRows := mydb.ParseRows(rows)

	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}

	return false
}

// UpdateToken 刷新用户登陆的token
func UpdateToken(username, token string) bool {
	stmt, err := mydb.DBConn().Prepare("replace into tbl_user_token(`user_name`,`user_token`) values (?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// 从数据库查询返回的数据格式
type User struct {
	Username   string
	Email      string
	Phone      string
	SignupAt   string
	LastActive string
	Status     int
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mydb.DBConn().Prepare("select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Printf(err.Error())
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Printf(err.Error())
		return user, err
	}

	return user, nil
}
