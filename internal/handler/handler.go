// Package handler provides
package handler

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"tests_app/internal/config"
	"tests_app/internal/service"
	"time"

	"github.com/gin-gonic/gin"

	_ "tests_app/docs"

	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Handler struct {
	services *service.Service
	config   config.Config
}

func New(services *service.Service, config config.Config) *Handler {
	return &Handler{services, config}
}

func (h *Handler) InitApiRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := router.Group("api")
	{
		apiAuth := api.Group("auth")
		{
			apiAuth.POST("sign-up", h.signUp)
			apiAuth.POST("sign-in", h.signIn)
			apiAuth.POST("sign-out", h.signOut)
			apiAuth.POST("refresh-token", h.refreshToken)
			apiAuth.PUT("update-user", h.userIdentity, h.updateUser)
		}

		apiTest := api.Group("test")
		{
			apiTest.GET("", h.getAllTests)
			apiTest.GET(":id", h.getTestById)

			apiTest.GET(":id/full", h.userIdentity, h.getFullTestById)
			apiTest.POST("", h.userIdentity, h.createTest)
			apiTest.PUT(":id", h.userIdentity, h.updateTest)
			apiTest.DELETE(":id", h.userIdentity, h.deleteTest)
		}

		apiQuestion := api.Group("question", h.userIdentity)
		{
			apiQuestion.GET(":id", h.getQuestionById)
			apiQuestion.GET("/next", h.getNextQuestion)
			apiQuestion.POST("", h.createQuestion)
			apiQuestion.PUT(":id", h.updateQuestion)
			apiQuestion.DELETE(":id", h.deleteQuestion)
		}

		apiQuestionResult := api.Group("answer", h.userIdentity)
		{
			apiQuestionResult.GET("", h.getAnswer)
			apiQuestionResult.POST("", h.setAnswer)
		}

		apiTestResult := api.Group("test-answer", h.userIdentity)
		{
			apiTestResult.POST("", h.completeTest)
		}
	}

	return router
}

func (h *Handler) InitFileRoutes(c *gin.Engine) {
	c.GET("/", func(ctx *gin.Context) {
		files.RLock()
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", files.mainPage)
		files.RUnlock()
	})
	c.GET("/vue_app.js", func(ctx *gin.Context) {
		files.RLock()
		ctx.Data(http.StatusOK, "application/javascript; charset=utf-8", files.js)
		files.RUnlock()
	})
	c.GET("/main.css", func(ctx *gin.Context) {
		files.RLock()
		ctx.Data(http.StatusOK, "text/css; charset=utf-8", files.css)
		files.RUnlock()
	})

	loadFiles()
	go watchFile("static/index.html")
	go watchFile("static/vue_app.js")
	go watchFile("static/main.css")
}

var files = struct {
	sync.RWMutex
	mainPage []byte
	js       []byte
	css      []byte
}{}

func loadFiles() {
	files.Lock()
	defer files.Unlock()

	list := []struct {
		file *[]byte
		path string
	}{
		{
			file: &files.mainPage,
			path: "static/index.html",
		},
		{
			file: &files.js,
			path: "static/vue_app.js",
		},
		{
			file: &files.css,
			path: "static/main.css",
		},
	}

	for _, x := range list {
		file, err := os.Open(x.path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		*x.file, err = io.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
}

func watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			loadFiles()
			initialStat = stat
			log.Print("files were changed")
		}

		time.Sleep(3 * time.Second)
	}
}
