package http

import (
	"testing"
)

// go test -v -count=1 ./pkg/http/ -run ^TestParams$
func TestParams(t *testing.T) {
	path := "/users/:username/:ages"
	prefix, maskedPath := maskPath(path)

	// t.Logf("prefix = [%s], maskedPath = [%s]\n", prefix, maskedPath)

	path = "/users/alex/32"
	params := makeParams(prefix, maskedPath, path)

	for key, value := range params {
		t.Logf("params[%s] = '%s'\n", key, value)
	}
}
