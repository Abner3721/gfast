/*
* @desc: 登录控制器
* @company: 云南奇讯科技有限公司
* @Author: yixiaohu
* @Date:   2022/4/27 21:52
 */

package controller

import (
	// 导入必要的包，包括密码加密、错误处理、日志和用户相关服务
	"context"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmode"
	"github.com/tiger1103/gfast/v3/api/v1/system"
	commonService "github.com/tiger1103/gfast/v3/internal/app/common/service"
	"github.com/tiger1103/gfast/v3/internal/app/system/model"
	"github.com/tiger1103/gfast/v3/internal/app/system/service"
	"github.com/tiger1103/gfast/v3/library/libUtils"
)

// 定义登录控制器实例
var (
	Login = loginController{}
)

// 定义登录控制器结构体，继承基础控制器
type loginController struct {
	BaseController
}

// Login 登录处理方法
// @param ctx 上下文对象
// @param req 用户登录请求结构体
// @return res 用户登录响应结构体
// @return err 错误信息
func (c *loginController) Login(ctx context.Context, req *system.UserLoginReq) (res *system.UserLoginRes, err error) {
	var (
		user        *model.LoginUserRes // 登录用户信息
		token       string              // 登录生成的Token
		permissions []string            // 用户权限
		menuList    []*model.UserMenus  // 用户菜单列表
	)

	// 判断验证码是否正确，开发模式下跳过验证码校验
	debug := gmode.IsDevelop()
	if !debug {
		// 验证码校验失败则返回错误
		if !commonService.Captcha().VerifyString(req.VerifyKey, req.VerifyCode) {
			err = gerror.New("验证码输入错误")
			return
		}
	}

	// 获取客户端IP地址和用户代理
	ip := libUtils.GetClientIp(ctx)
	userAgent := libUtils.GetUserAgent(ctx)

	// 验证用户名和密码，获取用户信息
	user, err = service.SysUser().GetAdminUserByUsernamePassword(ctx, req)
	if err != nil {
		// 保存登录失败的日志信息
		service.SysLoginLog().Invoke(gctx.New(), &model.LoginLogParams{
			Status:    0,            // 登录状态，0表示失败
			Username:  req.Username, // 登录用户名
			Ip:        ip,           // 登录IP地址
			UserAgent: userAgent,    // 用户代理信息
			Msg:       err.Error(),  // 错误信息
			Module:    "系统后台",       // 登录模块
		})
		return
	}

	// 更新用户登录信息（如IP地址）
	err = service.SysUser().UpdateLoginInfo(ctx, user.Id, ip)
	if err != nil {
		return
	}

	// 保存登录成功的日志信息
	service.SysLoginLog().Invoke(gctx.New(), &model.LoginLogParams{
		Status:    1,            // 登录状态，1表示成功
		Username:  req.Username, // 登录用户名
		Ip:        ip,           // 登录IP地址
		UserAgent: userAgent,    // 用户代理信息
		Msg:       "登录成功",       // 登录成功信息
		Module:    "系统后台",       // 登录模块
	})

	// 生成唯一的Token Key，用于Token验证
	key := gconv.String(user.Id) + "-" + gmd5.MustEncryptString(user.UserName) + gmd5.MustEncryptString(user.UserPassword)
	if g.Cfg().MustGet(ctx, "gfToken.multiLogin").Bool() {
		key = gconv.String(user.Id) + "-" + gmd5.MustEncryptString(user.UserName) + gmd5.MustEncryptString(user.UserPassword+ip+userAgent)
	}

	// 隐藏用户密码
	user.UserPassword = ""

	// 生成Token
	token, err = service.GfToken().GenerateToken(ctx, key, user)
	if err != nil {
		g.Log().Error(ctx, err)
		err = gerror.New("登录失败，后端服务出现错误")
		return
	}

	// 获取用户菜单和权限
	menuList, permissions, err = service.SysUser().GetAdminRules(ctx, user.Id)
	if err != nil {
		return
	}

	// 组织返回的响应数据
	res = &system.UserLoginRes{
		UserInfo:    user,        // 用户信息
		Token:       token,       // Token
		MenuList:    menuList,    // 用户菜单
		Permissions: permissions, // 用户权限
	}

	// 保存用户在线状态
	service.SysUserOnline().Invoke(gctx.New(), &model.SysUserOnlineParams{
		UserAgent: userAgent,               // 用户代理
		Uuid:      gmd5.MustEncrypt(token), // 加密的Token
		Token:     token,                   // Token
		Username:  user.UserName,           // 用户名
		Ip:        ip,                      // IP地址
	})
	return
}

// LoginOut 退出登录处理方法
// @param ctx 上下文对象
// @param req 退出登录请求结构体
// @return res 退出登录响应结构体
// @return err 错误信息
func (c *loginController) LoginOut(ctx context.Context, req *system.UserLoginOutReq) (res *system.UserLoginOutRes, err error) {
	// 移除Token
	err = service.GfToken().RemoveToken(ctx, service.GfToken().GetRequestToken(g.RequestFromCtx(ctx)))
	return
}
