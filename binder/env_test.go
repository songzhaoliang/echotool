package binder

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvBinder(t *testing.T) {
	os.Setenv("ID", "1")
	os.Setenv("NAME", "peter")

	u := &User{}
	err := EnvBinder.Bind(nil, u)

	assert.NoError(t, err)
	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "peter", u.Name)
}
