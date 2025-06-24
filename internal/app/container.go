package app

import (
	"eticket-api/internal/client"
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase/allocation"
	"eticket-api/internal/usecase/auth"
	"eticket-api/internal/usecase/booking"
	"eticket-api/internal/usecase/claim_session"
	"eticket-api/internal/usecase/class"
	"eticket-api/internal/usecase/fare"
	"eticket-api/internal/usecase/harbor"
	"eticket-api/internal/usecase/manifest"
	"eticket-api/internal/usecase/payment"
	"eticket-api/internal/usecase/role"
	"eticket-api/internal/usecase/route"
	"eticket-api/internal/usecase/schedule"
	"eticket-api/internal/usecase/ship"
	"eticket-api/internal/usecase/ticket"
	"eticket-api/internal/usecase/user"
)

// DependencyContainer manages all application dependencies efficiently
type Container struct {
	Bootstrap    *Bootstrap
	Clients      *Clients
	Repositories *Repositories
	UseCases     *UseCases
	Controllers  *Controllers
	Middlewares  *Middlewares
}

type Clients struct {
	Tripay *client.TripayClient
}

// RepositoryContainer holds all repositories
type Repositories struct {
	Allocation *repository.AllocationRepository
	Auth       *repository.AuthRepository
	Booking    *repository.BookingRepository
	Class      *repository.ClassRepository
	Fare       *repository.FareRepository
	Harbor     *repository.HarborRepository
	Manifest   *repository.ManifestRepository
	Role       *repository.RoleRepository
	Route      *repository.RouteRepository
	Schedule   *repository.ScheduleRepository
	Ship       *repository.ShipRepository
	Ticket     *repository.TicketRepository
	User       *repository.UserRepository
	Session    *repository.SessionRepository
}

