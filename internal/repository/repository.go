package repository

import (
	"github.com/jackc/pgx/v4"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) Repository {
	return Repository{conn: conn}
}
