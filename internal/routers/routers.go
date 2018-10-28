package routers

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/gzip"
	"github.com/go-macaron/toolbox"
	"github.com/ouqiang/gocron/internal/modules/app"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"github.com/ouqiang/gocron/internal/routers/host"
	"github.com/ouqiang/gocron/internal/routers/install"
	"github.com/ouqiang/gocron/internal/routers/loginlog"
	"github.com/ouqiang/gocron/internal/routers/manage"
	"github.com/ouqiang/gocron/internal/routers/task"
	"github.com/ouqiang/gocron/internal/routers/tasklog"
	"github.com/ouqiang/gocron/internal/routers/user"
	"github.com/rakyll/statik/fs"
	"gopkg.in/macaron.v1"

	_ "github.com/ouqiang/gocron/internal/statik"
)

const (
	urlPrefix = "/api"
	staticDir = "public"
)

var statikFS http.FileSystem

func init() {
	var err error
	statikFS, err = fs.New()
	if err != nil {
		log.Fatal(err)
	}
}

// Register 路由注册
func Register(m *macaron.Macaron) {
	m.SetURLPrefix(urlPrefix)
	// 所有GET方法，自动注册HEAD方法
	m.SetAutoHead(true)
	m.Get("/", func(ctx *macaron.Context) {
		file, err := statikFS.Open("/index.html")
		if err != nil {
			logger.Error("读取首页文件失败: %s", err)
			ctx.WriteHeader(http.StatusInternalServerError)
			return
		}

		io.Copy(ctx.Resp, file)

	})
	// 系统安装
	m.Group("/install", func() {
		m.Post("/store", binding.Bind(install.InstallForm{}), install.Store)
		m.Get("/status", func(ctx *macaron.Context) string {
			jsonResp := utils.JsonResponse{}
			return jsonResp.Success("", app.Installed)
		})
	})

	// 用户
	m.Group("/user", func() {
		m.Get("", user.Index)
		m.Get("/:id", user.Detail)
		m.Post("/store", binding.Bind(user.UserForm{}), user.Store)
		m.Post("/remove/:id", user.Remove)
		m.Post("/login", user.ValidateLogin)
		m.Post("/enable/:id", user.Enable)
		m.Post("/disable/:id", user.Disable)
		m.Post("/editMyPassword", user.UpdateMyPassword)
		m.Post("/editPassword/:id", user.UpdatePassword)
	})

	// 定时任务
	m.Group("/task", func() {
		m.Post("/store", binding.Bind(task.TaskForm{}), task.Store)
		m.Get("/:id", task.Detail)
		m.Get("", task.Index)
		m.Get("/log", tasklog.Index)
		m.Post("/log/clear", tasklog.Clear)
		m.Post("/log/stop", tasklog.Stop)
		m.Post("/remove/:id", task.Remove)
		m.Post("/enable/:id", task.Enable)
		m.Post("/disable/:id", task.Disable)
		m.Get("/run/:id", task.Run)
	})

	// 主机
	m.Group("/host", func() {
		m.Get("/:id", host.Detail)
		m.Post("/store", binding.Bind(host.HostForm{}), host.Store)
		m.Get("", host.Index)
		m.Get("/all", host.All)
		m.Get("/ping/:id", host.Ping)
		m.Post("/remove/:id", host.Remove)
	})

	// 管理
	m.Group("/system", func() {
		m.Group("/slack", func() {
			m.Get("", manage.Slack)
			m.Post("/update", manage.UpdateSlack)
			m.Post("/channel", manage.CreateSlackChannel)
			m.Post("/channel/remove/:id", manage.RemoveSlackChannel)
		})
		m.Group("/mail", func() {
			m.Get("", manage.Mail)
			m.Post("/update", binding.Bind(manage.MailServerForm{}), manage.UpdateMail)
			m.Post("/user", manage.CreateMailUser)
			m.Post("/user/remove/:id", manage.RemoveMailUser)
		})
		m.Group("/webhook", func() {
			m.Get("", manage.WebHook)
			m.Post("/update", manage.UpdateWebHook)
		})
		m.Get("/login-log", loginlog.Index)
	})

	// API
	m.Group("/v1", func() {
		m.Post("/tasklog/remove/:id", tasklog.Remove)
		m.Post("/task/enable/:id", task.Enable)
		m.Post("/task/disable/:id", task.Disable)
	}, apiAuth)

	// 404错误
	m.NotFound(func(ctx *macaron.Context) string {
		jsonResp := utils.JsonResponse{}

		return jsonResp.Failure(utils.NotFound, "您访问的页面不存在")
	})
	// 50x错误
	m.InternalServerError(func(ctx *macaron.Context) string {
		jsonResp := utils.JsonResponse{}

		return jsonResp.Failure(utils.ServerError, "服务器内部错误, 请稍后再试")
	})
}

