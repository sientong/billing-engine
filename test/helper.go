package test

import (
	"fmt"
	"testing"
)

func ExpectPanicMessage(t *testing.T, fn func(), expectedMessage string) {
	t.Helper()

	defer func() {

		if r := recover(); r != nil {
			var actual string

			switch v := r.(type) {
			case error:
				actual = v.Error()
			case string:
				actual = v
			default:
				actual = fmt.Sprintf("%v", v)
			}

			if actual != expectedMessage {
				t.Errorf("Expected panic message:\n  %q\nbut got:\n  %q", expectedMessage, actual)
			}
		} else {
			t.Errorf("Expected panic with message %q, but no panic occurred", expectedMessage)
		}
	}()

	fn()
}
