package validator_test

import (
	"app/platform/web/validator"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests for validator.RequiredJSON
func TestRequiredJSON(t *testing.T) {
	// Success cases
	// - success to validate required fields
	t.Run("success to validate required fields", func(t *testing.T) {
		// arrange
		// ...
		
		// act
		inputReader := strings.NewReader(`{"name":"John","age":20}`)
		inputKeys := []string{"name", "age"}
		err := validator.RequiredJSON(inputReader, inputKeys...)

		// assert
		require.NoError(t, err)
	})

	// Failure cases
	// - fail to validate required fields - key not found
	t.Run("fail to validate required fields - key not found", func(t *testing.T) {
		// arrange
		// ...

		// act
		inputReader := strings.NewReader(`{"name":"John"}`)
		inputKeys := []string{"name", "age"}
		err := validator.RequiredJSON(inputReader, inputKeys...)

		// assert
		require.ErrorIs(t, err, validator.ErrValidatorRequired)
		require.EqualError(t, err, "validator: required field not found - key age")
	})

	// - fail to validate required fields - invalid json
	t.Run("fail to validate required fields - invalid json", func(t *testing.T) {
		// arrange
		// ...

		// act
		inputReader := strings.NewReader(`{"name":"John","age":20`)
		inputKeys := []string{"name", "age"}
		err := validator.RequiredJSON(inputReader, inputKeys...)

		// assert
		require.Error(t, err)
	})
}