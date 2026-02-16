package service

import (
	"errors"
	"ignis/internal/domain"
	"math"
	"sort"
)

// PackageCalculatorService implements the domain.PackageCalculator interface
type PackageCalculatorService struct{}

// NewPackageCalculatorService creates a new instance of PackageCalculatorService
func NewPackageCalculatorService() *PackageCalculatorService {
	return &PackageCalculatorService{}
}

func (s *PackageCalculatorService) Calculate(req domain.CalculateRequest) (*domain.CalculateResult, error) {
	if len(req.PackSizes) == 0 {
		return nil, errors.New("pack sizes cannot be empty")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// 1. Prepare and sort sizes (ascending helps DP efficiency)
	sizes := make([]int, len(req.PackSizes))
	copy(sizes, req.PackSizes)
	sort.Ints(sizes)

	// 2. Setup DP arrays
	// dp[i] = min packs needed for amount i
	// parent[i] = the size of the pack used to get to amount i (for reconstruction)
	dp := make([]int, req.Amount+1)
	parent := make([]int, req.Amount+1)

	// Initialize DP with "Infinity"
	for i := 1; i <= req.Amount; i++ {
		dp[i] = math.MaxInt32
	}
	dp[0] = 0

	// 3. Fill DP table: O(Amount * PackSizes)
	//
	for _, size := range sizes {
		for i := size; i <= req.Amount; i++ {
			if dp[i-size] != math.MaxInt32 {
				// If using this pack results in FEWER total packs than what we had...
				if dp[i-size]+1 < dp[i] {
					dp[i] = dp[i-size] + 1
					parent[i] = size
				}
			}
		}
	}

	// 4. Check if a solution exists
	if dp[req.Amount] == math.MaxInt32 {
		return nil, errors.New("no exact combination possible for the requested amount")
	}

	// 5. Reconstruct the counts by walking backwards through 'parent'
	//
	resMap := make(map[int]int)
	curr := req.Amount
	for curr > 0 {
		size := parent[curr]
		resMap[size]++
		curr -= size
	}

	return &domain.CalculateResult{
		Packages: resMap,
		Total:    req.Amount,
	}, nil
}
