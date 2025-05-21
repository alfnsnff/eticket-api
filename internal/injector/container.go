package injector

import (
	"eticket-api/config"
	"eticket-api/internal/delivery/http/controller"
	authcontroller "eticket-api/internal/delivery/http/controller/auth"

	"eticket-api/internal/repository"
	authrepository "eticket-api/internal/repository/auth"

	"eticket-api/internal/usecase"
	authusecase "eticket-api/internal/usecase/auth"

	"eticket-api/pkg/jwt"

	"gorm.io/gorm"
)

type Container struct {
	Cfg          *config.Config
	DB           *gorm.DB
	TokenManager *jwt.TokenManager

	AuthRepository     *authrepository.AuthRepository
	RoleRepository     *authrepository.RoleRepository
	UserRepository     *authrepository.UserRepository
	UserRoleRepository *authrepository.UserRoleRepository

	ShipRepository       *repository.ShipRepository
	AllocationRepository *repository.AllocationRepository
	ManifestRepository   *repository.ManifestRepository
	TicketRepository     *repository.TicketRepository
	FareRepository       *repository.FareRepository
	ScheduleRepository   *repository.ScheduleRepository
	BookingRepository    *repository.BookingRepository
	SessionRepository    *repository.SessionRepository
	RouteRepository      *repository.RouteRepository
	HarborRepository     *repository.HarborRepository
	ClassRepository      *repository.ClassRepository

	AuthUsecase     *authusecase.AuthUsecase
	RoleUsecase     *authusecase.RoleUsecase
	UserUsecase     *authusecase.UserUsecase
	UserRoleUsecase *authusecase.UserRoleUsecase

	ShipUsecase       *usecase.ShipUsecase
	AllocationUsecase *usecase.AllocationUsecase
	ManifestUsecase   *usecase.ManifestUsecase
	TicketUsecase     *usecase.TicketUsecase
	FareUsecase       *usecase.FareUsecase
	ScheduleUsecase   *usecase.ScheduleUsecase
	BookingUsecase    *usecase.BookingUsecase
	SessionUsecase    *usecase.SessionUsecase
	RouteUsecase      *usecase.RouteUsecase
	HarborUsecase     *usecase.HarborUsecase
	ClassUsecase      *usecase.ClassUsecase

	AuthController     *authcontroller.AuthController
	RoleController     *authcontroller.RoleController
	UserController     *authcontroller.UserController
	UserRoleController *authcontroller.UserRoleController

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

func NewContainer(
	cfg *config.Config,
	db *gorm.DB,
	tm *jwt.TokenManager,

	// Auth
	ar *authrepository.AuthRepository,
	rr *authrepository.RoleRepository,
	ur *authrepository.UserRepository,
	urr *authrepository.UserRoleRepository,

	// App repos
	sr *repository.ShipRepository,
	ar2 *repository.AllocationRepository,
	mr *repository.ManifestRepository,
	tr *repository.TicketRepository,
	fr *repository.FareRepository,
	schr *repository.ScheduleRepository,
	br *repository.BookingRepository,
	sessr *repository.SessionRepository,
	rr2 *repository.RouteRepository,
	hr *repository.HarborRepository,
	cr *repository.ClassRepository,

	// Usecases
	au *authusecase.AuthUsecase,
	ru *authusecase.RoleUsecase,
	uu *authusecase.UserUsecase,
	uru *authusecase.UserRoleUsecase,

	su *usecase.ShipUsecase,
	alu *usecase.AllocationUsecase,
	mu *usecase.ManifestUsecase,
	tu *usecase.TicketUsecase,
	fu *usecase.FareUsecase,
	schu *usecase.ScheduleUsecase,
	bu *usecase.BookingUsecase,
	seu *usecase.SessionUsecase,
	rou *usecase.RouteUsecase,
	hu *usecase.HarborUsecase,
	cu *usecase.ClassUsecase,

	// Controllers
	ac *authcontroller.AuthController,
	rc *authcontroller.RoleController,
	uc *authcontroller.UserController,
	urc *authcontroller.UserRoleController,

	sc *controller.ShipController,
	alc *controller.AllocationController,
	mc *controller.ManifestController,
	tc *controller.TicketController,
	fc *controller.FareController,
	scc *controller.ScheduleController,
	bc *controller.BookingController,
	sec *controller.SessionController,
	roc *controller.RouteController,
	hc *controller.HarborController,
	cc *controller.ClassController,
) *Container {
	return &Container{
		Cfg:          cfg,
		DB:           db,
		TokenManager: tm,

		AuthRepository:     ar,
		RoleRepository:     rr,
		UserRepository:     ur,
		UserRoleRepository: urr,

		ShipRepository:       sr,
		AllocationRepository: ar2,
		ManifestRepository:   mr,
		TicketRepository:     tr,
		FareRepository:       fr,
		ScheduleRepository:   schr,
		BookingRepository:    br,
		SessionRepository:    sessr,
		RouteRepository:      rr2,
		HarborRepository:     hr,
		ClassRepository:      cr,

		AuthUsecase:     au,
		RoleUsecase:     ru,
		UserUsecase:     uu,
		UserRoleUsecase: uru,

		ShipUsecase:       su,
		AllocationUsecase: alu,
		ManifestUsecase:   mu,
		TicketUsecase:     tu,
		FareUsecase:       fu,
		ScheduleUsecase:   schu,
		BookingUsecase:    bu,
		SessionUsecase:    seu,
		RouteUsecase:      rou,
		HarborUsecase:     hu,
		ClassUsecase:      cu,

		AuthController:     ac,
		RoleController:     rc,
		UserController:     uc,
		UserRoleController: urc,

		ShipController:       sc,
		AllocationController: alc,
		ManifestController:   mc,
		TicketController:     tc,
		FareController:       fc,
		ScheduleController:   scc,
		BookingController:    bc,
		SessionController:    sec,
		RouteController:      roc,
		HarborController:     hc,
		ClassController:      cc,
	}
}
