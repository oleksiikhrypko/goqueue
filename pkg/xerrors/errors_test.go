package xerrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Consume(t *testing.T) {
	t.Parallel()
	t.Run("error", func(t *testing.T) {
		someErr := errors.New("some error")
		baseErr := New("base message").Consume(someErr).WithAdditionalInfo("test", map[string]any{"field": "value", "field2": "base_err_value2"})
		testErr := New("new test")
		newErr := testErr.Consume(baseErr).WithAdditionalInfo("test", map[string]any{"field2": "value2"})
		assert.True(t, errors.Is(newErr, testErr))
		assert.True(t, errors.Is(newErr, someErr))
		assert.True(t, errors.Is(newErr, baseErr))
		assert.Equal(t, "test", newErr.Error())

		resErr := &Error{}
		assert.True(t, errors.As(newErr, &resErr))
		assert.Equal(t, "test", resErr.Error())
		assert.Equal(t, "test", resErr.Message())

		ext := resErr.Extensions()

		assert.Equal(t, map[string]any{"field": "value", "field2": "value2"}, ext)

		testErr = New("new test")
		ext = testErr.Extensions()
		assert.Equal(t, map[string]any{}, ext)
	})
}
