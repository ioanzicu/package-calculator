package service

import (
	"errors"
	"ignis/internal/domain"
	"sort"
)

// PackageCalculatorService implements the domain.PackageCalculator interface
type PackageCalculatorService struct{}

// NewPackageCalculatorService creates a new instance of PackageCalculatorService
func NewPackageCalculatorService() *PackageCalculatorService {
	return &PackageCalculatorService{}
}

// Calculate implements the package optimization algorithm
// It finds the optimal distribution of packages to fulfill the requested amount exactly
func (s *PackageCalculatorService) Calculate(req domain.CalculateRequest) (*domain.CalculateResult, error) {
	if len(req.PackSizes) == 0 {
		return nil, errors.New("pack sizes cannot be empty")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Validate pack sizes
	for _, size := range req.PackSizes {
		if size <= 0 {
			return nil, errors.New("pack sizes must be greater than zero")
		}
	}

	// Sort pack sizes in descending order
	sizes := make([]int, len(req.PackSizes))
	copy(sizes, req.PackSizes)
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	// Find optimal combination that equals the target amount
	result := s.findExactCombination(sizes, req.Amount)

	return result, nil
}

// findExactCombination finds a package distribution that equals the target amount exactly
// Algorithm: Start with largest pack, then try to fill remainder with smaller packs
func (s *PackageCalculatorService) findExactCombination(sizes []int, target int) *domain.CalculateResult {
	result := &domain.CalculateResult{
		Packages: make(map[int]int),
		Total:    0,
	}

	if len(sizes) == 0 {
		return result
	}

	// Start with the largest pack size
	largestPack := sizes[0]

	// Try different counts of the largest pack, starting from the maximum possible
	maxLargestPacks := target / largestPack

	// Try from max down to 0 to find a combination that works
	for numLargest := maxLargestPacks; numLargest >= 0; numLargest-- {
		remainder := target - (numLargest * largestPack)

		if remainder == 0 {
			// Perfect fit with just the largest packs
			if numLargest > 0 {
				result.Packages[largestPack] = numLargest
				result.Total = numLargest * largestPack
			}
			return result
		}

		// Try to fill remainder with smaller packs
		if len(sizes) > 1 {
			smallerSizes := sizes[1:]
			smallerResult := s.findExactCombination(smallerSizes, remainder)

			// Check if we found an exact match
			if smallerResult.Total == remainder {
				// Found a valid combination
				if numLargest > 0 {
					result.Packages[largestPack] = numLargest
				}
				for size, count := range smallerResult.Packages {
					result.Packages[size] = count
				}
				result.Total = target
				return result
			}
		}
	}

	// If no exact combination found, return empty result
	// This shouldn't happen with valid pack sizes, but handle gracefully
	return result
}
