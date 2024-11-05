package middleware

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetUserStatus(t *testing.T) {
	tests := []struct {
		name string
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetUserStatus(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetUserStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthRequired(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AuthRequired(tt.args.c)
		})
	}
}

func TestIsAuthenticated(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IsAuthenticated(tt.args.c)
		})
	}
}
