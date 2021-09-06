package api

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gorilla/mux"
	"github.com/tsmweb/chat-service/server"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/httputil"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

// HandleWS entry point for chat (websocket).
func HandleWS(jwt auth.JWT, server *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)

		// upgrade connection
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Register incoming connection in server.
		if err = server.Register(userID, conn); err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	})
}

const chatApiVersion string = "v1"

var chatResource string

func init() {
	chatResource = fmt.Sprintf("/%s/ws", chatApiVersion)
}

// MakeChatRouter creates a router for chat.
func MakeChatRouter(
	r *mux.Router,
	jwt auth.JWT,
	auth middleware.Auth,
	server *server.Server) {

	// ws [GET]
	r.Handle(chatResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(HandleWS(jwt, server))),
	).Methods(http.MethodGet)
}
