package router

import (
	"MerchShop/internal/application/core/domain"
	"MerchShop/internal/ports"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const timeout = 200 * time.Millisecond

type Router struct {
	app  ports.APIPort
	srv  http.Server
	port string
}

func NewRouter(app ports.APIPort, port string) *Router {
	r := &Router{app: app, port: port}
	router := gin.Default()
	r.srv = http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}
	authorizedGroup := router.Group("/api")
	authorizedGroup.Use(r.WithAuth())
	{
		authorizedGroup.GET("/info", r.userInfo)
		authorizedGroup.GET("/buy/:item", r.buyItem)
		authorizedGroup.POST("/sendCoin", r.sendCoin)
	}

	unauthorizedGroup := router.Group("/api/auth")
	{
		unauthorizedGroup.POST("/", r.auth)
	}
	return r
}

func (r *Router) Start() error {
	return r.srv.ListenAndServe()
}
func (r *Router) Stop(ctx context.Context) error {
	return r.srv.Shutdown(ctx)
}

func (r *Router) userInfo(ctx *gin.Context) {
	var (
		user domain.User
	)

	val, ok := ctx.Get("user")
	log.Println(val, ok)
	if user, ok = val.(domain.User); !ok {
		log.Println(user, "not okay")
		ctx.JSON(400, gin.H{"error": "internal server error"})
		return
	}
	log.Println(user.Username, ok)
	apictx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	inventory, wallet, err := r.app.Info(apictx, user)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Println(inventory, wallet, err)
	user.Inventory = inventory
	data := ConvertDomainToUserData(user, wallet)
	ctx.JSON(200, data)
}

func (r *Router) buyItem(ctx *gin.Context) {
	var (
		user domain.User
		ok   bool
	)
	val, _ := ctx.Get("user")
	if user, ok = val.(domain.User); !ok {
		ctx.JSON(400, gin.H{"error": "internal server error"})
		return
	}
	item := ctx.Param("item")
	apictx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := r.app.BuyItem(apictx, user, item)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
}

func (r *Router) sendCoin(ctx *gin.Context) {
	var (
		user domain.User
		ok   bool
	)
	val, _ := ctx.Get("user")
	if user, ok = val.(domain.User); !ok {
		ctx.JSON(400, gin.H{"error": "internal server error"})
		return
	}

	info := SendCoinsInfo{}
	decoder := json.NewDecoder(ctx.Request.Body)

	err := decoder.Decode(&info)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer ctx.Request.Body.Close()
	apictx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err = r.app.SendCoin(apictx, user, info.Username, info.Amount)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
}

func (r *Router) auth(ctx *gin.Context) {
	userData := UserCredentials{}
	decoder := json.NewDecoder(ctx.Request.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer ctx.Request.Body.Close()

	apictx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	token, err := r.app.Authorize(apictx, userData.Username, userData.Password)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	response := map[string]string{
		"token": token,
	}
	ctx.JSON(200, response)
}
