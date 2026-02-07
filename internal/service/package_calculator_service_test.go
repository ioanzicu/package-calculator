package service

import (
	"ignis/internal/domain"
	"testing"
)

func TestPackageCalculatorService_Calculate(t *testing.T) {
	service := NewPackageCalculatorService()

	tests := []struct {
		name        string
		request     domain.CalculateRequest
		wantErr     bool
		errContains string
		validate    func(t *testing.T, result *domain.CalculateResult)
	}{
		{
			name: "example case - 500000 with pack sizes 23, 31, 53",
			request: domain.CalculateRequest{
				PackSizes: []int{23, 31, 53},
				Amount:    500000,
			},
			wantErr: false,
			validate: func(t *testing.T, result *domain.CalculateResult) {
				if result == nil {
					t.Fatal("result should not be nil")
				}
				// Verify we have packages
				if len(result.Packages) == 0 {
					t.Error("expected packages in result")
				}
				// Verify total is at least the requested amount
				if result.Total < 500000 {
					t.Errorf("total %d is less than requested amount 500000", result.Total)
				}
				t.Logf("Result: %+v", result.Packages)
				t.Logf("Total: %d", result.Total)
			},
		},
		{
			name: "small amount - 100 with pack sizes 23, 31, 53",
			request: domain.CalculateRequest{
				PackSizes: []int{23, 31, 53},
				Amount:    100,
			},
			wantErr: false,
			validate: func(t *testing.T, result *domain.CalculateResult) {
				if result == nil {
					t.Fatal("result should not be nil")
				}
				if result.Total < 100 {
					t.Errorf("total %d is less than requested amount 100", result.Total)
				}
				t.Logf("Result: %+v", result.Packages)
			},
		},
		{
			name: "exact match - 53 with pack sizes 23, 31, 53",
			request: domain.CalculateRequest{
				PackSizes: []int{23, 31, 53},
				Amount:    53,
			},
			wantErr: false,
			validate: func(t *testing.T, result *domain.CalculateResult) {
				if result == nil {
					t.Fatal("result should not be nil")
				}
				if result.Packages[53] != 1 {
					t.Errorf("expected 1 pack of size 53, got %d", result.Packages[53])
				}
				if result.Total != 53 {
					t.Errorf("expected total 53, got %d", result.Total)
				}
			},
		},
		{
			name: "single pack size",
			request: domain.CalculateRequest{
				PackSizes: []int{10},
				Amount:    100,
			},
			wantErr: false,
			validate: func(t *testing.T, result *domain.CalculateResult) {
				if result == nil {
					t.Fatal("result should not be nil")
				}
				if result.Packages[10] != 10 {
					t.Errorf("expected 10 packs of size 10, got %d", result.Packages[10])
				}
				if result.Total != 100 {
					t.Errorf("expected total 100, got %d", result.Total)
				}
			},
		},
		{
			name: "error - empty pack sizes",
			request: domain.CalculateRequest{
				PackSizes: []int{},
				Amount:    100,
			},
			wantErr:     true,
			errContains: "pack sizes cannot be empty",
		},
		{
			name: "error - zero amount",
			request: domain.CalculateRequest{
				PackSizes: []int{10, 20},
				Amount:    0,
			},
			wantErr:     true,
			errContains: "amount must be greater than zero",
		},
		{
			name: "error - negative amount",
			request: domain.CalculateRequest{
				PackSizes: []int{10, 20},
				Amount:    -100,
			},
			wantErr:     true,
			errContains: "amount must be greater than zero",
		},
		{
			name: "error - invalid pack size",
			request: domain.CalculateRequest{
				PackSizes: []int{10, 0, 20},
				Amount:    100,
			},
			wantErr:     true,
			errContains: "pack sizes must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Calculate(tt.request)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
