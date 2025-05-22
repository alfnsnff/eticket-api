package module

import (
	"eticket-api/config"
	"eticket-api/internal/delivery/http/controller"
	authctrl "eticket-api/internal/delivery/http/controller/auth"
	"eticket-api/pkg/jwt"
)

type ControllerModule struct {
	AuthController     *authctrl.AuthController
	RoleController     *authctrl.RoleController
	UserController     *authctrl.UserController
	UserRoleController *authctrl.UserRoleController

	ShipController       *controller.ShipController
	AllocationController *controller.AllocationController
	ManifestController   *controller.ManifestController
	TicketController     *controller.TicketController
	FareController       *controller.FareController
	ScheduleController   *controller.ScheduleController
	BookingController    *controller.BookingController
	SessionController    *controller.SessionController
	RouteController      *controller.RouteController
	HarborController     *controller.HarborController
	ClassController      *controller.ClassController
}

func NewControllerModule(cfg *config.Config, uc *UsecaseModule, tm *jwt.TokenManager) *ControllerModule {
	return &ControllerModule{
		AuthController:     authctrl.NewAuthController(cfg, tm, uc.AuthUsecase),
		RoleController:     authctrl.NewRoleController(uc.RoleUsecase),
		UserController:     authctrl.NewUserController(uc.UserUsecase),
		UserRoleController: authctrl.NewUserRoleController(uc.UserRoleUsecase),

		ShipController:       controller.NewShipController(uc.ShipUsecase),
		AllocationController: controller.NewAllocationController(uc.AllocationUsecase),
		ManifestController:   controller.NewManifestController(uc.ManifestUsecase),
		TicketController:     controller.NewTicketController(uc.TicketUsecase),
		FareController:       controller.NewFareController(uc.FareUsecase),
		ScheduleController:   controller.NewScheduleController(uc.ScheduleUsecase),
		BookingController:    controller.NewBookingController(uc.BookingUsecase),
		SessionController:    controller.NewSessionController(uc.SessionUsecase),
		RouteController:      controller.NewRouteController(uc.RouteUsecase),
		HarborController:     controller.NewHarborController(uc.HarborUsecase),
		ClassController:      controller.NewClassController(uc.ClassUsecase),
	}
}
