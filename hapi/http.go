package hapi

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hnit-acm/hfunc/hserver/hhttp"
	"golang.org/x/sync/errgroup"
)

// Server 启动gin hserver
// port 端口
// regFunc 注册的路由函数
func Serve(port string, g *gin.Engine, regFunc func(c *gin.Engine)) {
	if g == nil {
		g = gin.Default()
	}
	if regFunc != nil {
		regFunc(g)
	}
	httpServer := hhttp.NewServer("tcp", ":"+port, hhttp.WithHandler(g))

	ctx, cancel := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)
	sigs := []os.Signal{
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
	}

	group.Go(func() error {
		startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Start(startCtx)
	})

	group.Go(func() error {
		<-ctx.Done() // 等待退出信号
		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		log.Println("stop hserver")
		return httpServer.Stop(stopCtx)
	})

	c := make(chan os.Signal, len(sigs))
	signal.Notify(c, sigs...)
	group.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case sig := <-c:
				switch sig {
				case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
					// 可以处理一些退出逻辑
					cancel()
				default:
				}
			}
		}
	})
	group.Wait()
}

// Server 启动gin hserver
// port 端口
// regFunc 注册的路由函数
// clearFunc server退出逻辑处理
// timeOut 超时时间
func ServeTimeout(port string, g *gin.Engine, regFunc func(c *gin.Engine), clearFunc func(), timeOut time.Duration, opts ...hhttp.Option) {
	if g == nil {
		g = gin.Default()
	}
	if regFunc != nil {
		regFunc(g)
	}
	opts = append(opts, hhttp.WithHandler(g))
	httpServer := hhttp.NewServer("tcp", ":"+port, opts...)

	ctx, cancel := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)
	sigs := []os.Signal{
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
	}

	group.Go(func() error {
		startCtx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()
		return httpServer.Start(startCtx)
	})

	group.Go(func() error {
		<-ctx.Done() // 等待退出信号
		stopCtx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()
		log.Println("stop hserver")
		return httpServer.Stop(stopCtx)
	})

	c := make(chan os.Signal, len(sigs))
	signal.Notify(c, sigs...)
	group.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case sig := <-c:
				switch sig {
				case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
					// 可以处理一些退出逻辑
					clearFunc()
					cancel()
				default:
				}
			}
		}
	})
	group.Wait()
}
