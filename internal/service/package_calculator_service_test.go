package service

import (
	"ignis/internal/domain"
	"strings" // Added strings for cleaner error checking
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
			name: "minimal quantity test - 10 with sizes [6, 5, 2]",
			request: domain.CalculateRequest{
				PackSizes: []int{6, 5, 2},
				Amount:    10,
			},
			wantErr: false,
			validate: func(t *testing.T, result *domain.CalculateResult) {
				// IMPORTANT: A greedy algorithm would pick [6, 2, 2] (3 packs)
				// The DP algorithm MUST pick [5, 5] (2 packs)
				count := 0
				for _, v := range result.Packages {
					count += v
				}
				if count != 2 {
					t.Errorf("expected minimal pack count of 2 (5,5), got %d", count)
				}
				if result.Total != 10 {
					t.Errorf("expected total 10, got %d", result.Total)
				}
			},
		},
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
				if result.Total != 500000 {
					t.Errorf("expected exact total 500000, got %d", result.Total)
				}
				// Verify your specific requirement: {23: 2, 31: 7, 53: 9429}
				if result.Packages[53] != 9429 || result.Packages[31] != 7 || result.Packages[23] != 2 {
					t.Errorf("unexpected distribution: %+v", result.Packages)
				}
			},
		},
		{
			name: "no exact match possible",
			request: domain.CalculateRequest{
				PackSizes: []int{5, 10},
				Amount:    7,
			},
			wantErr:     true,
			errContains: "no exact combination possible",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Calculate(tt.request)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
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
