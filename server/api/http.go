package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gitchain/gitchain/server"
	"github.com/gitchain/gitchain/ui"
	"github.com/gorilla/mux"
)

var srv *server.T

func Start(srvr *server.T) {
	srv = srvr

	r := mux.NewRouter()

	// Gitchain API
	r.Methods("POST").Path("/rpc").HandlerFunc(jsonRpcService().ServeHTTP)
	r.Methods("GET").Path("/info").HandlerFunc(info)

	setupGitRoutes(r)

	// UI
	r.Methods("GET").Path("/websocket").HandlerFunc(websocketHandler)
	r.Methods("GET").Path("/").HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Content-Type", "text/html")
		var content []byte

		if len(srv.LiveUI) > 0 {
			content, _ = ioutil.ReadFile(path.Join(srv.LiveUI, "index.html"))
		} else {
			content, _ = ui.Asset("index.html")
		}
		resp.Write(content)
	})
	r.Methods("GET").Path("/{path}").HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		file := mux.Vars(req)["path"]
		ext := filepath.Ext(file)
		resp.Header().Add("Content-Type", mime.TypeByExtension(ext))
		var content []byte
		if len(srv.LiveUI) > 0 {
			content, _ = ioutil.ReadFile(path.Join(srv.LiveUI, file))
		} else {
			content, _ = ui.Asset(file)
		}
		if content == nil {
			resp.WriteHeader(404)
			resp.Write([]byte{})
		} else {
			resp.Write(content)
		}
	})

	http.Handle("/", r)

	err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", srv.HttpPort), nil)
	if err != nil {
		log.Fatal(err)
	}
}
