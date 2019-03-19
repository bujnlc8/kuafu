package kuafu

import "fmt"


// format string
func FormatString(s string, args ...interface{}) string {
	return fmt.Sprintf(s, args...)
}