package database

import(
	"context"
	"time"
	"net/url"
	"fmt"
	//"os"
	"github.com/jackc/pgx/v5/pgxpool"
	"goschool/models"
)

func ConnToPostgres(p *models.PostgresReqs) (*pgxpool.Pool, error){
	if p.Host==""{
		p.Host="localhost:5432"
	}

	if p.User== "" || p.Pass== "" || p.Name == "" {
		return nil, fmt.Errorf("POSTGRES_USER,POSTGRES_PASSWORD and POSTGRES_DB must be set")
	}

	u:=url.URL{
		Scheme: "postgres",
		User: url.UserPassword(p.User,p.Pass),
		Host: p.Host,
		Path: p.Name,
		RawQuery: "sslmode=disable",
	}


	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return pgxpool.New(ctx, u.String())
}

func Connect(p *models.PostgresReqs) (*pgxpool.Pool,error) {
	pool,err:=ConnToPostgres(p)
	if err!=nil{
		return nil,err
	} else {
		fmt.Println("Postgres is working!")
	}

	if err:=pool.Ping(context.Background());err!=nil{
		fmt.Println("Ping failed:",err)
	}
	
	_,err=pool.Exec(context.Background(),`CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		password TEXT NOT NULL
	)`)
	if err!=nil{
		return nil,err
	} else {
		fmt.Println("Table created!")
	}
	return pool, nil
}
