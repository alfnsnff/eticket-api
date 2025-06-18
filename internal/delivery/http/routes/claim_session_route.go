package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type ClaimSessionRoute struct {
	Controller   *controller.ClaimSessionController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewClaimSessionRoute(session_controller *controller.ClaimSessionController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *ClaimSessionRoute {
	return &ClaimSessionRoute{Controller: session_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i ClaimSessionRoute) Set(router *gin.Engine) {
	sc := i.Controller

	public := router.Group("") // No middleware
	public.POST("/session/ticket/lock", sc.CreateClaimSession)
	public.GET("/sessions", sc.GetAllClaimSessions)
	public.GET("/session/:id", sc.GetSessionByID)
	public.POST("/session/ticket/data/entry", sc.UpdateClaimSession)
	public.GET("/session/uuid/:sessionuuid", sc.GetClaimSessionByUUID)
	public.DELETE("/session/:id", sc.DeleteClaimSession)

	protected := router.Group("") // No middleware
	protected.Use(i.Authorized.Handle())
	protected.Use(i.Authorized.Handle())
}
