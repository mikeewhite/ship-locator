package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShip_TrimsWhitespaceFromName(t *testing.T) {
	ship := NewShip(1234, "\t   CALL SIGN   \n\n", 80, 100)
	assert.Equal(t, "CALL SIGN", ship.Name)
}