// UseCaseContainer holds all use cases
type UseCases struct {
	Allocation   *allocation.AllocationUsecase
	Auth         *auth.AuthUsecase
	Booking      *booking.BookingUsecase
	Class        *class.ClassUsecase
	ClaimSession *claim_session.ClaimSessionUsecase
	Fare         *fare.FareUsecase
	Harbor       *harbor.HarborUsecase
	Manifest     *manifest.ManifestUsecase
	Payment      *payment.PaymentUsecase
	Role         *role.RoleUsecase
	Route        *route.RouteUsecase
	Schedule     *schedule.ScheduleUsecase
	Ship         *ship.ShipUsecase
	Ticket       *ticket.TicketUsecase
	User         *user.UserUsecase
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

// MiddlewareContainer holds all middlewares
type Middlewares struct {
	Authenticate *middleware.AuthenticateMiddleware
	Authorize    *middleware.AuthorizeMiddleware
}

// NewDependencyContainer creates a new dependency container
func NewContainer(bootstrap *Bootstrap) *Container {
	return &Container{
		Bootstrap:    bootstrap,
		Repositories: &Repositories{},
		UseCases:     &UseCases{},
		Middlewares:  &Middlewares{},
	}
}

// InitializeRepositories initializes all repositories
func (i *Container) InitializeRepositories() {
	i.Repositories.Allocation = repository.NewAllocationRepository()
	i.Repositories.Auth = repository.NewAuthRepository()
	i.Repositories.Booking = repository.NewBookingRepository()
	i.Repositories.Class = repository.NewClassRepository()
	i.Repositories.Fare = repository.NewFareRepository()
	i.Repositories.Harbor = repository.NewHarborRepository()
	i.Repositories.Manifest = repository.NewManifestRepository()
	i.Repositories.Role = repository.NewRoleRepository()
	i.Repositories.Route = repository.NewRouteRepository()
	i.Repositories.Schedule = repository.NewScheduleRepository()
	i.Repositories.Session = repository.NewSessionRepository()
	i.Repositories.Ship = repository.NewShipRepository()
	i.Repositories.Ticket = repository.NewTicketRepository()
	i.Repositories.User = repository.NewUserRepository()
}

// InitializeUseCases initializes all use cases with their dependencies
func (i *Container) InitializeUseCases() {
	i.Clients.Tripay = client.NewTripayClient(
		i.Bootstrap.Client,
		&i.Bootstrap.Config.Tripay)
	i.UseCases.Allocation = allocation.NewAllocationUsecase(
		i.Bootstrap.DB,
		i.Repositories.Allocation,
		i.Repositories.Fare,
	)
	i.UseCases.Auth = auth.NewAuthUsecase(
		i.Bootstrap.DB,
		i.Repositories.Auth,
		i.Repositories.User,
		i.Bootstrap.Mailer,
		i.Bootstrap.Token,
	)
	i.UseCases.Booking = booking.NewBookingUsecase(
		i.Bootstrap.DB,
		i.Repositories.Booking,
	)
	i.UseCases.Class = class.NewClassUsecase(
		i.Bootstrap.DB,
		i.Repositories.Class,
	)
	i.UseCases.ClaimSession = claim_session.NewClaimSessionUsecase(
		i.Bootstrap.DB,
		i.Repositories.Session,
		i.Repositories.Ticket,
		i.Repositories.Schedule,
		i.Repositories.Allocation,
		i.Repositories.Manifest,
		i.Repositories.Fare,
		i.Repositories.Booking,
	)
	i.UseCases.Fare = fare.NewFareUsecase(
		i.Bootstrap.DB,
		i.Repositories.Fare,
	)
	i.UseCases.Harbor = harbor.NewHarborUsecase(
		i.Bootstrap.DB,
		i.Repositories.Harbor,
	)
	i.UseCases.Manifest = manifest.NewManifestUsecase(
		i.Bootstrap.DB,
		i.Repositories.Manifest,
	)
	i.UseCases.Payment = payment.NewPaymentUsecase(
		i.Bootstrap.DB,
		i.Clients.Tripay,
		i.Repositories.Session,
		i.Repositories.Booking,
		i.Repositories.Ticket,
		i.Bootstrap.Mailer,
	)
	i.UseCases.Role = role.NewRoleUsecase(
		i.Bootstrap.DB,
		i.Repositories.Role,
	)
	i.UseCases.Route = route.NewRouteUsecase(
		i.Bootstrap.DB,
		i.Repositories.Route,
	)
	i.UseCases.Schedule = schedule.NewScheduleUsecase(
		i.Bootstrap.DB,
		i.Repositories.Allocation,
		i.Repositories.Class,
		i.Repositories.Fare,
		i.Repositories.Manifest,
		i.Repositories.Ship,
		i.Repositories.Schedule,
		i.Repositories.Ticket,
	)
	i.UseCases.Ship = ship.NewShipUsecase(
		i.Bootstrap.DB,
		i.Repositories.Ship,
	)
	i.UseCases.Ticket = ticket.NewTicketUsecase(
		i.Bootstrap.DB,
		i.Repositories.Ticket,
		i.Repositories.Schedule,
		i.Repositories.Manifest,
		i.Repositories.Fare,
	)
	i.UseCases.User = user.NewUserUsecase(
		i.Bootstrap.DB,
		i.Repositories.User,
	)
}

// RegisterControllersToRouter registers all controllers to the router
func (i *Container) InitializeControllers() *Container {
	i.Controllers.Allocation = controller.NewAllocationController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Allocation,
	)
	i.Controllers.Auth = controller.NewAuthController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Auth,
	)
	i.Controllers.Booking = controller.NewBookingController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Booking,
	)
	i.Controllers.Class = controller.NewClassController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Class,
	)
	i.Controllers.ClaimSession = controller.NewClaimSessionController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.ClaimSession,
	)
	i.Controllers.Fare = controller.NewFareController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Fare,
	)
	i.Controllers.Harbor = controller.NewHarborController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Harbor,
	)
	i.Controllers.Manifest = controller.NewManifestController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Manifest,
	)
	i.Controllers.Payment = controller.NewPaymentController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Payment,
	)
	i.Controllers.Role = controller.NewRoleController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Role,
	)
	i.Controllers.Route = controller.NewRouteController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Route,
	)
	i.Controllers.Schedule = controller.NewScheduleController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Schedule,
	)
	i.Controllers.Ship = controller.NewShipController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Ship,
	)
	i.Controllers.Ticket = controller.NewTicketController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.Ticket,
	)
	i.Controllers.User = controller.NewUserController(
		i.Bootstrap.Log,
		i.Bootstrap.Validate,
		i.UseCases.User,
	)

	return i
}
