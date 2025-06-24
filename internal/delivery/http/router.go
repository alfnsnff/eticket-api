package http

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Router handles HTTP routing with all controllers
type Router struct {
	Engine      *gin.Engine
	Log         logger.Logger
	Controllers *Controllers
	Middlewares *Middlewares
}

type Controllers struct {
	Allocation   *controller.AllocationController
	Auth         *controller.AuthController
	Booking      *controller.BookingController
	Class        *controller.ClassController
	Fare         *controller.FareController
	Harbor       *controller.HarborController
	Manifest     *controller.ManifestController
	Role         *controller.RoleController
	Route        *controller.RouteController
	Schedule     *controller.ScheduleController
	Ship         *controller.ShipController
	Ticket       *controller.TicketController
	User         *controller.UserController
	Payment      *controller.PaymentController
	ClaimSession *controller.ClaimSessionController
}

type Middlewares struct {
	Authenticate *middleware.AuthenticateMiddleware
	Authorize    *middleware.AuthorizeMiddleware
	Log          *middleware.LoggerMiddleware
}

// NewRouter creates a router with all controllers
func NewRouter(engine *gin.Engine, controllers *Controllers, middlewares *Middlewares) *Router {
	return &Router{
		Engine:      engine,
		Controllers: controllers,
		Middlewares: middlewares,
	}
}

// Setup configures all routes
func (r *Router) Setup() {
	r.SetupSystemRoutes()
	r.SetupRoutes()
}

// SetupSystemRoutes registers system-level routes
func (r *Router) SetupSystemRoutes() {
	r.Engine.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	r.Engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "eticket-api",
		})
	})
}

