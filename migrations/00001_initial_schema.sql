-- +goose Up
CREATE TABLE calculations (
  id SERIAL PRIMARY KEY,
  pack_sizes text NOT NULL,
  target_amount integer NOT NULL,
  result_json jsonb NOT NULL,
  total_items integer NOT NULL,
  created_at timestamp NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE calculations;
