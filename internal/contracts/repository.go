package contracts

import (
	"eticket-api/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	ScheduleRepository interface {
		Create(db *gorm.DB, entity *entity.Schedule) error
		Update(db *gorm.DB, entity *entity.Schedule) error
		Delete(db *gorm.DB, entity *entity.Schedule) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Schedule, error)
		GetByID(db *gorm.DB, id uint) (*entity.Schedule, error)
		GetAllScheduled(db *gorm.DB) ([]*entity.Schedule, error)
		GetActiveSchedule(db *gorm.DB) ([]*entity.Schedule, error)
	}

	TicketRepository interface {
		Create(db *gorm.DB, entity *entity.Ticket) error
		Update(db *gorm.DB, entity *entity.Ticket) error
		Delete(db *gorm.DB, entity *entity.Ticket) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Ticket, error)
		GetByScheduleID(db *gorm.DB, id, limit, offset int, sort, search string) ([]*entity.Ticket, error)
		GetByID(db *gorm.DB, id uint) (*entity.Ticket, error)
		GetByBookingID(db *gorm.DB, id uint) ([]*entity.Ticket, error)
		CountByScheduleClassAndStatuses(db *gorm.DB, scheduleID uint, classID uint) (int64, error)
		CreateBulk(db *gorm.DB, tickets []*entity.Ticket) error
		UpdateBulk(db *gorm.DB, tickets []*entity.Ticket) error
		FindManyByIDs(db *gorm.DB, ids []uint) ([]*entity.Ticket, error)
		FindManyBySessionID(db *gorm.DB, sessionID uint) ([]*entity.Ticket, error)
		CancelManyBySessionID(db *gorm.DB, sessionID uint) error
		Paid(db *gorm.DB, id uint) error
		CheckIn(db *gorm.DB, id uint) error
	}

	SessionRepository interface {
		Create(db *gorm.DB, entity *entity.ClaimSession) error
		Update(db *gorm.DB, entity *entity.ClaimSession) error
		Delete(db *gorm.DB, entity *entity.ClaimSession) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.ClaimSession, error)
		GetByID(db *gorm.DB, id uint) (*entity.ClaimSession, error)
		GetByUUID(db *gorm.DB, uuid string) (*entity.ClaimSession, error)
		GetByUUIDWithLock(db *gorm.DB, uuid string, forUpdate bool) (*entity.ClaimSession, error)
		FindExpired(db *gorm.DB, expiryTime time.Time, limit int) ([]*entity.ClaimSession, error)
	}

	ClaimSessionRepository interface {
		Create(db *gorm.DB, entity *entity.ClaimSession) error
		Update(db *gorm.DB, entity *entity.ClaimSession) error
		Delete(db *gorm.DB, entity *entity.ClaimSession) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.ClaimSession, error)
		GetByID(db *gorm.DB, id uint) (*entity.ClaimSession, error)
		GetByUUID(db *gorm.DB, uuid string) (*entity.ClaimSession, error)
		GetByUUIDWithLock(db *gorm.DB, uuid string, forUpdate bool) (*entity.ClaimSession, error)
		FindExpired(db *gorm.DB, expiryTime time.Time, limit int) ([]*entity.ClaimSession, error)
	}

	AllocationRepository interface {
		Create(db *gorm.DB, entity *entity.Allocation) error
		Update(db *gorm.DB, entity *entity.Allocation) error
		Delete(db *gorm.DB, entity *entity.Allocation) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Allocation, error)
		GetByID(db *gorm.DB, id uint) (*entity.Allocation, error)
		LockByScheduleAndClass(db *gorm.DB, scheduleID uint, classID uint) (*entity.Allocation, error)
		GetByScheduleAndClass(db *gorm.DB, scheduleID uint, classID uint) (*entity.Allocation, error)
		FindByScheduleID(db *gorm.DB, scheduleID uint) ([]*entity.Allocation, error)
		CreateBulk(db *gorm.DB, allocations []*entity.Allocation) error
	}

	ManifestRepository interface {
		Create(db *gorm.DB, entity *entity.Manifest) error
		Update(db *gorm.DB, entity *entity.Manifest) error
		Delete(db *gorm.DB, entity *entity.Manifest) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Manifest, error)
		GetByID(db *gorm.DB, id uint) (*entity.Manifest, error)
		GetByShipAndClass(db *gorm.DB, shipID uint, classID uint) (*entity.Manifest, error)
		FindByShipID(db *gorm.DB, shipID uint) ([]*entity.Manifest, error)
	}

	FareRepository interface {
		Create(db *gorm.DB, entity *entity.Fare) error
		Update(db *gorm.DB, entity *entity.Fare) error
		Delete(db *gorm.DB, entity *entity.Fare) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Fare, error)
		GetByID(db *gorm.DB, id uint) (*entity.Fare, error)
		GetByManifestAndRoute(db *gorm.DB, manifestID uint, routeID uint) (*entity.Fare, error)
	}

	BookingRepository interface {
		Create(db *gorm.DB, entity *entity.Booking) error
		Update(db *gorm.DB, entity *entity.Booking) error
		Delete(db *gorm.DB, entity *entity.Booking) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Booking, error)
		GetByID(db *gorm.DB, id uint) (*entity.Booking, error)
		GetByOrderID(db *gorm.DB, id string) (*entity.Booking, error)
		PaidConfirm(db *gorm.DB, id uint) error
		UpdateReferenceNumber(tx *gorm.DB, bookingID uint, reference *string) error
	}

	ClassRepository interface {
		Create(db *gorm.DB, entity *entity.Class) error
		Update(db *gorm.DB, entity *entity.Class) error
		Delete(db *gorm.DB, entity *entity.Class) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Class, error)
		GetByID(db *gorm.DB, id uint) (*entity.Class, error)
	}

	ShipRepository interface {
		Create(db *gorm.DB, entity *entity.Ship) error
		Update(db *gorm.DB, entity *entity.Ship) error
		Delete(db *gorm.DB, entity *entity.Ship) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Ship, error)
		GetByID(db *gorm.DB, id uint) (*entity.Ship, error)
	}

	RouteRepository interface {
		Create(db *gorm.DB, entity *entity.Route) error
		Update(db *gorm.DB, entity *entity.Route) error
		Delete(db *gorm.DB, entity *entity.Route) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Route, error)
		GetByID(db *gorm.DB, id uint) (*entity.Route, error)
	}

	HarborRepository interface {
		Create(db *gorm.DB, entity *entity.Harbor) error
		Update(db *gorm.DB, entity *entity.Harbor) error
		Delete(db *gorm.DB, entity *entity.Harbor) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Harbor, error)
		GetByID(db *gorm.DB, id uint) (*entity.Harbor, error)
	}

	UserRepository interface {
		Create(db *gorm.DB, entity *entity.User) error
		Update(db *gorm.DB, entity *entity.User) error
		Delete(db *gorm.DB, entity *entity.User) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.User, error)
		GetByID(db *gorm.DB, id uint) (*entity.User, error)
		GetByUsername(db *gorm.DB, username string) (*entity.User, error)
		GetByEmail(db *gorm.DB, email string) (*entity.User, error)
		UpdatePassword(db *gorm.DB, userID uint, password string) error
	}

	RoleRepository interface {
		Create(db *gorm.DB, entity *entity.Role) error
		Update(db *gorm.DB, entity *entity.Role) error
		Delete(db *gorm.DB, entity *entity.Role) error
		Count(db *gorm.DB) (int64, error)
		GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Role, error)
		GetByID(db *gorm.DB, id uint) (*entity.Role, error)
	}

	AuthRepository interface {
		Create(db *gorm.DB, refreshToken *entity.RefreshToken) error
		Count(db *gorm.DB) (int64, error)
		GetAllRefreshToken(db *gorm.DB) ([]*entity.RefreshToken, error)
		GetRefreshToken(db *gorm.DB, id string) (*entity.RefreshToken, error)
		RevokeRefreshTokenByID(db *gorm.DB, id uuid.UUID) error
		CreatePasswordReset(db *gorm.DB, pr *entity.PasswordReset) error
		GetByToken(db *gorm.DB, token string) (*entity.PasswordReset, error)
		MarkAsUsed(db *gorm.DB, token string) error
		DeleteExpired(db *gorm.DB) error
	}
)
