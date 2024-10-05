/*
* @desc: 后台路由
* @company: 云南奇讯科技有限公司
* @Author: yixiaohu
* @Date:   2022/2/18 17:34
 */

package router

import (
	// 导入GoFrame框架的HTTP包，用于处理路由
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	// 导入后台控制器
	"github.com/tiger1103/gfast/v3/internal/app/system/controller"
	// 导入系统服务，如Token验证和中间件
	"github.com/tiger1103/gfast/v3/internal/app/system/service"
	// 导入自动路由绑定库
	"github.com/tiger1103/gfast/v3/library/libRouter"
)

// 定义全局的路由结构体实例
var R = new(Router)

// Router 路由结构体
type Router struct{}

// BindController 绑定后台控制器路由
// @param ctx   上下文对象
// @param group 路由组对象，用于绑定路由和中间件
func (router *Router) BindController(ctx context.Context, group *ghttp.RouterGroup) {
	// 定义路由组，路径前缀为 /system
	group.Group("/system", func(group *ghttp.RouterGroup) {
		// 绑定登录控制器，处理登录相关路由
		group.Bind(
			// 登录控制器
			controller.Login,
		)
		// 登录验证拦截器，所有后续路由都需要通过Token验证
		service.GfToken().Middleware(group)
		// 上下文拦截器，用于注入请求上下文和权限验证
		group.Middleware(service.Middleware().Ctx, service.Middleware().Auth)
		// 后台操作日志钩子，在请求完成后记录操作日志
		group.Hook("/*", ghttp.HookAfterOutput, service.OperateLog().OperationLog)
		// 绑定多个后台控制器，处理不同的后台功能
		group.Bind(
			controller.User,       // 用户管理
			controller.Menu,       // 菜单管理
			controller.Role,       // 角色管理
			controller.Dept,       // 部门管理
			controller.Post,       // 岗位管理
			controller.DictType,   // 字典类型管理
			controller.DictData,   // 字典数据管理
			controller.Config,     // 配置管理
			controller.Monitor,    // 系统监控
			controller.LoginLog,   // 登录日志管理
			controller.OperLog,    // 操作日志管理
			controller.Personal,   // 个人中心
			controller.UserOnline, // 在线用户管理
			controller.Cache,      // 缓存处理
		)
		// 自动绑定定义的控制器
		if err := libRouter.RouterAutoBind(ctx, router, group); err != nil {
			// 如果自动绑定发生错误，抛出异常
			panic(err)
		}
	})
}