// SetupRoutes registers all API routes
func (r *Router) SetupRoutes() {
	// Get controller references for cleaner code
	ac := r.Controllers.Allocation
	auc := r.Controllers.Auth
	bc := r.Controllers.Booking
	sc := r.Controllers.ClaimSession
	cc := r.Controllers.Class
	fc := r.Controllers.Fare
	hc := r.Controllers.Harbor
	mc := r.Controllers.Manifest
	pc := r.Controllers.Payment
	roc := r.Controllers.Role
	rc := r.Controllers.Route
	scc := r.Controllers.Schedule
	shc := r.Controllers.Ship
	tc := r.Controllers.Ticket
	uc := r.Controllers.User

	// Create route groups
	public := r.Engine.Group("/api/v1")
	protected := r.Engine.Group("/api/v1")
	protected.Use(r.Middlewares.Authenticate.Set())

	// ==================== ALLOCATION ROUTES ====================
	public.GET("/allocations", ac.GetAllAllocations)
	public.GET("/allocation/:id", ac.GetAllocationByID)
	protected.POST("/allocation/create", ac.CreateAllocation)
	protected.PUT("/allocation/update/:id", ac.UpdateAllocation)
	protected.DELETE("/allocation/:id", ac.DeleteAllocation)

	// ==================== AUTH ROUTES ====================
	public.GET("/auth/me", auc.Me)
	public.POST("/auth/login", auc.Login)
	public.POST("/auth/refresh", auc.RefreshToken)
	public.POST("/auth/forget-password", auc.ForgetPassword)
	protected.POST("/auth/logout", auc.Logout)

	// ==================== BOOKING ROUTES ====================
	public.GET("/bookings", bc.GetAllBookings)
	public.GET("/booking/:id", bc.GetBookingByID)
	public.GET("/booking/order/:id", bc.GetBookingByOrderID)
	public.GET("/booking/payment/callback", bc.GetBookingByID)
	protected.POST("/booking/create", bc.CreateBooking)
	protected.PUT("/booking/update/:id", bc.UpdateBooking)
	protected.DELETE("/booking/:id", bc.DeleteBooking)

	// ==================== CLAIM SESSION ROUTES ====================
	public.POST("/session/ticket/lock", sc.CreateClaimSession)
	public.GET("/sessions", sc.GetAllClaimSessions)
	public.GET("/session/:id", sc.GetSessionByID)
	public.POST("/session/ticket/data/entry", sc.UpdateClaimSession)
	public.GET("/session/uuid/:sessionuuid", sc.GetClaimSessionByUUID)
	public.DELETE("/session/:id", sc.DeleteClaimSession)

	// ==================== CLASS ROUTES ====================
	public.GET("/classes", cc.GetAllClasses)
	public.GET("/class/:id", cc.GetClassByID)
	protected.POST("/class/create", cc.CreateClass)
	protected.PUT("/class/update/:id", cc.UpdateClass)
	protected.DELETE("/class/:id", cc.DeleteClass)

	// ==================== FARE ROUTES ====================
	public.GET("/fares", fc.GetAllFares)
	public.GET("/fare/:id", fc.GetFareByID)
	protected.POST("/fare/create", fc.CreateFare)
	protected.PUT("/fare/update/:id", fc.UpdateFare)
	protected.DELETE("/fare/:id", fc.DeleteFare)

	// ==================== HARBOR ROUTES ====================
	public.GET("/harbors", hc.GetAllHarbors)
	public.GET("/harbor/:id", hc.GetHarborByID)
	protected.POST("/harbor/create", hc.CreateHarbor)
	protected.PUT("/harbor/update/:id", hc.UpdateHarbor)
	protected.DELETE("/harbor/:id", hc.DeleteHarbor)

	// ==================== MANIFEST ROUTES ====================
	public.GET("/manifests", mc.GetAllManifests)
	public.GET("/manifest/:id", mc.GetManifestByID)
	protected.POST("/manifest/create", mc.CreateManifest)
	protected.PUT("/manifest/update/:id", mc.UpdateManifest)
	protected.DELETE("/manifest/:id", mc.DeleteManifest)

	// ==================== PAYMENT ROUTES ====================
	public.GET("/payment-channels", pc.GetPaymentChannels)
	public.GET("/payment/transaction/detail/:id", pc.GetTransactionDetail)
	public.POST("/payment/transaction/create", pc.CreatePayment)
	public.POST("/payment/callback", pc.HandleCallback)

	// ==================== ROLE ROUTES ====================
	public.GET("/roles", roc.GetAllRoles)
	public.GET("/role/:id", roc.GetRoleByID)
	protected.POST("/role/create", roc.CreateRole)
	protected.PUT("/role/update/:id", roc.UpdateRole)
	protected.DELETE("/role/:id", roc.DeleteRole)

	// ==================== ROUTE ROUTES ====================
	public.GET("/routes", rc.GetAllRoutes)
	public.GET("/route/:id", rc.GetRouteByID)
	protected.POST("/route/create", rc.CreateRoute)
	protected.PUT("/route/update/:id", rc.UpdateRoute)
	protected.DELETE("/route/:id", rc.DeleteRoute)

	// ==================== SCHEDULE ROUTES ====================
	public.GET("/schedules", scc.GetAllSchedules)
	public.GET("/schedules/active", scc.GetAllScheduled)
	public.GET("/schedule/:id", scc.GetScheduleByID)
	public.GET("/schedule/:id/quota", scc.GetQuotaByScheduleID)
	protected.POST("/schedule/create", scc.CreateSchedule)
	protected.POST("/schedule/allocation/create", scc.CreateScheduleWithAllocation)
	protected.PUT("/schedule/update/:id", scc.UpdateSchedule)
	protected.DELETE("/schedule/:id", scc.DeleteSchedule)

	// ==================== SHIP ROUTES ====================
	public.GET("/ships", shc.GetAllShips)
	public.GET("/ship/:id", shc.GetShipByID)
	protected.POST("/ship/create", shc.CreateShip)
	protected.PUT("/ship/update/:id", shc.UpdateShip)
	protected.DELETE("/ship/:id", shc.DeleteShip)

	// ==================== TICKET ROUTES ====================
	public.GET("/tickets", tc.GetAllTickets)
	public.GET("/ticket/:id", tc.GetTicketByID)
	public.GET("/ticket/schedule/:id", tc.GetAllTicketsByScheduleID)
	protected.PATCH("/ticket/check-in/:id", tc.CheckIn)
	protected.POST("/ticket/create", tc.CreateTicket)
	protected.PUT("/ticket/update/:id", tc.UpdateTicket)
	protected.DELETE("/ticket/:id", tc.DeleteTicket)

	// ==================== USER ROUTES ====================
	public.GET("/users", uc.GetAllUsers)
	public.GET("/user/:id", uc.GetUserByID)
	public.POST("/user/create", uc.CreateUser)
	protected.PUT("/user/update/:id", uc.UpdateUser)
	protected.DELETE("/user/:id", uc.DeleteUser)
}
