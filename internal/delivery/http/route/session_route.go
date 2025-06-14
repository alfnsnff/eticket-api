package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type SessionRouter struct {
	Controller   *controller.SessionController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewSessionRouter(session_controller *controller.SessionController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *SessionRouter {
	return &SessionRouter{Controller: session_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i SessionRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {
	sc := i.Controller

	public := rg.Group("") // No middleware
	public.POST("/session/ticket/lock", sc.SessionTicketLock)
	public.POST("/session/ticket/data/entry", sc.TicketDataEntry)
	public.GET("/sessions", sc.GetAllSessions)
	public.GET("/session/:id", sc.GetSessionByID)
	public.GET("/session/uuid/:sessionid", sc.GetSessionBySessionID)
	public.POST("/session/create", sc.CreateSession)
	public.PUT("/session/update/:id", sc.UpdateSession)
	public.DELETE("/session/:id", sc.DeleteSession)

	protected := rg.Group("") // No middleware
	protected.Use(i.Authorized.Handle())
	protected.Use(i.Authorized.Handle())
}
