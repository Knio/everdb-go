package everdb

import "testing"

func TestHelloer(t *testing.T) {
	t.Parallel()

	h := Database{}
	if val := h.HelloWorld(); val != 5 {
		t.Error("Value was wrong:", val)
	}
}
