package utils

import (
	"github.com/go-playground/validator/v10"
)

func GetCustomErrorMessage(field, tag string) string {
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
		return field + " phải lớn hơn hoặc bằng 1000"
	case "lt":
		return field + " phải nhỏ hơn " + extractValue(field, tag) + " ký tự"
	case "lte":
		return field + " phải nhỏ hơn hoặc bằng 2100"
	case "oneof":
		return field + "chỉ được 1 hoặc 2 hoặc 3"
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
	if (field == "Name" && tag == "min") || (field == "MajorName" && tag == "min") {
		return "2"
	} else if (field == "Name" && tag == "max") || (field == "MajorName" && tag == "max") {
		return "100"
	}
	if field == "MajorId" && tag == "min" {
		return "4"
	} else if field == "MajorId" && tag == "max" {
		return "20"
	}
	if field == "Phone" && tag == "len" {
		return "10"
	}
	if field == "OTPCode" && tag == "len" {
		return "6"
	}
	if field == "SubjectCode" && tag == "min" {
		return "4"
	} else if field == "SubjectCode" && tag == "max" {
		return "20"
	}
	if field == "SubjectName" && tag == "min" {
		return "2"
	} else if field == "SubjectName" && tag == "max" {
		return "100"
	}
	if field == "Department" && tag == "min" {
		return "2"
	} else if field == "Department" && tag == "max" {
		return "100"
	}
	if field == "Credits" && tag == "gt" {
		return "0"
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
