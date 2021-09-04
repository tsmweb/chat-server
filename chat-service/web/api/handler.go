package api

import (
	"github.com/gobwas/ws"
	"github.com/tsmweb/chat-service/server"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/httputil"
	"log"
	"net/http"
)

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
