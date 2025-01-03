package internal

type ContextKey string

const (
	UserIDKey       ContextKey = "user_id"
	SessionIDKey    ContextKey = "session_id"
	RestaurantIDKey ContextKey = "restaurant_id"
	RoleKey         ContextKey = "role"
)
