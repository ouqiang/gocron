package user

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/Unknwon/paginater"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"github.com/ouqiang/gocron/internal/routers/base"
	"gopkg.in/macaron.v1"
)

// @author qiang.ou<qingqianludao@gmail.com>
// @date 2017/4/23-14:16

// UserForm 用户表单
type UserForm struct {
	Id              int
	Name            string `binding:"Required;MaxSize(32)"` // 用户名
	Password        string // 密码
	ConfirmPassword string // 确认密码
	Email           string `binding:"Required;MaxSize(50)"` // 邮箱
	IsAdmin         int8   // 是否是管理员 1:管理员 0:普通用户
	Status          models.Status
}

// Index 用户列表页
func Index(ctx *macaron.Context) {
	queryParams := parseQueryParams(ctx)
	userModel := new(models.User)
	users, err := userModel.List(queryParams)
	if err != nil {
		logger.Error(err)
	}
	total, err := userModel.Total()
	if err != nil {
		logger.Error(err)
	}
	PageParams := fmt.Sprintf("page_size=%d", queryParams["PageSize"])
	queryParams["PageParams"] = template.URL(PageParams)
	p := paginater.New(int(total), queryParams["PageSize"].(int), queryParams["Page"].(int), 5)
	ctx.Data["Pagination"] = p
	ctx.Data["Params"] = queryParams
	ctx.Data["Title"] = "用户列表"
	ctx.Data["Users"] = users
	ctx.HTML(200, "user/index")
}

// 解析查询参数
func parseQueryParams(ctx *macaron.Context) models.CommonMap {
	var params models.CommonMap = models.CommonMap{}
	base.ParsePageAndPageSize(ctx, params)

	return params
}

// Create 新增用户页
func Create(ctx *macaron.Context) {
	userModel := new(models.User)
	userModel.Status = models.Enabled
	userModel.IsAdmin = 0
	ctx.Data["User"] = userModel
	ctx.Data["Title"] = "添加用户"
	ctx.HTML(200, "user/user_form")
}

// 编辑页面
func Edit(ctx *macaron.Context) {
	ctx.Data["Title"] = "编辑用户"
	userModel := new(models.User)
	id := ctx.ParamsInt(":id")
	err := userModel.Find(id)
	if err != nil {
		logger.Error(err)
	}
	ctx.Data["User"] = userModel
	ctx.HTML(200, "user/user_form")
}

// 保存任务
func Store(ctx *macaron.Context, form UserForm) string {
	form.Name = strings.TrimSpace(form.Name)
	form.Email = strings.TrimSpace(form.Email)
	form.Password = strings.TrimSpace(form.Password)
	form.ConfirmPassword = strings.TrimSpace(form.ConfirmPassword)
	json := utils.JsonResponse{}
	userModel := models.User{}
	nameExists, err := userModel.UsernameExists(form.Name, form.Id)
	if err != nil {
		return json.CommonFailure(utils.FailureContent, err)
	}
	if nameExists > 0 {
		return json.CommonFailure("用户名已存在")
	}

	emailExists, err := userModel.EmailExists(form.Email, form.Id)
	if err != nil {
		return json.CommonFailure(utils.FailureContent, err)
	}
	if emailExists > 0 {
		return json.CommonFailure("邮箱已存在")
	}

	if form.Id == 0 {
		if form.Password == "" {
			return json.CommonFailure("请输入密码")
		}
		if form.ConfirmPassword == "" {
			return json.CommonFailure("请再次输入密码")
		}
		if form.Password != form.ConfirmPassword {
			return json.CommonFailure("两次密码输入不一致")
		}
	}
	userModel.Name = form.Name
	userModel.Email = form.Email
	userModel.Password = form.Password
	userModel.IsAdmin = form.IsAdmin
	userModel.Status = form.Status

	if form.Id == 0 {
		_, err = userModel.Create()
		if err != nil {
			return json.CommonFailure("添加失败", err)
		}
	} else {
		_, err = userModel.Update(form.Id, models.CommonMap{
			"name":     form.Name,
			"email":    form.Email,
			"status":   form.Status,
			"is_admin": form.IsAdmin,
		})
		if err != nil {
			return json.CommonFailure("修改失败", err)
		}
	}

	return json.Success("保存成功", nil)
}

// 删除用户
func Remove(ctx *macaron.Context) string {
	id := ctx.ParamsInt(":id")
	json := utils.JsonResponse{}

	userModel := new(models.User)
	_, err := userModel.Delete(id)
	if err != nil {
		return json.CommonFailure(utils.FailureContent, err)
	}

	return json.Success(utils.SuccessContent, nil)
}

// 激活用户
func Enable(ctx *macaron.Context) string {
	return changeStatus(ctx, models.Enabled)
}

// 禁用用户
func Disable(ctx *macaron.Context) string {
	return changeStatus(ctx, models.Disabled)
}

