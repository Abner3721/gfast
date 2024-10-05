/*
* @desc: 登录相关的请求与响应结构体定义
* @company: 云南奇讯科技有限公司
* @Author: yixiaohu
* @Date:   2022/4/27 21:51
 */

package system

import (
	// 导入GoFrame框架，用于声明元信息
	"github.com/gogf/gf/v2/frame/g"
	// 导入公共API接口定义
	commonApi "github.com/tiger1103/gfast/v3/api/v1/common"
	// 导入系统模型，包含用户信息、菜单等结构体
	"github.com/tiger1103/gfast/v3/internal/app/system/model"
)

// UserLoginReq 定义用户登录请求的结构体
type UserLoginReq struct {
	// 定义API的元信息，包括路径、标签、HTTP方法和摘要
	g.Meta `path:"/login" tags:"登录" method:"post" summary:"用户登录"`
	// 用户名字段，使用参数绑定和验证，不能为空
	Username string `p:"username" v:"required#用户名不能为空"`
	// 密码字段，使用参数绑定和验证，不能为空
	Password string `p:"password" v:"required#密码不能为空"`
	// 验证码字段，使用参数绑定和验证，不能为空
	VerifyCode string `p:"verifyCode" v:"required#验证码不能为空"`
	// 验证码对应的key，用于验证用户输入的验证码
	VerifyKey string `p:"verifyKey"`
}

// UserLoginRes 定义用户登录响应的结构体
type UserLoginRes struct {
	// 响应的元信息，定义返回内容的MIME类型为JSON
	g.Meta `mime:"application/json"`
	// 登录成功后的用户信息，引用外部的LoginUserRes结构体
	UserInfo *model.LoginUserRes `json:"userInfo"`
	// 登录成功后的Token，用于后续的身份验证
	Token string `json:"token"`
	// 登录成功后返回的用户菜单列表
	MenuList []*model.UserMenus `json:"menuList"`
	// 登录成功后返回的用户权限列表
	Permissions []string `json:"permissions"`
}

// UserLoginOutReq 定义用户退出登录请求的结构体
type UserLoginOutReq struct {
	// 定义API的元信息，包括路径、标签、HTTP方法和摘要
	g.Meta `path:"/logout" tags:"登录" method:"get" summary:"退出登录"`
	// 引用公共API接口中的Author结构体，用于获取当前登录用户信息
	commonApi.Author
}

// UserLoginOutRes 定义用户退出登录响应的结构体
type UserLoginOutRes struct {
	// 空的响应结构体，表示退出登录时不返回具体内容
}
