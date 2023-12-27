package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"strings"
)

type CarbTypePostgres struct {
	conn *pgx.Conn
}

func (r CarbTypePostgres) GetCarbIds(s string) ([]int, error) {
	if s == "" {
		return nil, nil
	}
	var carbIds []int
	var err error
	for _, carbNames := range strings.Split(s, " ") {
		normalizedCarbStr := strings.ToLower(carbNames)
		query := `SELECT carb_id FROM carb_types WHERE carb_name = $1;`
		row := r.conn.QueryRow(context.Background(), query, normalizedCarbStr)
		var id int
		err = row.Scan(&id)
		if err != nil {
			return nil, err
		}
		carbIds = append(carbIds, id)
	}
	return carbIds, err
}
