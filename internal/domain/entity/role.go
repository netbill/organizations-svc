package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID              uuid.UUID `json:"id"`
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Head            bool      `json:"head"`
	Editable        bool      `json:"editable"`
	Rank            uint      `json:"rank"`
	Name            string    `json:"name"`
	//Permissions     []Permission `json:"permissions"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//func (r Role) findCode(code CodeRolePermission) bool {
//	for _, perm := range r.Permissions {
//		if perm.Code == code {
//			return true
//		}
//	}
//
//	return false
//}
//
//func (r Role) CanManageAgglomeration() bool {
//	return r.findCode(RolePermissionManageAgglomeration)
//}
//
//func (r Role) CanManageCities() bool {
//	return r.findCode(RolePermissionManageCities)
//}
//
//func (r Role) CanManageRoles() bool {
//	return r.findCode(RolePermissionManageRoles)
//}
//
//func (r Role) CanManageInvites() bool {
//	return r.findCode(RolePermissionManageInvites)
//}
//
//func (r Role) CanManageMembers() bool {
//	return r.findCode(RolePermissionManageMembers)
//}
//
//func (r Role) CanEdit() bool {
//	return r.Editable
//}
//
//func (r Role) CanAddForUser() bool {
//	return !r.Head
//}
//
//func (r Role) CanRemoveForUser() bool {
//	return !r.Head
//}
//
//func (r Role) CompareRight(rank uint) bool {
//	return r.Rank > rank
//}
