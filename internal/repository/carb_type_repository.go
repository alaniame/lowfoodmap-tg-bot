package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
)

type CarbTypeRepository struct {
	pool *pgxpool.Pool
}

func (r CarbTypeRepository) GetCarbIds(s string) ([]int, error) {
	if s == "" {
		return nil, nil
	}
	var carbIds []int
	var err error
	for _, carbNames := range strings.Split(s, " ") {
		normalizedCarbStr := strings.ToLower(carbNames)
		query := `SELECT carb_id FROM carb_types WHERE carb_name = $1;`
		row := r.pool.QueryRow(context.Background(), query, normalizedCarbStr)
		var id int
		err = row.Scan(&id)
		if err != nil {
			return nil, err
		}
		carbIds = append(carbIds, id)
	}
	return carbIds, err
}
