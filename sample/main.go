package sample

import (
	"net/http"
	"nicego"
	"fmt"
	"log"
	"time"
	"context"
)

func main() {
	// 获取route句柄
	r := nicego.NewRoute(context.Background())

	r.Static("/", "./public")	// 静态文件服务
	r.From("/hello").Do(HelloCtrl)	// 单一控制器路由注册
	r.From("/news").Use(Logger).Do(NewsCtrl)	// 单中间件路由注册
	r.From("/user").Use(Logger, Auth).Do(UserCtrl)	// 多中间件路由注册

	// 开启http服务
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}
}

// Hello 控制器
func HelloCtrl(ctx context.Context) {
	w, r := nicego.GetMeta(ctx)
	fmt.Fprintf(w, "Hello, world.")
}

// Logger 日志中间件
func Logger(ctx context.Context, next func(context.Context)) {
	t := time.Now()
	_, r := nicego.GetMeta(ctx)

	comFmt := "|---REQ ---[%-4s] %-20s |\n"
	log.Printf(comFmt, meta.Req.Method, meta.Req.URL.Path)

	next(ctx)

	outFmt := "|---RESP---[%-4s] %-20s | %-10s\n"
	log.Printf(outFmt, meta.Req.Method, meta.Req.URL.Path, time.Now().Sub(t))
}

// Auth 中间件
func Auth(ctx context.Context, next func(context.Context)) {
		authrizedUsers := map[string]struct{}{
				"Alice": struct{}{},
				"Bob": struct{}{},
		}
		w, r := nicego.GetMeta(ctx)
		params := r.URL.Query()
		if username := params.Get("username"); username != "" {
				if _, ok := authrizedUsers[username]; ok {
						next(context.WithValue(ctx, nicego.CtxKey("username"), username))
				} else {
						fmt.Fprintf(w, "Unauthrized user: %s.\n", username)
				}
		} else {
				fmt.Frintf(w, "Please provide your name.\n")
		}
}

// NewsCtrl 中间件
func NewsCtrl(ctx context.Context) {
		w, _ := GetMeta(ctx)
		fmt.Fprintf(w, "Welcome to news category.\n")
}

// UserCtrl 中间件
func UserCtrl(ctx context.Context) {
		w, _ := GetMeta(ctx)
		user := ctx.Value(nicego.CtxKey("username").(string)
		fmt.Fprintf(w, "Hello, %s.\n", user)
}