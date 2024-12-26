package utils

func GetQueryStringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
