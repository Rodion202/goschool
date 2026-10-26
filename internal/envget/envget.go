package envget

import (
	"fmt"
	"goschool/internal/models"
	"os"

	"github.com/joho/godotenv"
)

func ConfigInit() models.Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	return models.Config{
		PGXpass:        os.Getenv("POSTGRES_PASSWORD"),
		PGXuser:        os.Getenv("POSTGRES_USER"),
		PGXdb:          os.Getenv("POSTGRES_DB"),
		JWTsecret:      os.Getenv("JWT_SECRET"),
		PGXHost:        os.Getenv("POSTGRES_HOST"),
		ADMIN_EMAIL:    os.Getenv("ADMIN_EMAIL"),
		ADMIN_NAME:     os.Getenv("ADMIN_NAME"),
		ADMIN_PASSWORD: os.Getenv("ADMIN_PASSWORD"),
	}
}
