package ddl

import (
	"strings"
)

func trim(str string) string {
	return strings.Trim(str, " \n\r\t")
}
