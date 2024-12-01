package db

import (
	"context"
	"database/sql"
	"dynamic-proxy/config"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// var RedisClient *redis.Client
var PostgresDB *sql.DB
var ctx = context.Background()

func InitPostgres(pgUrlData config.PostgresAddr) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		pgUrlData.Host, pgUrlData.Port, pgUrlData.User, pgUrlData.Password, pgUrlData.Dbname)
	var err error
	PostgresDB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	if err = PostgresDB.Ping(); err != nil {
		return err
	}

	log.Println("Connected to PostgreSQL")
	return nil
}

type Tenant struct {
	ID string
}

func GetTenant(hostname string) (*Tenant, error) {
	query := `SELECT id FROM tenants WHERE hostname = $1`
	// Remove t- from beginning of hostname (key)
	hostName := hostname[2:]

	row := PostgresDB.QueryRow(query, hostName)

	var tenant Tenant
	err := row.Scan(&tenant.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no tenant found with hostname %s", hostname)
		}
		return nil, err
	}

	return &tenant, nil
}
