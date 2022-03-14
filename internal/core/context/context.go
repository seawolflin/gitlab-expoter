package context

import (
	"flag"
	"github.com/patrickmn/go-cache"
	"github.com/xanzy/go-gitlab"
	"net/url"
	"sync"
	"time"
)

type context struct {
	gitlabUrl                  string
	privateToken               string
	gitlabClient               *gitlab.Client
	cache                      *cache.Cache // 缓存gitlab上的数据，不需要每次都请求
	onCacheEvictedHandlers     map[string]func(c *cache.Cache)
	onCacheEvictedHandlersLock *sync.RWMutex
}

func (ctx *context) Cache() *cache.Cache {
	return ctx.cache
}

func (ctx *context) GitlabClient() *gitlab.Client {
	return ctx.gitlabClient
}

var ins *context
var once = sync.Once{}

func GetInstance() *context {
	once.Do(func() {
		ins = &context{}
	})
	return ins
}

func init() {
	flag.StringVar(&GetInstance().gitlabUrl, "url", "", "Gitlab Url")
	flag.StringVar(&GetInstance().privateToken, "token", "", "Gitlab Private Token")
}

func (ctx *context) Parse() {
	flag.Parse()
	ctx.check()

	client, err := gitlab.NewClient(ins.privateToken, gitlab.WithBaseURL(ins.gitlabUrl))
	if err != nil {
		panic(err.Error())
	}
	ctx.gitlabClient = client

	ctx.cache = cache.New(5*time.Minute, 5*time.Minute)
	ctx.onCacheEvictedHandlers = make(map[string]func(c *cache.Cache), 16)
	ctx.onCacheEvictedHandlersLock = &sync.RWMutex{}
	ctx.cache.OnEvicted(func(key string, _ interface{}) {
		ctx.onCacheEvictedHandlersLock.RLock()
		defer ctx.onCacheEvictedHandlersLock.RUnlock()

		go ctx.onCacheEvictedHandlers[key](ctx.cache)
	})
}

func (ctx context) OnCacheEvicted(key string, f func(c *cache.Cache)) {
	ctx.onCacheEvictedHandlersLock.Lock()
	defer ctx.onCacheEvictedHandlersLock.Unlock()

	ctx.onCacheEvictedHandlers[key] = f
}

func (ctx *context) check() {
	if len(ctx.gitlabUrl) <= 0 {
		panic("url不能为空")
	}
	_, err := url.Parse(ctx.gitlabUrl)
	if err != nil {
		panic("无效的Gitlab url")
	}
	if len(ctx.privateToken) <= 0 {
		panic("token不能为空")
	}
}