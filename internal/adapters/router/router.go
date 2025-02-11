package router

import (
	"MerchShop/internal/ports"
	"encoding/json"
	"net/http"
)

type Router struct {
	app  ports.APIPort
	srv  http.Server
	port string
}

func NewRouter(app ports.APIPort, port string) *Router {
	r := &Router{app: app, port: port}
	r.srv = http.Server{
		Handler: http.DefaultServeMux,
	}
	http.HandleFunc("GET /api/info", r.UserInfo)
	http.Handle("", nil)
	http.Handle("", nil)
	http.Handle("", nil)
	return r
}

func (r *Router) Start() error {
	return r.srv.ListenAndServe()
}

func (r *Router) UserInfo(w http.ResponseWriter, req *http.Request) {
	token := req.FormValue("access_token")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := r.app.Authenticate(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}
	opps, err := r.app.Info(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	data := ConvertToUserData(user, opps)
	json.NewEncoder(w).Encode(data)
}
