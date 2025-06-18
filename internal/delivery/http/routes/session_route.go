package routes

// import (
// 	"eticket-api/internal/delivery/http/controller"
// 	"eticket-api/internal/delivery/http/middleware"

// 	"github.com/gin-gonic/gin"
// )

// type SessionRoute struct {
// 	Controller   *controller.SessionController
// 	Authenticate *middleware.AuthenticateMiddleware
// 	Authorized   *middleware.AuthorizeMiddleware
// }

// func NewSessionRoute(session_controller *controller.SessionController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *SessionRoute {
// 	return &SessionRoute{Controller: session_controller, Authenticate: authtenticate, Authorized: authorized}
// }

// func (i SessionRoute) Set(router *gin.Engine) {
// 	sc := i.Controller

// 	public := router.Group("") // No middleware
// 	public.POST("/session/ticket/lock", sc.SessionTicketLock)
// 	public.POST("/session/ticket/data/entry", sc.SessionTicketDataEntry)
// 	public.GET("/sessions", sc.GetAllSessions)
// 	public.GET("/session/:id", sc.GetSessionByID)
// 	public.GET("/session/uuid/:sessionid", sc.GetSessionBySessionID)
// 	public.POST("/session/create", sc.CreateSession)
// 	public.PUT("/session/update/:id", sc.UpdateSession)
// 	public.DELETE("/session/:id", sc.DeleteSession)

// 	protected := router.Group("") // No middleware
// 	protected.Use(i.Authorized.Handle())
// 	protected.Use(i.Authorized.Handle())
// }
