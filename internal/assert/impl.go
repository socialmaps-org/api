package assert

import "testing"

func Equal[T comparable](t *testing.T, exp, act T, fmtArgs ...any) {
	if act != exp {
		if len(fmtArgs) >= 1 {
			fmt, ok := fmtArgs[0].(string)
			if ok {
				t.Errorf(fmt, fmtArgs[1:])
			} else {
				t.Error(fmtArgs...)
			}
		} else {
			t.Errorf("Expected: %+v\nGot     : %+v", exp, act)
		}
	}
}
