package api_test

import (
	"ignis/internal/adapter/api"
	"ignis/internal/service"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCalculatorE2E(t *testing.T) {
	// Setup service and handler
	calcService := service.NewPackageCalculatorService()
	h := api.NewCalculatorHandler(calcService, nil) // nil repo for now, or mock if needed

	// Prepare form data
	formData := url.Values{}
	formData.Set("packSizes", "23, 31, 53")
	formData.Set("amount", "500000")

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true") // Simulate HTMX request

	// Create recorder
	w := httptest.NewRecorder()

	// Call handler
	h.Calculate(w, req)

	// Check status code
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.Status)
	}

	// Read body
	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	// Verify expected results in HTML and check the order (descending)
	expectedLines := []string{
		"<td>53</td><td>9429</td>",
		"<td>31</td><td>7</td>",
		"<td>23</td><td>2</td>",
	}

	lastIdx := -1
	for _, expected := range expectedLines {
		idx := strings.Index(html, expected)
		if idx == -1 {
			t.Errorf("expected HTML to contain %q, but it didn't.\nBody: %s", expected, html)
		}
		if idx < lastIdx {
			t.Errorf("expected HTML to have %q after previous result, but it appeared before.\nBody: %s", expected, html)
		}
		lastIdx = idx
	}

	// Verify other content
	if !strings.Contains(html, "Results for 500000 items:") {
		t.Errorf("expected HTML to contain total header")
	}
	if !strings.Contains(html, "Total items: <strong>500000</strong>") {
		t.Errorf("expected HTML to contain total amount")
	}
}
