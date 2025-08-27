package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// SetupDatabase
// Configura a conexao com o banco de dados Postgres
// (Criando um singleton utilizando um ponteiro)
func SetupDatabase() (*sql.DB, error) {
	godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("problema na conexao com o banco de dados", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("nao foi possivel realizar o ping com o db", err)
		return nil, err
	}

	log.Println("Conexao com o banco de dados estabelecida com sucesso")

	return db, nil
}
