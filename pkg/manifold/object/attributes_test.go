package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttributeSet(t *testing.T) {
	sys := New("")
	sys.Attrs().Set("test", "value")

	assert.True(t, sys.Attrs().Has("test"))
	assert.Equal(t, "value", sys.Attrs().Get("test").(string))

	sys.Attrs().Set("test", 42)

	assert.True(t, sys.Attrs().Has("test"))
	assert.Equal(t, 42, sys.Attrs().Get("test").(int))

	sys.Attrs().Del("test")

	assert.False(t, sys.Attrs().Has("test"))
	assert.Nil(t, sys.Attrs().Get("test"))
}
