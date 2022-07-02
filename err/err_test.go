package err

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestErr_Error(t *testing.T) {
	tableTest := []struct {
		Name        string
		Err         *err
		ExpectedMsg string
	}{
		{
			Name:        "should return empty string error message",
			Err:         nil,
			ExpectedMsg: "",
		},
		{
			Name:        "should return error message",
			Err:         &err{Err: fmt.Errorf("something went wrong")},
			ExpectedMsg: fmt.Errorf("something went wrong").Error(),
		},
	}

	for _, tt := range tableTest {
		t.Run(tt.Name, func(t *testing.T) {
			msg := tt.Err.Error()
			assert.Equal(t, tt.ExpectedMsg, msg)
		})
	}
}

func TestErr_GRPCStatusCode(t *testing.T) {
	tableTest := []struct {
		Name           string
		Err            *err
		ExpectedStatus codes.Code
	}{
		{
			Name:           "should return gRPC status code 0",
			Err:            nil,
			ExpectedStatus: 0,
		},
		{
			Name:           fmt.Sprintf("should return %v", codes.Internal),
			Err:            &err{GRPCCode: codes.Internal},
			ExpectedStatus: codes.Internal,
		},
		{
			Name:           fmt.Sprintf("should return %v", codes.OK),
			Err:            &err{GRPCCode: codes.OK},
			ExpectedStatus: codes.OK,
		},
	}
	for _, tt := range tableTest {
		t.Run(tt.Name, func(t *testing.T) {
			code := tt.Err.GRPCStatusCode()
			assert.Equal(t, tt.ExpectedStatus, code)
		})
	}
}
