package everdb

import "testing"
import "os"

func TestMain(m *testing.M) {
	code := m.Run()

	os.Remove("test.db")
	os.Exit(code)
}

func equals(a, b []byte) bool {
	l := len(a)
	if l != len(b) {
		return false
	}
	for i := 0; i < l; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFileBlockDevice(t *testing.T) {
	t.Parallel()

	h, err := NewFileBlockDevice("test.db", false, true)

	if nil != err {
		t.Error(err)
	}

	if val := h.Len(); val != 0 {
		t.Error("Initial block device was not empty", val)
	}

	err = h.Resize(1)

	if nil != err {
		t.Errorf("Resize failed: %v", err)
	}

	if h.Len() != 1 {
		t.Error("Resize failed")
	}

	zero := make([]byte, BLOCK_SIZE)
	for i := 0; i < BLOCK_SIZE; i++ {
		zero[i] = 0
	}

	b, err := h.Get(0)

	if nil != err {
		t.Error("Failed to get block")
	}

	if !equals(b, zero) {
		t.Error("fresh block is not zeroed")
	}

}