// 中间件注册
func RegisterMiddleware(m *macaron.Macaron) {
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	if macaron.Env != macaron.DEV {
		m.Use(gzip.Gziper())
	}
	m.Use(
		macaron.Static(
			"",
			macaron.StaticOptions{
				Prefix:     staticDir,
				FileSystem: statikFS,
			},
		),
	)
	if macaron.Env == macaron.DEV {
		m.Use(toolbox.Toolboxer(m))
	}
	m.Use(macaron.Renderer())
	m.Use(checkAppInstall)
	m.Use(ipAuth)
	m.Use(userAuth)
	m.Use(urlAuth)
}

// region 自定义中间件

/** 检测应用是否已安装 **/
func checkAppInstall(ctx *macaron.Context) {
	if app.Installed {
		return
	}
	if strings.HasPrefix(ctx.Req.URL.Path, "/install") || ctx.Req.URL.Path == "/" {
		return
	}
	jsonResp := utils.JsonResponse{}

	data := jsonResp.Failure(utils.AppNotInstall, "应用未安装")
	ctx.Write([]byte(data))
}

// IP验证, 通过反向代理访问gocron，需设置Header X-Real-IP才能获取到客户端真实IP
func ipAuth(ctx *macaron.Context) {
	if !app.Installed {
		return
	}
	allowIpsStr := app.Setting.AllowIps
	if allowIpsStr == "" {
		return
	}
	clientIp := ctx.RemoteAddr()
	allowIps := strings.Split(allowIpsStr, ",")
	if utils.InStringSlice(allowIps, clientIp) {
		return
	}
	logger.Warnf("非法IP访问-%s", clientIp)
	jsonResp := utils.JsonResponse{}

	data := jsonResp.Failure(utils.UnauthorizedError, "您无权限访问")

	ctx.Write([]byte(data))
}

// 用户认证
func userAuth(ctx *macaron.Context) {
	if !app.Installed {
		return
	}
	user.RestoreToken(ctx)
	if user.IsLogin(ctx) {
		return
	}
	uri := strings.TrimRight(ctx.Req.URL.Path, "/")
	if strings.HasPrefix(uri, "/v1") {
		return
	}
	excludePaths := []string{"", "/user/login", "/install/status"}
	for _, path := range excludePaths {
		if uri == path {
			return
		}
	}
	jsonResp := utils.JsonResponse{}
	data := jsonResp.Failure(utils.AuthError, "认证失败")
	ctx.Write([]byte(data))

}

// URL权限验证
func urlAuth(ctx *macaron.Context) {
	if !app.Installed {
		return
	}
	if user.IsAdmin(ctx) {
		return
	}
	uri := strings.TrimRight(ctx.Req.URL.Path, "/")
	if strings.HasPrefix(uri, "/v1") {
		return
	}
	// 普通用户允许访问的URL地址
	allowPaths := []string{
		"",
		"/install/status",
		"/task",
		"/task/log",
		"/host",
		"/host/all",
		"/user/login",
		"/user/editMyPassword",
	}
	for _, path := range allowPaths {
		if path == uri {
			return
		}
	}

	jsonResp := utils.JsonResponse{}

	data := jsonResp.Failure(utils.UnauthorizedError, "您无权限访问")
	ctx.Write([]byte(data))
}

/** API接口签名验证 **/
func apiAuth(ctx *macaron.Context) {
	if !app.Installed {
		return
	}
	if !app.Setting.ApiSignEnable {
		return
	}
	apiKey := strings.TrimSpace(app.Setting.ApiKey)
	apiSecret := strings.TrimSpace(app.Setting.ApiSecret)
	json := utils.JsonResponse{}
	if apiKey == "" || apiSecret == "" {
		msg := json.CommonFailure("使用API前, 请先配置密钥")
		ctx.Write([]byte(msg))
		return
	}
	currentTimestamp := time.Now().Unix()
	time := ctx.QueryInt64("time")
	if time <= 0 {
		msg := json.CommonFailure("参数time不能为空")
		ctx.Write([]byte(msg))
		return
	}
	if time < (currentTimestamp - 1800) {
		msg := json.CommonFailure("time无效")
		ctx.Write([]byte(msg))
		return
	}
	sign := ctx.QueryTrim("sign")
	if sign == "" {
		msg := json.CommonFailure("参数sign不能为空")
		ctx.Write([]byte(msg))
		return
	}
	raw := apiKey + strconv.FormatInt(time, 10) + strings.TrimSpace(ctx.Req.URL.Path) + apiSecret
	realSign := utils.Md5(raw)
	if sign != realSign {
		msg := json.CommonFailure("签名验证失败")
		ctx.Write([]byte(msg))
		return
	}
}

// endregion
