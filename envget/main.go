package envget

import(
	"fmt"
	"os"
	"goschool/models"
	"github.com/joho/godotenv"
)


func GetPostgresReqs() *models.PostgresReqs{
	if err:=godotenv.Load();err!=nil{
		fmt.Println("No .env file found")
	}
	return &models.PostgresReqs{
		User: os.Getenv("POSTGRES_USER"),
		Pass: os.Getenv("POSTGRES_PASSWORD"),
		Name: os.Getenv("POSTGRES_DB"),
		Host: os.Getenv("POSTGRES_HOST"),
	}
}
