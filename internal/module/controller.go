package module

import (
	"eticket-api/internal/delivery/http/controller"
)

type ControllerModule struct {
	AuthController *controller.AuthController
	RoleController *controller.RoleController
	UserController *controller.UserController

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
	PaymentController    *controller.PaymentController
}

func NewControllerModule(uc *UsecaseModule) *ControllerModule {
	return &ControllerModule{
		AuthController: controller.NewAuthController(uc.AuthUsecase),
		RoleController: controller.NewRoleController(uc.RoleUsecase),
		UserController: controller.NewUserController(uc.UserUsecase),

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
		PaymentController:    controller.NewPaymentController(uc.PaymentUsecase),
	}
}
