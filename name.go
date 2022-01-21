package log

import (
	"strings"
)

func splitName(name string) (first, rest string) {
	i := strings.IndexByte(name, '.')
	if i == -1 {
		first = name
		return
	}
	return name[:i], name[i+1:]
}
