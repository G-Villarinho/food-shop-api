package models

type Permission string

const (
	CreateOrderPermission            Permission = "create_order"
	CancelOrderPermission            Permission = "cancel_order"
	ApproveOrderPermission           Permission = "approve_order"
	DispatchOrderPermission          Permission = "dispatch_order"
	DeliverOrderPermission           Permission = "deliver_order"
	ListOrdersPermission             Permission = "list_orders"
	CreateEvaluationPermission       Permission = "create_evaluation"
	ListEvaluationsPermission        Permission = "list_evaluations"
	UpdateEvaluationAnswerPermission Permission = "update_evaluation_answer"
	GetEvaluationSummaryPermission   Permission = "get_evaluation_summary"
	UpdateMenuPermission             Permission = "update_menu"
)

var rolePermissions = map[Role][]Permission{
	Manager: {ListOrdersPermission, CancelOrderPermission, ApproveOrderPermission, DispatchOrderPermission, ListEvaluationsPermission,
		UpdateEvaluationAnswerPermission, GetEvaluationSummaryPermission, UpdateMenuPermission},
	Customer: {CreateOrderPermission, DeliverOrderPermission, CreateEvaluationPermission},
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
