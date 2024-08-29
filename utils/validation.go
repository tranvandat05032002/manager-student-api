package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func GetCustomErrorMessage(field, tag string) string {
	fmt.Println("Field --> ", field)
	fmt.Println("Tag --> ", tag)
	switch tag {
	case "required":
		return field + " không được để trống"
	case "email":
		return field + " phải là một địa chỉ email hợp lệ"
	case "url":
		return field + "có định dạng là một URL"
	case "alpha":
		return field + " chỉ được chứa ký tự chữ cái"
	case "min":
		return field + " phải có ít nhất " + extractValue(field, tag) + " ký tự"
	case "max":
		return field + " không được quá " + extractValue(field, tag) + " ký tự"
	case "len":
		return field + " phải có đúng " + extractValue(field, tag) + " ký tự"
	case "gt":
		return field + " phải lớn hơn " + extractValue(field, tag) + " ký tự"
	case "gte":
		return field + " phải lớn hơn hoặc bằng " + extractValue(field, tag) + " ký tự"
	case "lt":
		return field + " phải nhỏ hơn " + extractValue(field, tag) + " ký tự"
	case "lte":
		return field + " phải nhỏ hơn hoặc bằng " + extractValue(field, tag) + " ký tự"
	default:
		return field + " không hợp lệ"
	}
}

func extractValue(field, tag string) string {
	if field == "Password" && tag == "min" {
		return "6"
	} else if field == "Password" && tag == "max" {
		return "30"
	}
	if field == "Name" && tag == "min" {
		return "2"
	} else if field == "Name" && tag == "max" {
		return "100"
	}
	if field == "Phone" && tag == "len" {
		return "10"
	}
	if field == "OTPCode" && tag == "len" {
		return "6"
	}
	return ""
}

func GetErrorMessagesResponse(err error) map[string]string {
	validatorErrors := err.(validator.ValidationErrors)
	errorMessages := make(map[string]string)
	for _, e := range validatorErrors {
		field := e.StructField()
		tag := e.Tag()
		errorMessages[field] = GetCustomErrorMessage(field, tag)
	}
	return errorMessages
}
