package api_test

import (
	"context"
	"ignis/internal/adapter/api"
	dbsqlc "ignis/internal/adapter/db/sqlc"
	"ignis/internal/domain"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// MockRepository implements db.Repository
type MockRepository struct {
	Calculations []dbsqlc.Calculation
	CreateErr    error
	ListErr      error
	LastCreated  dbsqlc.CreateCalculationParams
}

func (m *MockRepository) CreateCalculation(ctx context.Context, arg dbsqlc.CreateCalculationParams) (dbsqlc.Calculation, error) {
	m.LastCreated = arg
	if m.CreateErr != nil {
		return dbsqlc.Calculation{}, m.CreateErr
	}
	calc := dbsqlc.Calculation{
		ID:           1,
		PackSizes:    arg.PackSizes,
		TargetAmount: arg.TargetAmount,
		ResultJson:   arg.ResultJson,
		TotalItems:   arg.TotalItems,
		CreatedAt:    pgtype.Timestamp{Time: time.Now(), Valid: true},
	}
	m.Calculations = append(m.Calculations, calc)
	return calc, nil
}

func (m *MockRepository) ListCalculations(ctx context.Context) ([]dbsqlc.Calculation, error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	return m.Calculations, nil
}

func (m *MockRepository) Close() {}

// MockCalculator implements domain.PackageCalculator
type MockCalculator struct {
	Result *domain.CalculateResult
	Err    error
}

func (m *MockCalculator) Calculate(req domain.CalculateRequest) (*domain.CalculateResult, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Result, nil
}

func TestCalculatorHandler_Calculate_Persistence(t *testing.T) {
	mockCalc := &MockCalculator{
		Result: &domain.CalculateResult{
			Packages: map[int]int{53: 1},
			Total:    53,
		},
	}
	mockRepo := &MockRepository{}
	h := api.NewCalculatorHandler(mockCalc, mockRepo)

	formData := url.Values{}
	formData.Set("packSizes", "23, 31, 53")
	formData.Set("amount", "53")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.Calculate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", w.Code)
	}

	// Verify persistence call
	if mockRepo.LastCreated.TargetAmount != 53 {
		t.Errorf("expected repo to be called with target amount 53, got %v", mockRepo.LastCreated.TargetAmount)
	}
	if mockRepo.LastCreated.PackSizes != "23, 31, 53" {
		t.Errorf("expected repo to be called with pack sizes '23, 31, 53', got %v", mockRepo.LastCreated.PackSizes)
	}
}

func TestCalculatorHandler_History(t *testing.T) {
	mockRepo := &MockRepository{
		Calculations: []dbsqlc.Calculation{
			{
				ID:           1,
				PackSizes:    "23, 31, 53",
				TargetAmount: 500000,
				TotalItems:   500000,
				CreatedAt:    pgtype.Timestamp{Time: time.Now(), Valid: true},
			},
		},
	}
	h := api.NewCalculatorHandler(nil, mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/history", nil)
	w := httptest.NewRecorder()

	h.History(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "23, 31, 53") {
		t.Errorf("expected history to contain pack sizes, but it didn't")
	}
	if !strings.Contains(body, "500000") {
		t.Errorf("expected history to contain amount, but it didn't")
	}
}
