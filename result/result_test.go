package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abhissng/internal/result" // Replace with the actual path to your package
)

func TestNewSuccess(t *testing.T) {
	value := "success value"
	successResult := result.NewSuccess(&value)

	assert.True(t, successResult.IsSuccess())
	assert.False(t, successResult.IsError())

	val, err := successResult.Value()
	assert.NoError(t, err)
	assert.Equal(t, value, *val)
}

func TestNewError(t *testing.T) {
	testErr := errors.New("test error")
	errorResult := result.NewError[any](testErr)

	assert.False(t, errorResult.IsSuccess())
	assert.True(t, errorResult.IsError())

	_, err := errorResult.Value()
	assert.Error(t, err)
	assert.Equal(t, testErr, err)

	assert.Equal(t, testErr, errorResult.Error())
}

func TestToResult(t *testing.T) {
	value := "success value"
	successResult := result.ToResult(&value, nil)

	assert.IsType(t, &result.Success[string]{}, successResult)

	errorResult := result.ToResult[string](nil, errors.New("test error"))
	assert.IsType(t, &result.Failure[string]{}, errorResult)
}

func TestCastError(t *testing.T) {
	value := "success value"
	successResult := result.NewSuccess(&value)

	// Cast to same type should work
	castResult := result.CastError[string, string](successResult)
	assert.IsType(t, &result.Failure[string]{}, castResult)

	// Cast to different type should fail
	castErrorResult := result.CastError[string, int](successResult)
	assert.IsType(t, &result.Failure[int]{}, castErrorResult)
	assert.EqualError(t, castErrorResult.Error(), "cannot cast a success result")

	// Cast error result should return a new error result
	testErr := errors.New("test error")
	errorResult := result.NewError[string](testErr)
	castErrorResult = result.CastError[string, int](errorResult)
	assert.IsType(t, &result.Failure[int]{}, castErrorResult)

	// Cast error result with specific type should return a new error result
	specificErr := fmt.Errorf("specific error: %w", testErr)
	errorResult = result.NewError[string](specificErr)
	castErrorResultint := result.CastError[string, error](errorResult)
	assert.IsType(t, &result.Failure[error]{}, castErrorResultint)
	assert.EqualError(t, castErrorResultint.Error(), "specific error: test error")
}

func TestMapError(t *testing.T) {
	value := "success value"
	successResult := result.NewSuccess(&value)

	// Mapping a success result should fail
	mappedResult := result.MapError[string, error](successResult, func(err error) error {
		return fmt.Errorf("mapped error")
	})
	assert.IsType(t, &result.Failure[error]{}, mappedResult)
	assert.EqualError(t, mappedResult.Error(), "cannot map a success result")

	testErr := errors.New("test error")
	errorResult := result.NewError[string](testErr)

	// Mapping an error result should return a new error result
	mappedResult = result.MapError[string, error](errorResult, func(err error) error {
		return fmt.Errorf("mapped error: %w", err)
	})
	assert.IsType(t, &result.Failure[error]{}, mappedResult)
	assert.EqualError(t, mappedResult.Error(), "mapped error: test error")
}
