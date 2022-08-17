package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/mail"
	"goblog/pkg/view"
	"net/http"

	"github.com/thedevsaddam/govalidator"
)

// AuthController 处理静态页面
type AuthController struct {
}

// Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{
		"User": user.User{},
	}, "auth.register")
}

// DoRegister 处理注册逻辑
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {

	// 1. 初始化数据
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}

	// 2. 表单规则
	errs := requests.ValidateRegistrationForm(_user)

	if len(errs) > 0 {
		// 3. 表单不通过 —— 重新显示表单
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")
	} else {
		// 4. 验证成功，创建数据
		_user.Create()

		if _user.ID > 0 {
			// 登录用户并跳转到首页
			auth.Login(_user)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "注册失败，请联系管理员")
		}
	}
}

// Login 显示登录表单
func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{
		"Email":    "",
		"Password": "",
	}, "auth.login")
}

// DoLogin 处理登录表单提交
func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化表单数据
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// 2. 尝试登录
	if err := auth.Attempt(email, password); err == nil {
		// 登录成功
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		// 3. 失败，显示错误提示
		view.RenderSimple(w, view.D{
			"Error":    err.Error(),
			"Email":    email,
			"Password": password,
		}, "auth.login")
	}
}

// Logout 退出登录
func (*AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout()
	http.Redirect(w, r, "/", http.StatusFound)
}

// sendmail
func (*AuthController) SendMail(w http.ResponseWriter, r *http.Request) {
	// sendmail
	view.RenderSimple(w, view.D{
		"Email": "",
	}, "auth.send_mail")
}

// dosendmail
func (*AuthController) DoSendMail(w http.ResponseWriter, r *http.Request) {

	email := r.PostFormValue("email")

	rules := govalidator.MapData{
		"email": []string{"required", "min:4", "max:20", "email"},
	}

	messages := govalidator.MapData{
		"email": []string{"required:必须填写", "min:最少4个字符", "max:最多20个字符", "email:必须是邮箱"},
	}

	opts := govalidator.Options{
		Request:         r,        // request object
		Rules:           rules,    // rules map
		Messages:        messages, // custom message map (Optional)
		RequiredDefault: true,     // all the field to be pass the rules
	}

	v := govalidator.New(opts)
	err := v.Validate()

	if len(err) > 0 {

		// 存在错误的情况
		view.RenderSimple(w, view.D{
			"Email":  email,
			"Errors": err,
		}, "auth.send_mail")
	} else {

		// 发送邮件
		subject := "用户注册邮箱验证"
		body := `<a herf="#">点击认证</a>`
		result := mail.SendMail(email, subject, body)

		if result {
			fmt.Fprintf(w, "发送成功，请前往邮箱验证")
		} else {

			ers := make(map[string][]string)
			ers["email"] = []string{"发送失败"}
			// 发送失败
			view.RenderSimple(w, view.D{
				"Email":  email,
				"Errors": ers,
			}, "auth.send_mail")
		}
	}
}
