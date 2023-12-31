package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"sync"
	"todo/internal/manager/interfaces/processor/rest"
)

type Handler struct {
	engine       *gin.Engine
	engineRunner sync.Once

	ctl rest.IControllersManager
}

func InitHandler(controller rest.IControllersManager) *Handler {
	newHandler := &Handler{
		ctl: controller,
	}
	newHandler.InitRoutes()
	return newHandler
}

func (h *Handler) InitRoutes() {
	h.Engine().GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "ok"})
	})

	h.Engine().GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h.AddRoutesV1()
}

func (h *Handler) Engine() *gin.Engine {
	h.engineRunner.Do(func() {
		h.engine = InitEngine(viper.GetString("ENV_MODE"))
	})
	return h.engine
}

func InitEngine(env string) *gin.Engine {
	setGinMode(env)

	router := gin.New()
	router.Use(gin.Recovery())
	router.RedirectTrailingSlash = true
	corsConfig := cors.Config{
		AllowOriginFunc: func(origin string) bool { return true },
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           604800,
	}
	router.Use(cors.New(corsConfig))

	return router
}

func setGinMode(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
		return
	}
	gin.SetMode(gin.DebugMode)
}
