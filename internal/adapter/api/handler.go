package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"ignis/internal/adapter/db"
	dbsqlc "ignis/internal/adapter/db/sqlc"
	"ignis/internal/domain"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type PageData struct {
	Title   string
	Message string
}

type CalculatorHandler struct {
	calculator domain.PackageCalculator
	repo       db.Repository
}

func NewCalculatorHandler(calculator domain.PackageCalculator, repo db.Repository) *CalculatorHandler {
	return &CalculatorHandler{
		calculator: calculator,
		repo:       repo,
	}
}

func (h *CalculatorHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	packSizesStr := r.FormValue("packSizes")
	amountStr := r.FormValue("amount")

	// Parse pack sizes
	packSizesStrSlice := strings.Split(packSizesStr, ",")
	packSizes := make([]int, 0, len(packSizesStrSlice))
	for _, sizeStr := range packSizesStrSlice {
		sizeStr = strings.TrimSpace(sizeStr)
		if sizeStr == "" {
			continue
		}
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf("<div class='error'>Invalid pack size: %s</div>", sizeStr)))
			return
		}
		packSizes = append(packSizes, size)
	}

	// Parse amount
	amount, err := strconv.Atoi(strings.TrimSpace(amountStr))
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf("<div class='error'>Invalid amount: %s</div>", amountStr)))
		return
	}

	// Calculate
	result, err := h.calculator.Calculate(domain.CalculateRequest{
		PackSizes: packSizes,
		Amount:    amount,
	})
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf("<div class='error'>Calculation error: %s</div>", err.Error())))
		return
	}

	// Build HTML response
	var html strings.Builder
	html.WriteString("<div class='result-success'>")
	html.WriteString(fmt.Sprintf("<h3>Results for %d items:</h3>", amount))
	html.WriteString("<table class='result-table'>")
	html.WriteString("<tr><th>Pack Size</th><th>Quantity</th></tr>")

	// Sort pack sizes for consistent output
	sortedSizes := make([]int, 0, len(result.Packages))
	for size := range result.Packages {
		sortedSizes = append(sortedSizes, size)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedSizes)))

	for _, packSize := range sortedSizes {
		count := result.Packages[packSize]
		html.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%d</td></tr>", packSize, count))
	}

	html.WriteString("</table>")
	html.WriteString(fmt.Sprintf("<p class='total'>Total items: <strong>%d</strong></p>", result.Total))
	html.WriteString("</div>")

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Trigger", "calculation-done")
	w.Write([]byte(html.String()))

	if h.repo == nil {
		return
	}

	// Save to history asynchronously or synchronously? Let's do it synchronously for simplicity for now
	ctx := r.Context()
	resultJson, _ := json.Marshal(result.Packages)
	_, err = h.repo.CreateCalculation(ctx, dbsqlc.CreateCalculationParams{
		PackSizes:    packSizesStr,
		TargetAmount: int32(amount),
		ResultJson:   resultJson,
		TotalItems:   int32(result.Total),
	})
	if err != nil {
		fmt.Printf("failed to save calculation: %v\n", err)
	}
}

func (h *CalculatorHandler) History(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	calculations, err := h.repo.ListCalculations(ctx)
	if err != nil {
		http.Error(w, "Failed to load history", http.StatusInternalServerError)
		return
	}

	var html strings.Builder
	html.WriteString("<div class='history-container'>")
	html.WriteString("<h3>Recent Calculations</h3>")
	if len(calculations) == 0 {
		html.WriteString("<p>No history yet.</p>")
	} else {
		html.WriteString("<table class='history-table'>")
		html.WriteString("<tr><th>Date</th><th>Packs</th><th>Amount</th><th>Total</th></tr>")
		for _, calc := range calculations {
			html.WriteString("<tr>")
			html.WriteString(fmt.Sprintf("<td>%s</td>", calc.CreatedAt.Time.Format("2006-01-02 15:04")))
			html.WriteString(fmt.Sprintf("<td>%s</td>", calc.PackSizes))
			html.WriteString(fmt.Sprintf("<td>%d</td>", calc.TargetAmount))
			html.WriteString(fmt.Sprintf("<td>%d</td>", calc.TotalItems))
			html.WriteString("</tr>")
		}
		html.WriteString("</table>")
	}
	html.WriteString("</div>")

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html.String()))
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "index.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:   "Interactive Package Calculator",
		Message: "Enter pack sizes (comma-separated) and the amount to calculate the optimal package distribution.",
	}

	tmpl.Execute(w, data)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
