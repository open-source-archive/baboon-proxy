package ltm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorMessage(t *testing.T) {
	err := ErrorLTM{Code: 404, Message: "Could not connect to LB"}
	expected := "Could not connect to LB"
	assert.Equal(t, expected, err.ErrorMessage(), "Should be equal")
}
