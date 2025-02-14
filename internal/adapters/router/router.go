package router

import (
	"MerchShop/internal/application/core/domain"
	"MerchShop/internal/ports"
	"encoding/json"
	"fmt"
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
		Addr:    fmt.Sprintf(":%s", port),
		Handler: http.DefaultServeMux,
	}
	http.Handle("GET /api/info", r.WithAuth(r.UserInfo))
	http.HandleFunc("GET /api/buy/{item}", r.WithAuth(r.BuyItem))
	http.HandleFunc("POST /api/sendCoin", r.WithAuth(r.sendCoin))
	http.HandleFunc("POST /api/auth", r.Auth)
	//r.srv.Handler = Logger(r.srv.Handler)
	return r
}

func (r *Router) Start() error {
	return r.srv.ListenAndServe()
}

func (r *Router) UserInfo(w http.ResponseWriter, req *http.Request) {
	var (
		user domain.User
		ok   bool
	)
	if user, ok = req.Context().Value("user").(domain.User); !ok {
		WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("internal server error"))
		return
	}
	inventory, wallet, err := r.app.Info(user)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	user.Inventory = inventory
	data := ConvertDomainToUserData(user, wallet)
	json.NewEncoder(w).Encode(data)
}

func (r *Router) BuyItem(w http.ResponseWriter, req *http.Request) {
	var (
		user domain.User
		ok   bool
	)
	if user, ok = req.Context().Value("user").(domain.User); !ok {
		WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("internal server error"))
		return
	}
	item := req.PathValue("item")
	err := r.app.BuyItem(user, item)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
}

func (r *Router) sendCoin(w http.ResponseWriter, req *http.Request) {
	var (
		user domain.User
		ok   bool
	)
	if user, ok = req.Context().Value("user").(domain.User); !ok {
		WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("internal server error"))
		return
	}

	info := SendCoinsInfo{}
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&info)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	_, err = r.app.SendCoin(user, info.Username, info.Amount)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
}

func (r *Router) Auth(w http.ResponseWriter, req *http.Request) {
	userData := UserCredentials{}
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&userData)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	defer req.Body.Close()

	token, err := r.app.Authorize(userData.Username, userData.Password)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	response := map[string]string{
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
