/*
* @desc: 路由绑定
* @company: 云南奇讯科技有限公司
* @Author: yixiaohu
* @Date:   2022/2/18 16:23
 */

package router

import (
	// 导入GoFrame框架的HTTP包，用于处理路由和请求
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	// 导入公共模块的路由定义
	commonRouter "github.com/tiger1103/gfast/v3/internal/app/common/router"
	// 导入公共服务，如中间件服务
	commonService "github.com/tiger1103/gfast/v3/internal/app/common/service"
	// 导入系统模块的路由定义
	systemRouter "github.com/tiger1103/gfast/v3/internal/app/system/router"
	// 导入自动路由绑定库
	"github.com/tiger1103/gfast/v3/library/libRouter"
)

// 定义全局的路由结构体实例
var R = new(Router)

// Router 路由结构体
type Router struct{}

// BindController 绑定控制器路由
// @param ctx   上下文对象
// @param group 路由组对象，用于绑定路由和中间件
func (router *Router) BindController(ctx context.Context, group *ghttp.RouterGroup) {
	// 定义路由组，路径前缀为 /api/v1
	group.Group("/api/v1", func(group *ghttp.RouterGroup) {
		// 跨域处理中间件，安全起见正式环境应注释该行
		group.Middleware(commonService.Middleware().MiddlewareCORS)
		// 全局响应处理器中间件，处理返回数据的格式
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		// 绑定后台模块的路由
		systemRouter.R.BindController(ctx, group)
		// 绑定公共模块的路由
		commonRouter.R.BindController(ctx, group)
		// 自动绑定自定义的路由模块
		if err := libRouter.RouterAutoBind(ctx, router, group); err != nil {
			// 如果自动绑定发生错误，抛出异常
			panic(err)
		}
	})
}
