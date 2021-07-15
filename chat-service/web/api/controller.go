package api

import (
	"github.com/gobwas/ws"
	"github.com/tsmweb/chat-service/chat"
	"github.com/tsmweb/go-helper-api/auth"
	ctrl "github.com/tsmweb/go-helper-api/controller"
	"log"
	"net/http"
)

// Controller provides the end point for the routers.
type Controller struct {
	*ctrl.Controller
	chat *chat.Chat
}

// NewController creates a new instance of Controller.
func NewController(jwt auth.JWT, chat *chat.Chat) *Controller {
	return &Controller{
		Controller: ctrl.New(jwt),
		chat: chat,
	}
}

// Connect 
func (c *Controller) Connect() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		// upgrade connection
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Register incoming connection in chat.
		if err = c.chat.Register(userID, conn); err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	})
}
