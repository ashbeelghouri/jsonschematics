package api

import "strings"

func GetPathRegex(path string) string {
	path = strings.ReplaceAll(path, "*", ".*")
	path = strings.ReplaceAll(path, ":", "[^/]+")
	return "^" + path + "$"
}
