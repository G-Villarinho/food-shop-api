package validation

var ValidationMessages = map[string]string{
	"required":        "This field is required",
	"email":           "Invalid email format",
	"min":             "Value is too short",
	"max":             "Value is too long",
	"eqfield":         "Fields do not match",
	"gt":              "The value must be greater than zero",
	"datetime":        "Invalid date format",
	StrongPasswordTag: "Password must be at least 8 characters long, contain an uppercase letter, a number, and a special character",
}
