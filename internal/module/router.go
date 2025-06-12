package module

import (
	"eticket-api/internal/delivery/http/route"
)

type RouterModule struct {
	AuthRouter *route.AuthRouter
	RoleRouter *route.RoleRouter
	UserRouter *route.UserRouter

	ShipRouter       *route.ShipRouter
	AllocationRouter *route.AllocationRouter
	ManifestRouter   *route.ManifestRouter
	TicketRouter     *route.TicketRouter
	FareRouter       *route.FareRouter
	ScheduleRouter   *route.ScheduleRouter
	BookingRouter    *route.BookingRouter
	SessionRouter    *route.SessionRouter
	RouteRouter      *route.RouteRouter
	HarborRouter     *route.HarborRouter
	ClassRouter      *route.ClassRouter
	PaymentRouter    *route.PaymentRouter
}

func NewRouteModule(uc *ControllerModule, m *MiddlewareModule) *RouterModule {
	return &RouterModule{
		AuthRouter: route.NewAuthRouter(uc.AuthController, m.Authenticate, m.Authorize),
		RoleRouter: route.NewRoleRouter(uc.RoleController, m.Authenticate, m.Authorize),
		UserRouter: route.NewUserRouter(uc.UserController, m.Authenticate, m.Authorize),

		ShipRouter:       route.NewShipRouter(uc.ShipController, m.Authenticate, m.Authorize),
		AllocationRouter: route.NewAllocationRouter(uc.AllocationController, m.Authenticate, m.Authorize),
		ManifestRouter:   route.NewManifestRouter(uc.ManifestController, m.Authenticate, m.Authorize),
		TicketRouter:     route.NewTicketRouter(uc.TicketController, m.Authenticate, m.Authorize),
		FareRouter:       route.NewFareRouter(uc.FareController, m.Authenticate, m.Authorize),
		ScheduleRouter:   route.NewScheduleRouter(uc.ScheduleController, m.Authenticate, m.Authorize),
		BookingRouter:    route.NewBookingRouter(uc.BookingController, m.Authenticate, m.Authorize),
		SessionRouter:    route.NewSessionRouter(uc.SessionController, m.Authenticate, m.Authorize),
		RouteRouter:      route.NewRouteRouter(uc.RouteController, m.Authenticate, m.Authorize),
		HarborRouter:     route.NewHarborRouter(uc.HarborController, m.Authenticate, m.Authorize),
		ClassRouter:      route.NewClassRouter(uc.ClassController, m.Authenticate, m.Authorize),
		PaymentRouter:    route.NewPaymentRouter(uc.PaymentController, m.Authenticate, m.Authorize),
	}
}
