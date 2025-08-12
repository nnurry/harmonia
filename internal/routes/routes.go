package routes

import (
	"net/http"

	"github.com/nnurry/harmonia/internal/handler"
)

type Router struct {
	*http.ServeMux
}

func (router *Router) VirtualMachineHandler() http.Handler {
	mux := http.NewServeMux()

	handler := handler.NewVirtualMachine()

	mux.HandleFunc("POST /create", handler.Create)
	mux.HandleFunc("POST /create/fleet", handler.CreateFleet)

	mux.HandleFunc("POST /format", handler.FormatRequest)

	return mux
}

func (router *Router) V1Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/virtual-machine/", http.StripPrefix("/virtual-machine", router.VirtualMachineHandler()))

	return mux
}

func SetupMux() *Router {
	router := Router{http.NewServeMux()}

	router.ServeMux.Handle("/api/v1/", http.StripPrefix("/api/v1", router.V1Handler()))

	router.ServeMux.HandleFunc("/heartbeat", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte("i have not exploded"))
	})
	return &router
}
