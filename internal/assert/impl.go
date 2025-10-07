package assert

import "testing"

func Equal[T comparable](t *testing.T, exp, act T, fargs ...any) {
	if act != exp {
		if len(fargs) >= 1 {
			format, ok := fargs[0].(string)
			if ok {
				t.Errorf(format, fargs[1:])
			} else {
				t.Error(fargs)
			}
		} else {
			t.Errorf("Expected: %+v\nGot     : %+v", exp, act)
		}
	}
}
