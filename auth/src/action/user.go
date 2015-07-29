package action

import (
	"errors"
	_"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/dchest/authcookie"
	"time"
	"net/http"
)
const (
	//cookie加密、解密使用
	KEY        string = "QAZWERT4556"
	COOKIENAME string = "MNBVCXZ"
)

//生成cookie，放到reponse对象
func GenerateCookie(w http.ResponseWriter, r *http.Request, userNmae string, number int) {
	timeLength := 24 * time.Hour
	cookieValue := authcookie.NewSinceNow(userNmae, timeLength, []byte(KEY))
	expire := time.Now().Add(timeLength)
	cookie := http.Cookie{Name: COOKIENAME, Value: cookieValue, Path: "/", Expires: expire, MaxAge: 86400}
	http.SetCookie(w, &cookie)
}

func GetCookieName(req *http.Request) string {
	cookie, err := req.Cookie(COOKIENAME)
	if err == nil {
		return authcookie.Login(cookie.Value, []byte(KEY))
	}
	return ""
}

func Register(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	acname := req.FormValue("acname")
	password := req.FormValue("password")
	account := Account{Ac_name: acname, Ac_password: password}
	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	ok := RegisterInsert(&account)
	if !ok {
		strBody = []byte("{\"Code\":1,\"Message\":\"insert db faild!!\"}")
	}
	w.Write(strBody)
}

func RegisterMulti(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()

	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	for k, v := range req.Form {
		ok := MultiRegister(k,v[0])
		if !ok {
			strBody = []byte("{\"Code\":1,\"Message\":\"insert db faild!!\"}")
		}
		break
	}
	w.Write(strBody)
}

/*
//检查是否登录，未登录，则返回登录页
func checkLogin(oauth *OAuth, w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("checkLogin\r\n")
	acname := Logged(w, r)
	fmt.Println("checkLogin acname", acname)
	if acname == "" {
		common.ForwardPage(w, "./static/public/oauth2/login.html", map[string]string{"RequestURI": "/oauth2/login?" + r.URL.RawQuery})
		return false
	}
	return true
}

func Logged(w http.ResponseWriter, req *http.Request) string {
	cookie, err := req.Cookie(COOKIENAME)
	if err == nil {
		return authcookie.Login(cookie.Value, []byte(KEY))
	}
	return ""
}
*/

func Login(w http.ResponseWriter, req *http.Request) (string, error) {
	acname := req.FormValue("acname")
	password := req.FormValue("password")
	if acname == "" || password == "" {
		return "", errors.New("未输入用户名和密码！")
	}
	user := User{Acname: acname, Password: password}
	ok := LoginQuery(&user)
	if ok {
		GenerateCookie(w, req, user.Acname, 1)
		return acname, nil
	} else {
		return "", errors.New("用户名或密码错误！")
	}
}

func Logout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	cookie := http.Cookie{Name: COOKIENAME, Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)
}

//通过openId获取用户资源权限列表
func SetUserInfo(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")
	acname := GetCookieName(req)
	if acname == "" {
		strBody = []byte("{\"Code\":1,\"Message\":\"user need login\"}")
		w.Write(strBody)
		return
	}

	//acname:="18585816540"
	acid := GetAcId(acname)
	if acid == -1 {
		strBody = []byte("{\"Code\":1,\"Message\":\"user not exist\"}")
		w.Write(strBody)
		return
	}

	//var  info map[string]string
	req.ParseForm()
	info := make(map[string]string)
	for k, v := range req.Form {
		info[k] = v[0]
	}

	var UserInfo ATUserInfo
	UserInfo.Ac_id = acid
	UserInfo.Info = info
	ok := UpdateUserInfo(&UserInfo)
	if !ok {
		strBody = []byte("{\"Code\":1,\"Message\":\"save data faild\"}")
	}
	w.Write(strBody)
}

func MultiLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")
	strName := req.FormValue("acname")
	strPassword := req.FormValue("password")
	if strName == "" || strPassword == "" {
		strBody = []byte("{\"Code\":1,\"Message\":\"password is empty!!\"}")
		w.Write(strBody)
		return
	}

	strACName, _ := LoginMulti(strName, strPassword)
	if len(strACName) > 0 {
		GenerateCookie(w, req, strACName, 1)
		w.Write(strBody)
		return
	} else {
		strBody = []byte("{\"Code\":1,\"Message\":\"password is empty!!\"}")
		return
	}
}