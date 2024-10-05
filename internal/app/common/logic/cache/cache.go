/*
* @desc: 缓存处理
* @company: 云南奇讯科技有限公司
* @Author: yixiaohu <yxh669@qq.com>
* @Date:   2022/9/27 16:33
 */

package cache

import (
	// 引入GoFrame框架的包，用于获取配置和上下文
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	// 引入自定义的缓存处理库
	"github.com/tiger1103/gfast-cache/cache"
	// 引入项目中的常量定义包
	"github.com/tiger1103/gfast/v3/internal/app/common/consts"
	// 引入服务注册包，用于注册缓存服务
	"github.com/tiger1103/gfast/v3/internal/app/common/service"
)

// init 初始化函数，用于注册缓存服务
func init() {
	// 调用服务注册方法，注册缓存服务
	service.RegisterCache(New())
}

// New 创建一个新的缓存服务实例
func New() *sCache {
	var (
		// 创建一个新的上下文对象
		ctx = gctx.New()
		// 定义缓存容器的变量
		cacheContainer *cache.GfCache
	)
	// 从配置中获取缓存的前缀
	prefix := g.Cfg().MustGet(ctx, "system.cache.prefix").String()
	// 从配置中获取缓存的模式（Redis或内存）
	model := g.Cfg().MustGet(ctx, "system.cache.model").String()

	// 如果缓存模式为Redis，则创建Redis缓存实例
	if model == consts.CacheModelRedis {
		// 使用Redis作为缓存容器
		cacheContainer = cache.NewRedis(prefix)
	} else {
		// 否则使用内存作为缓存容器
		cacheContainer = cache.New(prefix)
	}

	// 返回缓存服务实例
	return &sCache{
		GfCache: cacheContainer, // 初始化缓存容器
		prefix:  prefix,         // 设置缓存前缀
	}
}

// sCache 缓存服务的结构体
type sCache struct {
	*cache.GfCache        // 嵌入缓存容器
	prefix         string // 缓存的前缀
}
