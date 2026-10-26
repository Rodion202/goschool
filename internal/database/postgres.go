package database

import (
	"context"
	"fmt"
	"goschool/internal/models"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func (d Database) Ping() error {
	return d.Pool.Ping(context.Background())
}

func (d Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}

func (d Database) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	row := d.Pool.QueryRow(ctx, `SELECT id,name,email,password,role FROM users WHERE email = $1`, email)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d Database) AddUser(ctx context.Context, u *models.User) error {
	if u.Role == "" {
		u.Role = "student"
	}
	return d.Pool.QueryRow(
		ctx,
		`INSERT INTO users (name, email, password,role) VALUES ($1, $2, $3, $4) RETURNING id,role`,
		u.Name, u.Email, u.Password, u.Role,
	).Scan(&u.Id, &u.Role)
}

func (d Database) GetUserById(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	row := d.Pool.QueryRow(ctx, `SELECT id,email,name,password,role FROM users WHERE id = $1`, id)
	err := row.Scan(&u.Id, &u.Email, &u.Name, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d Database) GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := d.Pool.Query(ctx, `SELECT id, email, name,role FROM users`)
	if err != nil {
		return nil, fmt.Errorf("error %w", err)
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User

		err := rows.Scan(&u.Id, &u.Email, &u.Name, &u.Role)
		if err != nil {
			return nil, fmt.Errorf("error %w", err)
		}

		users = append(users, u)
	}
	return users, nil
}

func (d Database) DeleteUserById(ctx context.Context, id int64) error {
	_, err := d.Pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func ConnToPostgres(p *models.Config) (*pgxpool.Pool, error) {
	if p.PGXHost == "" {
		p.PGXHost = "localhost:5432"
	}

	if p.PGXuser == "" || p.PGXpass == "" || p.PGXdb == "" {
		return nil, fmt.Errorf("POSTGRES_USER,POSTGRES_PASSWORD and POSTGRES_DB must be set")
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(p.PGXuser, p.PGXpass),
		Host:     p.PGXHost,
		Path:     p.PGXdb,
		RawQuery: "sslmode=disable",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return pgxpool.New(ctx, u.String())
}

func Connect(drop bool, config *models.Config) (*Database, error) {
	pool, err := ConnToPostgres(config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		fmt.Println("Ping failed:", err)
		return nil, err
	}
	fmt.Println("Postgres is working!")

	if drop {
		_, err = pool.Exec(ctx, `DROP TABLE IF EXISTS users`)
		if err != nil {
			return nil, err
		}
		fmt.Println("Table users dropped!")
	}

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'student'
	)`)

	if err != nil {
		return nil, err
	}
	fmt.Println("Table created!")

	db := Database{Pool: pool}

	return &db, nil
}
