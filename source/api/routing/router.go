package routing

import (
	"net/http"

	"github.com/go-chi/chi"
)

var URLParam = chi.URLParam

type Route struct {
	Method  string
	Pattern string
	Handler Handler
}
type Routable interface {
	AuthenticatedRoutes() []Route
	UnauthenticatedRoutes() []Route
}

type RouterBuilder struct {
	unauthRoutes     []Route
	authRoutes       []Route
	notFound         Handler
	methodNotAllowed Handler
	middlewares      []func(http.Handler) http.Handler
	authMiddlewares  []func(http.Handler) http.Handler
}

func NewRouterBuilder() *RouterBuilder {
	return &RouterBuilder{}
}

func (b *RouterBuilder) LoadRoutes(routable Routable) {
	b.unauthRoutes = append(b.unauthRoutes, routable.UnauthenticatedRoutes()...)
	b.authRoutes = append(b.authRoutes, routable.AuthenticatedRoutes()...)
}

func (b *RouterBuilder) SetNotFoundHandler(handler Handler) {
	b.notFound = handler
}

func (b *RouterBuilder) SetMethodNotAllowedHandler(handler Handler) {
	b.methodNotAllowed = handler
}

func (b *RouterBuilder) Build() *chi.Mux {
	router := chi.NewRouter()
	setupRouter(router, b.middlewares, b.unauthRoutes)
	router.Group(func(r chi.Router) {
		setupRouter(r, b.authMiddlewares, b.authRoutes)
	})
	if b.notFound != nil {
		router.NotFound(b.notFound.ServeHTTP)
	}
	if b.methodNotAllowed != nil {
		router.MethodNotAllowed(b.methodNotAllowed.ServeHTTP)
	}
	return router
}

func setupRouter(router chi.Router, middlewares []func(http.Handler) http.Handler, routes []Route) {
	for _, middleware := range middlewares {
		router.Use(middleware)
	}
	for _, route := range routes {
		router.Method(route.Method, route.Pattern, route.Handler)
	}
}

func (b *RouterBuilder) UseMiddleware(middleware ...func(http.Handler) http.Handler) {
	b.middlewares = append(b.middlewares, middleware...)
}

func (b *RouterBuilder) UseAuthMiddleware(middleware ...func(http.Handler) http.Handler) {
	b.authMiddlewares = append(b.authMiddlewares, middleware...)
}
