package user

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/modules/utils"
    "github.com/ouqiang/gocron/models"
    "github.com/go-macaron/session"
    "github.com/ouqiang/gocron/modules/logger"
    "time"
)

// @author qiang.ou<qingqianludao@gmail.com>
// @date 2017/4/23-14:16

func Login(ctx *macaron.Context)  {
    ctx.Data["Title"] = "用户登录"
    ctx.HTML(200, "user/login")
}

func ValidateLogin(ctx *macaron.Context, sess session.Store) string {
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

    sess.Set("username", username)

    return json.Success("登录成功", nil)
}

func Logout(ctx *macaron.Context, sess session.Store) {
    if IsLogin(sess) {
        err := sess.Delete("username")
        if err != nil {
            logger.Error("用户退出登录失败", err)
        }
    }

    ctx.SetSecureCookie("MacaronSession", "", 0, "/", "", nil, nil, time.Now().AddDate(-1, 0, 0))
    Login(ctx)
}

func Username(sess session.Store) string  {
    username,ok := sess.Get("username").(string)
    if ok {
        return username
    }

    return ""
}

func IsLogin(sess session.Store) bool  {
    username,ok := sess.Get("username").(string)
    if ok && username != "" {
        return true
    }

    return false
}