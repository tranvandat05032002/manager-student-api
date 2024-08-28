package utils

import "strings"

func IsAllowedImageExt(ext string) bool {
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, e := range allowedExts {
		if strings.EqualFold(ext, e) {
			return true
		}
	}
	return false
}
