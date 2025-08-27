package main

import (
	"log"
	"net/http"

	"github.com/dankalighieri/crud-api/config"
	"github.com/dankalighieri/crud-api/handlers"
	"github.com/dankalighieri/crud-api/models"
	"github.com/gorilla/mux"
)

// Responsavel por iniciar a conexao com o banco de dados e mapear as rotas
func main() {
	dbConnection, err := config.SetupDatabase()
	if err != nil {
		log.Fatal("problema na conexao com o banco de dados", err)
	}

	_, err = dbConnection.Exec(models.CreateTableSQL)
	if err != nil {
		log.Fatal("problema na query", err)
	}

	// Realiza o fechamento da conexao com o banco de dados ao finalizar a execucao do main
	defer dbConnection.Close()

	// Criando as rotas

	router := mux.NewRouter()

	UserHandler := handlers.NewUserHandler(dbConnection)
	// CreateUserHandler

	router.HandleFunc("/users", UserHandler.ReadUsers).Methods("GET")
	router.HandleFunc("/users", UserHandler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", UserHandler.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users/{id}", UserHandler.UpdateUser).Methods("POST")

	// Inicializando o servidor
	log.Fatal(http.ListenAndServe(":8080", router))
}
