package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"goschool/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct{
	Pool *pgxpool.Pool
}

func (d Database)Ping()error{
	err:=d.Pool.Ping(context.Background())
	if err!=nil{
		return err
	}
	return nil
}

func (d Database)Close(){
	if d.Pool !=nil{
		d.Pool.Close()
	}
}

func (d Database)AddUser(u models.User)error {
	_,err:=d.Pool.Exec(context.Background(),`INSERT INTO users (email, name, password) VALUES ($1,$2,$3)`,u.Email,u.Name,u.Password)
	if err!=nil{
		return err
	}
	return nil
}

func (d Database)GetUserById(id string)(*models.User,error){
	var u models.User
	err:=d.Pool.QueryRow(context.Background(),`SELECT id,email,name,password FROM users WHERE id = $1`,id).Scan(&u.Id,&u.Email,&u.Name,&u.Password)
	if err!=nil{
		if errors.Is(err,pgx.ErrNoRows){
			return nil,fmt.Errorf("There is no user with id %s",id)
		}
	}
	return &u,nil
}

func (d Database)GetAllUsers()([]models.User, error){
	rows,err:=d.Pool.Query(context.Background(),`SELECT id, email, name, password FROM users`)
	if err!=nil{
		return nil,fmt.Errorf("Error %w",err) 
	}
	defer rows.Close()
	var users []models.User
	for rows.Next(){
		var u models.User

		err:=rows.Scan(&u.Id,&u.Email,&u.Name,&u.Password)
		if err!=nil{
			return nil, fmt.Errorf("Error %w",err)
		}

		users = append(users, u)
	}
	return users, nil
}

func (d Database)DeleteUserById(id string)error{
	_,err:=d.Pool.Exec(context.Background(),`DELETE FROM users WHERE id = $1`,id)
	if err!=nil{
		return err
	}
	return nil
}

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

func Connect(p *models.PostgresReqs, drop bool) (models.Storage,error) {
	pool,err:=ConnToPostgres(p)
	if err!=nil{
		return nil,err
	} else {
		fmt.Println("Postgres is working!")
	}

	if err:=pool.Ping(context.Background());err!=nil{
		fmt.Println("Ping failed:",err)
	}

	if drop {
		_,err=pool.Exec(context.Background(),`DROP TABLE IF EXISTS users`)
	if err!=nil{
		return nil,err
	} else {
		fmt.Println("Table users dropped!")
	}
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
	db:=Database{Pool:pool}
	return db, nil
}
