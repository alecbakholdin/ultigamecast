package assert

import (
	"fmt"
)

func That(condition bool, msgFmt string, args ...any) {
	if !condition {
		panic(fmt.Sprintf(msgFmt, args...))
	}
}
