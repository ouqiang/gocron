package routers

import (
    "github.com/go-macaron/binding"
    "github.com/ouqiang/gocron/routers/install"
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/routers/task"
    "github.com/ouqiang/gocron/routers/host"
    "github.com/ouqiang/gocron/routers/tasklog"
    "github.com/ouqiang/gocron/modules/utils"
    "github.com/go-macaron/session"
    "github.com/go-macaron/toolbox"
    "strings"
    "github.com/ouqiang/gocron/modules/app"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/routers/user"
    "github.com/go-macaron/gzip"
    "github.com/ouqiang/gocron/routers/manage"
    "github.com/go-macaron/csrf"
    "github.com/ouqiang/gocron/routers/loginlog"
)

// 静态文件目录
const StaticDir = "public"

// 路由注册
func Register(m *macaron.Macaron) {
    // 所有GET方法，自动注册HEAD方法
    m.SetAutoHead(true)
    // 首页
    m.Get("/", Home)
    // 系统安装
    m.Group("/install", func() {
        m.Get("", install.Create)
        m.Post("/store", binding.Bind(install.InstallForm{}), install.Store)
    })

    // 用户
    m.Group("/user", func() {
        m.Get("/login", user.Login)
        m.Post("/login", user.ValidateLogin)
        m.Get("/logout", user.Logout)
        m.Get("/editPassword", user.EditPassword)
        m.Post("/editPassword", user.UpdatePassword)
    })

    // 任务
    m.Group("/task", func() {
        m.Get("/create", task.Create)
        m.Post("/store", binding.Bind(task.TaskForm{}), task.Store)
        m.Get("/edit/:id", task.Edit)
        m.Get("", task.Index)
        m.Get("/log", tasklog.Index)
        m.Post("/log/clear", tasklog.Clear)
        m.Post("/remove/:id", task.Remove)
        m.Post("/enable/:id", task.Enable)
        m.Post("/disable/:id", task.Disable)
        m.Get("/run/:id", task.Run)
    })

    // 主机
    m.Group("/host", func() {
        m.Get("/create", host.Create)
        m.Get("/edit/:id", host.Edit)
        m.Get("/ping/:id", host.Ping)
        m.Post("/store", binding.Bind(host.HostForm{}), host.Store)
        m.Get("", host.Index)
        m.Post("/remove/:id", host.Remove)
    })

    // 管理
    m.Group("/manage", func() {
        m.Group("/slack", func() {
            m.Get("/", manage.Slack)
            m.Get("/edit", manage.EditSlack)
            m.Post("/url", manage.UpdateSlackUrl)
            m.Post("/channel", manage.CreateSlackChannel)
            m.Post("/channel/remove/:id", manage.RemoveSlackChannel)
        })
        m.Group("/mail", func() {
            m.Get("/", manage.Mail)
            m.Get("/edit", manage.EditMail)
            m.Post("/server", binding.Bind(manage.MailServerForm{}), manage.UpdateMailServer)
            m.Post("/user", manage.CreateMailUser)
            m.Post("/user/remove/:id", manage.RemoveMailUser)
        })
        m.Get("/login-log", loginlog.Index)
    })

    // API
    m.Group("/api/v1", func() {
       m.Get("/tasklog/update-status", tasklog.UpdateStatus)
    });

    // 404错误
    m.NotFound(func(ctx *macaron.Context) {
        if isGetRequest(ctx) && !isAjaxRequest(ctx) {
            ctx.Data["Title"] = "404 - NOT FOUND"
            ctx.HTML(404, "error/404")
        } else {
            json := utils.JsonResponse{}
            ctx.Resp.Write([]byte(json.Failure(utils.NotFound, "您访问的地址不存在")))
        }
    })
    // 50x错误
    m.InternalServerError(func(ctx *macaron.Context) {
        logger.Debug("500错误")
        if isGetRequest(ctx) && !isAjaxRequest(ctx) {
            ctx.Data["Title"] = "500 - INTERNAL SERVER ERROR"
            ctx.HTML(500, "error/500")
        } else {
            json := utils.JsonResponse{}
            ctx.Resp.Write([]byte(json.Failure(utils.ServerError, "网站暂时无法访问,请稍后再试")))
        }
    })
}

// 中间件注册
func RegisterMiddleware(m *macaron.Macaron) {
    m.Use(macaron.Logger())
    m.Use(macaron.Recovery())
    if macaron.Env != macaron.DEV {
        m.Use(gzip.Gziper())
    }
    m.Use(macaron.Static(StaticDir))
    m.Use(macaron.Renderer(macaron.RenderOptions{
        Directory:  "templates",
        Extensions: []string{".html"},
        // 模板语法分隔符，默认为 ["{{", "}}"]
        Delims: macaron.Delims{"{{{", "}}}"},
        // 追加的 Content-Type 头信息，默认为 "UTF-8"
        Charset: "UTF-8",
        // 渲染具有缩进格式的 JSON，默认为不缩进
        IndentJSON: true,
        // 渲染具有缩进格式的 XML，默认为不缩进
        IndentXML: true,
    }))
    m.Use(session.Sessioner(session.Options{
        Provider:       "file",
        ProviderConfig: app.DataDir + "/sessions",
    }))
    m.Use(csrf.Csrfer())
    m.Use(toolbox.Toolboxer(m))
    checkAppInstall(m)
    userAuth(m)
    setShareData(m)
}

// 系统未安装，重定向到安装页面
func checkAppInstall(m *macaron.Macaron)  {
    m.Use(func(ctx *macaron.Context) {
        installUrl := "/install"
        if strings.HasPrefix(ctx.Req.URL.Path, installUrl) {
            return
        }
        if !app.Installed {
            ctx.Redirect(installUrl)
        }
    })
}

// 用户认证
func userAuth(m *macaron.Macaron)  {
    m.Use(func(ctx *macaron.Context, sess session.Store) {
        if user.IsLogin(sess) {
            return
        }
        uri := ctx.Req.URL.Path
        found := false
        excludePaths := []string{"/install", "/user/login", "/"}
        for _, path := range excludePaths {
            if strings.HasPrefix(uri, path) {
                found = true
                break
            }
        }
        if !found {
            ctx.Redirect("/user/login")
        }
    })
}

// 设置共享数据
func setShareData(m *macaron.Macaron)  {
    m.Use(func(ctx *macaron.Context, sess session.Store) {
        ctx.Data["URI"] = ctx.Req.URL.Path
        urlPath := strings.TrimPrefix(ctx.Req.URL.Path, "/")
        paths := strings.Split(urlPath, "/")
        ctx.Data["Controller"] = ""
        ctx.Data["Action"] = ""
        if len(paths) > 0 {
            ctx.Data["Controller"] = paths[0]
        }
        if len(paths) > 1 {
            ctx.Data["Action"] = paths[1]
        }
        ctx.Data["LoginUsername"] = user.Username(sess)
        ctx.Data["LoginUid"] = user.Uid(sess)
    })
}

// 管理员认证
func adminAuth(ctx *macaron.Context, sess session.Store)  {
    if !user.IsAdmin(sess) {
        ctx.Data["Title"] = "无权限访问此页面"
        ctx.HTML(403, "error/no_permission")
    }
}

func isAjaxRequest(ctx *macaron.Context) bool {
    req := ctx.Req.Header.Get("X-Requested-With")
    if req == "XMLHttpRequest" {
        return true
    }

    return false
}

func isGetRequest(ctx *macaron.Context) bool {
    return ctx.Req.Method == "GET"
}