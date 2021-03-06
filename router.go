// ===============================================
// Description  : nciego/router.go
// Author       : StevE.Z
// Email        : stevzhang01@gmail.com
// Date         : 2019-04-22 10:32:05
// ================================================
package nicego

import (
	"context"
	"net/http"
)

// router
type Router struct {
	route       *Route
	pattern     string
	middlewares []func(context.Context, func(context.Context))
}

// injectMiddlewares
func injectMiddlewares(rtr *Router, controller func(context.Context)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			next func(context.Context)
			i    int
		)
		next = func(ctx context.Context) {
			if i < len(rtr.middlewares) {
				i++
				rtr.middlewares[i-1](ctx, next)
			} else {
				if controller != nil {
					controller(ctx)
				}
			}
		}
		ctx := context.WithValue(rtr.route.ctx, metaKey{}, metaVal{w: w, r: r})
		next(ctx)
	}
}

// Use
func (rtr *Router) Use(middlewares ...func(context.Context, func(context.Context))) *Router {
	rtr.middlewares = append(rtr.middlewares, middlewares...)
	return rtr
}

// Do
func (rtr *Router) Do(controller func(context.Context)) {
	rtr.route.mux.HandleFunc(rtr.pattern, injectMiddlewares(rtr, controller))
}

// Static
func (rtr *Router) Static(dir string) {
	fileHandler := http.StripPrefix(rtr.pattern, http.FileServer(http.Dir(dir)))
	staticContrller := func(ctx context.Context) {
		w, r := GetMeta(ctx)
		fileHandler.ServeHTTP(w, r)
	}
	rtr.route.mux.HandleFunc(rtr.pattern, injectMiddlewares(rtr, staticContrller))
}
