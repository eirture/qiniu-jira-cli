package cmdutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPureServiceName(t *testing.T) {
	testcases := []struct {
		name string
		want string
	}{
		{"qboxrspub", "qboxrspub"},
		{"qboxrspub(已发布)", "qboxrspub"},
		{"qboxrsf （已发布）", "qboxrsf"},
		{"qboxrsf test", "qboxrsf"},
		{"sisyphus_sche_v2", "sisyphus_sche_v2"},
	}

	for _, tc := range testcases {
		assert.Equal(t, tc.want, GetPureServiceName(tc.name))
	}
}
