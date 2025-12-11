package db

import (
	"context"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool() (*pgxpool.Pool, error){
	dsn:= os.Getenv("DATABASE_URL")
	cfg, err :=pgxpool.ParseConfig(dsn)
	if err!=nil{
		return nil, err
	}

	pool, err:= pgxpool.NewWithConfig(context.Background(),cfg)
	if err!=nil{
		return nil,err
	}

	return pool, nil
}

