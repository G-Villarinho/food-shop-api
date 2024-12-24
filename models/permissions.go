package models

type Permission string

const (
	CreateOrderPermission Permission = "create_order"
	ListOrdersPermission  Permission = "list_orders"
)

var rolePermissions = map[Role][]Permission{
	Manager:  {ListOrdersPermission},
	Customer: {CreateOrderPermission},
}

func CheckPermission(role Role, permission Permission) bool {
	permissions, exists := rolePermissions[role]
	if !exists {
		return false
	}
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}
