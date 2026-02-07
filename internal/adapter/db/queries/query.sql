-- name: CreateCalculation :one
INSERT INTO calculations (
  pack_sizes, target_amount, result_json, total_items
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: ListCalculations :many
SELECT * FROM calculations
ORDER BY created_at DESC;
