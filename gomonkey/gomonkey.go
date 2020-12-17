package funcs

import (
	"strings"
)

func MyFunc(p string) []string {
	return strings.Split(p, ",")
}
