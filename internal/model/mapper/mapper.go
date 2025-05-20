package mapper

import (
	"eticket-api/internal/domain/entity"
	authentity "eticket-api/internal/domain/entity/auth"
	"eticket-api/internal/model"
	authmodel "eticket-api/internal/model/auth"

	"github.com/jinzhu/copier"
)

// Generic mapper between two types (Entity <-> DTO)
type Mapper[S any, D any] struct{}

// Convert from source to destination (e.g., Entity -> DTO)
func (m *Mapper[S, D]) ToModel(input *S) *D {
	dst := new(D)
	copier.Copy(dst, input)
	return dst
}

// Convert a slice from source to destination
func (m *Mapper[S, D]) ToModels(inputs []*S) []*D {
	result := make([]*D, 0, len(inputs))
	for _, item := range inputs {
		result = append(result, m.ToModel(item))
	}
	return result
}

// Convert from destination to source (e.g., DTO -> Entity)
func (m *Mapper[S, D]) ToEntity(input *D) *S {
	dst := new(S)
	copier.Copy(dst, input)
	return dst
}

// Convert a slice from destination to source
func (m *Mapper[S, D]) ToEntities(inputs []*D) []*S {
	result := make([]*S, 0, len(inputs))
	for _, item := range inputs {
		result = append(result, m.ToEntity(item))
	}
	return result
}

// TriMapper handles read, write, and update mappings
type TriMapper[E any, R any, W any, U any] struct {
	read   *Mapper[E, R]
	write  *Mapper[E, W]
	update *Mapper[E, U]
}

func NewMapper[E any, R any, W any, U any]() *TriMapper[E, R, W, U] {
	return &TriMapper[E, R, W, U]{
		read:   &Mapper[E, R]{},
		write:  &Mapper[E, W]{},
		update: &Mapper[E, U]{},
	}
}

func (m *TriMapper[E, R, W, U]) ToModel(e *E) *R {
	return m.read.ToModel(e)
}

func (m *TriMapper[E, R, W, U]) ToModels(es []*E) []*R {
	return m.read.ToModels(es)
}

func (m *TriMapper[E, R, W, U]) FromWrite(w *W) *E {
	return m.write.ToEntity(w)
}

func (m *TriMapper[E, R, W, U]) FromWrites(ws []*W) []*E {
	return m.write.ToEntities(ws)
}

func (m *TriMapper[E, R, W, U]) FromUpdate(u *U) *E {
	return m.update.ToEntity(u)
}

func (m *TriMapper[E, R, W, U]) FromUpdates(us []*U) []*E {
	return m.update.ToEntities(us)
}

// Global mappers: One per entity (read, write, update)
var (
	TicketMapper                    = NewMapper[entity.Ticket, model.ReadTicketResponse, model.WriteTicketRequest, model.UpdateTicketRequest]()
	RouteMapper                     = NewMapper[entity.Route, model.ReadRouteResponse, model.WriteRouteRequest, model.UpdateRouteRequest]()
	FareMapper                      = NewMapper[entity.Fare, model.ReadFareResponse, model.WriteFareRequest, model.UpdateFareRequest]()
	ManifestMapper                  = NewMapper[entity.Manifest, model.ReadManifestResponse, model.WriteManifestRequest, model.UpdateManifestRequest]()
	ShipMapper                      = NewMapper[entity.Ship, model.ReadShipResponse, model.WriteShipRequest, model.UpdateShipRequest]()
	SessionMapper                   = NewMapper[entity.ClaimSession, model.ReadClaimSessionResponse, model.WriteClaimSessionRequest, model.UpdateClaimSessionRequest]()
	TicketClassToSessionClassMapper = NewMapper[entity.Class, model.ClaimSessionTicketClassItem, model.WriteClaimSessionRequest, model.UpdateClaimSessionRequest]()
	ScheduleMapper                  = NewMapper[entity.Schedule, model.ReadScheduleResponse, model.WriteScheduleRequest, model.UpdateScheduleRequest]()
	ScheduleSessionMapper           = NewMapper[entity.Schedule, model.ClaimSessionSchedule, model.WriteScheduleRequest, model.UpdateScheduleRequest]()
	ClassMapper                     = NewMapper[entity.Class, model.ReadClassResponse, model.WriteClassRequest, model.UpdateClassRequest]()
	BookingMapper                   = NewMapper[entity.Booking, model.ReadBookingResponse, model.WriteBookingRequest, model.UpdateBookingRequest]()
	HarborMapper                    = NewMapper[entity.Harbor, model.ReadHarborResponse, model.WriteHarborRequest, model.UpdateHarborRequest]()
	AllocationMapper                = NewMapper[entity.Allocation, model.ReadAllocationResponse, model.WriteAllocationRequest, model.UpdateAllocationRequest]()
	UserMapper                      = NewMapper[authentity.User, authmodel.ReadUserResponse, authmodel.WriteUserRequest, authmodel.UpdateUserRequest]()
	RoleMapper                      = NewMapper[authentity.Role, authmodel.ReadRoleResponse, authmodel.WriteRoleRequest, authmodel.UpdateRoleRequest]()
	UserRoleMapper                  = NewMapper[authentity.UserRole, authmodel.ReadUserRoleResponse, authmodel.WriteUserRoleRequest, authmodel.UpdateUserRoleRequest]()
)
