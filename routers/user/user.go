package user

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/modules/utils"
    "github.com/ouqiang/gocron/models"
    "github.com/go-macaron/session"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/go-macaron/captcha"
)

// @author qiang.ou<qingqianludao@gmail.com>
// @date 2017/4/23-14:16

func Login(ctx *macaron.Context)  {
    ctx.Data["Title"] = "用户登录"
    ctx.HTML(200, "user/login")
}

func EditPassword(ctx *macaron.Context)  {
    ctx.Data["Title"] = "修改密码"
    ctx.HTML(200, "user/editPassword")
}

func UpdatePassword(ctx *macaron.Context, sess session.Store) string  {
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

func ValidateLogin(ctx *macaron.Context, sess session.Store, cpt *captcha.Captcha) string {
    username := ctx.QueryTrim("username")
    password := ctx.QueryTrim("password")
    json := utils.JsonResponse{}
    if username == "" || password == "" {
        return json.CommonFailure("用户名、密码不能为空")
    }
    userModel := new (models.User)
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

func Logout(ctx *macaron.Context, sess session.Store) {
    if IsLogin(sess) {
        err := sess.Destory(ctx)
        if err != nil {
            logger.Error("用户退出登录失败", err)
        }
    }

    Login(ctx)
}

func Username(sess session.Store) string  {
    username,ok := sess.Get("username").(string)
    if ok {
        return username
    }

    return ""
}

func Uid(sess session.Store) int  {
    uid,ok := sess.Get("uid").(int)
    if ok {
        return uid
    }

    return 0
}

func IsLogin(sess session.Store) bool  {
    uid, ok := sess.Get("uid").(int)
    if ok && uid > 0 {
        return true
    }

    return false
}

func IsAdmin(sess session.Store) bool  {
    isAdmin, ok := sess.Get("isAdmin").(int8)
    if ok && isAdmin > 0 {
        return true
    }

    return false
}