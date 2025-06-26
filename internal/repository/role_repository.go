package repository

import (
	"eticket-api/internal/domain"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleRepository struct{}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
}

func (rr *RoleRepository) Count(db *gorm.DB) (int64, error) {
	var total int64
	result := db.Model(&domain.Role{}).Count(&total)
	return total, result.Error
}

func (rr *RoleRepository) Insert(db *gorm.DB, role *domain.Role) error {
	result := db.Create(role)
	return result.Error
}

func (rr *RoleRepository) InsertBulk(db *gorm.DB, roles []*domain.Role) error {
	result := db.Create(&roles)
	return result.Error
}

func (ar *RoleRepository) Update(db *gorm.DB, role *domain.Role) error {
	result := db.Save(role)
	return result.Error
}

func (rr *RoleRepository) UpdateBulk(db *gorm.DB, roles []*domain.Role) error {
	result := db.Save(&roles)
	return result.Error
}

func (ar *RoleRepository) Delete(db *gorm.DB, role *domain.Role) error {
	result := db.Select(clause.Associations).Delete(role)
	return result.Error
}

func (rr *RoleRepository) FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.Role, error) {
	roles := []*domain.Role{}
	query := db
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("role_name ILIKE ?", search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&roles).Error
	return roles, err
}
func (rr *RoleRepository) FindByID(db *gorm.DB, id uint) (*domain.Role, error) {
	role := new(domain.Role)
	result := db.First(&role, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return role, nil
}
