package acceptance

import "strings"

func toUnix(t string) string {
	return strings.ReplaceAll(t, "\r", "")
}