// 改变任务状态
func changeStatus(ctx *macaron.Context, status models.Status) string {
	id := ctx.ParamsInt(":id")
	json := utils.JsonResponse{}
	userModel := new(models.User)
	_, err := userModel.Update(id, models.CommonMap{
		"Status": status,
	})
	if err != nil {
		return json.CommonFailure(utils.FailureContent, err)
	}

	return json.Success(utils.SuccessContent, nil)
}

// Login 用户登录
func Login(ctx *macaron.Context) {
	ctx.Data["Title"] = "用户登录"
	ctx.HTML(200, "user/login")
}

// EditPassword 修改密码页面
func EditPassword(ctx *macaron.Context) {
	id := ctx.ParamsInt(":id")
	ctx.Data["Title"] = "修改密码"
	ctx.Data["Id"] = id
	ctx.HTML(200, "user/editPassword")
}

// UpdatePassword 更新我的密码
func UpdatePassword(ctx *macaron.Context) string {
	id := ctx.ParamsInt(":id")
	newPassword := ctx.QueryTrim("new_password")
	confirmNewPassword := ctx.QueryTrim("confirm_new_password")
	json := utils.JsonResponse{}
	if newPassword == "" || confirmNewPassword == "" {
		return json.CommonFailure("请输入密码")
	}
	userModel := new(models.User)
	_, err := userModel.UpdatePassword(id, newPassword)
	if err != nil {
		return json.CommonFailure("修改失败")
	}

	return json.Success("修改成功", nil)
}

// EditMyPassword 修改我的密码页面
func EditMyPassword(ctx *macaron.Context) {
	ctx.Data["Title"] = "修改密码"
	ctx.HTML(200, "user/editMyPassword")
}

// UpdateMyPassword 更新我的密码
func UpdateMyPassword(ctx *macaron.Context, sess session.Store) string {
	oldPassword := ctx.QueryTrim("old_password")
	newPassword := ctx.QueryTrim("new_password")
	confirmNewPassword := ctx.QueryTrim("confirm_new_password")
	json := utils.JsonResponse{}
	if oldPassword == "" || newPassword == "" || confirmNewPassword == "" {
		return json.CommonFailure("原密码和新密码均不能为空")
	}
	if newPassword != confirmNewPassword {
		return json.CommonFailure("两次输入密码不一致")
	}
	if oldPassword == newPassword {
		return json.CommonFailure("原密码与新密码不能相同")
	}
	userModel := new(models.User)
	if !userModel.Match(Username(sess), oldPassword) {
		return json.CommonFailure("原密码输入错误")
	}
	_, err := userModel.UpdatePassword(Uid(sess), newPassword)
	if err != nil {
		return json.CommonFailure("修改失败")
	}

	return json.Success("修改成功", nil)
}

// ValidateLogin 验证用户登录
func ValidateLogin(ctx *macaron.Context, sess session.Store, cpt *captcha.Captcha) string {
	username := ctx.QueryTrim("username")
	password := ctx.QueryTrim("password")
	json := utils.JsonResponse{}
	if username == "" || password == "" {
		return json.CommonFailure("用户名、密码不能为空")
	}
	userModel := new(models.User)
	if !userModel.Match(username, password) {
		return json.CommonFailure("用户名或密码错误")
	}
	if !cpt.VerifyReq(ctx.Req) {
		return json.Failure(utils.CaptchaError, "验证码错误")
	}

	loginLogModel := new(models.LoginLog)
	loginLogModel.Username = userModel.Name
	loginLogModel.Ip = ctx.RemoteAddr()
	_, err := loginLogModel.Create()
	if err != nil {
		logger.Error("记录用户登录日志失败", err)
	}

	sess.Set("username", userModel.Name)
	sess.Set("uid", userModel.Id)
	sess.Set("isAdmin", userModel.IsAdmin)

	return json.Success("登录成功", nil)
}

// Logout 用户退出
func Logout(ctx *macaron.Context, sess session.Store) {
	if IsLogin(sess) {
		err := sess.Destory(ctx)
		if err != nil {
			logger.Error("用户退出登录失败", err)
		}
	}

	Login(ctx)
}

// Username 获取session中的用户名
func Username(sess session.Store) string {
	username, ok := sess.Get("username").(string)
	if ok {
		return username
	}

	return ""
}

// Uid 获取session中的Uid
func Uid(sess session.Store) int {
	uid, ok := sess.Get("uid").(int)
	if ok {
		return uid
	}

	return 0
}

// IsLogin 判断用户是否已登录
func IsLogin(sess session.Store) bool {
	uid, ok := sess.Get("uid").(int)
	if ok && uid > 0 {
		return true
	}

	return false
}

// IsAdmin 判断当前用户是否是管理员
func IsAdmin(sess session.Store) bool {
	isAdmin, ok := sess.Get("isAdmin").(int8)
	if ok && isAdmin > 0 {
		return true
	}

	return false
}
