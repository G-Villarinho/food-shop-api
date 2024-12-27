package utils

import "strconv"

func GetQueryStringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func GetQueryIntPointer(value string) *int {
	if value == "" {
		return nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}

	return &intValue
}
