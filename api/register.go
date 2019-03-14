package api

import (
	"strings"
)

func tableOf(values string) string {
	return strings.Replace(values, ":", "", 99999)
}
