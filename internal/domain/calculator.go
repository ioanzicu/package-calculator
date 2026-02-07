package domain

// CalculateRequest represents the input for package calculation
type CalculateRequest struct {
	PackSizes []int
	Amount    int
}

// CalculateResult represents the output of package calculation
type CalculateResult struct {
	Packages map[int]int // map[packSize]count
	Total    int         // total items in all packages
}

// PackageCalculator defines the interface for package calculation service
type PackageCalculator interface {
	Calculate(req CalculateRequest) (*CalculateResult, error)
}
