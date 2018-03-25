// +build gijit

package rand

import "testing"

func TestFloat32(t *testing.T) {
	t.Skip("slow")
}

func TestConcurrent(t *testing.T) {
	t.Skip("using nosync")
}
