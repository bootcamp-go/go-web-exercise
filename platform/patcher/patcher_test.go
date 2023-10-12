package patcher_test

import (
	"app/platform/patcher"
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests for patcher.Patch
func TestPatch(t *testing.T) {
	// Success cases
	// - success to patch a struct - patch not adjusting to tags
	// - success to patch a struct - patch adjusting to tags
	// - success to patch a struct - assiganble case is json: patch field is float64 and struct field is int
	// - success to patch a struct - unexported field not set
	// - success to patch a struct - key not found in patch
	t.Run("success to patch a struct - patch not adjusting to tags", func(t *testing.T) {
		// arrange
		// ...

		// act
		type Schema struct {
			Name string `patcher:"name"`
			Age  int    `patcher:"age"`
		}
		inputPtr   := &Schema{}
		inputPatch := map[string]any{
			"name": "John",
			"age":  20,
		}
		err := patcher.Patch(inputPtr, inputPatch)

		// assert
		require.NoError(t, err)
		require.Equal(t, "John", inputPtr.Name)
		require.Equal(t, 20, inputPtr.Age)
	})

	t.Run("success to patch a struct - patch adjusting to tags", func(t *testing.T) {
		// arrange
		// ...

		// act
		type Schema struct {
			Name string `patcher:"name"`
			Age  int    `patcher:"age"`
		}
		inputPtr   := &Schema{}
		inputPatch := map[string]any{
			"name": "John",
			"age":  20,
		}
		err := patcher.Patch(inputPtr, inputPatch)

		// assert
		require.NoError(t, err)
		require.Equal(t, "John", inputPtr.Name)
		require.Equal(t, 20, inputPtr.Age)
	})

	t.Run("success to patch a struct - patch not adjusting to tags", func(t *testing.T) {
		// arrange
		// ...

		// act
		type Schema struct {
			Name string `patcher:"name"`
			Age  int    `patcher:"age"`
		}
		inputPtr   := &Schema{}
		inputPatch := map[string]any{
			"name": "John",
			"age":  20.0,
		}
		err := patcher.Patch(inputPtr, inputPatch)

		// assert
		require.NoError(t, err)
		require.Equal(t, "John", inputPtr.Name)
		require.Equal(t, 20, inputPtr.Age)
	})

	t.Run("success to patch a struct - unexported field not set", func(t *testing.T) {
		// arrange
		// ...

		// act
		type Schema struct {
			Name string `patcher:"name"`
			age  int    `patcher:"age"`
		}
		inputPtr   := &Schema{}
		inputPatch := map[string]any{
			"name": "John",
			"age":  20,
		}
		err := patcher.Patch(inputPtr, inputPatch)

		// assert
		require.NoError(t, err)
		require.Equal(t, "John", inputPtr.Name)
		require.Equal(t, 0, inputPtr.age)	
	})

	t.Run("success to patch a struct - key not found in patch", func(t *testing.T) {
		// arrange
		// ...

		// act
		type Schema struct {
			Name string `patcher:"name"`
			Age  int    `patcher:"age"`
		}
		inputPtr   := &Schema{}
		inputPatch := map[string]any{
			"name": "John",
		}
		err := patcher.Patch(inputPtr, inputPatch)

		// assert
		require.NoError(t, err)
		require.Equal(t, "John", inputPtr.Name)
		require.Equal(t, 0, inputPtr.Age)
	})

	// Failure cases
	// - failure to patch a struct - invalid type
	t.Run("failure to patch a struct - invalid type", func(t *testing.T) {
		// arrange
		// ...

		// act
		type Schema struct {
			Name string `patcher:"name"`
			Age  int    `patcher:"age"`
		}
		inputPtr   := &Schema{}
		inputPatch := map[string]any{
			"name": "John",
			"age":  "20",
		}
		err := patcher.Patch(inputPtr, inputPatch)

		// assert
		require.ErrorIs(t, err, patcher.ErrPatcherInvalidType)
		require.EqualError(t, err, "patcher: invalid type - fieldName Age - fieldTag age")
	})
}