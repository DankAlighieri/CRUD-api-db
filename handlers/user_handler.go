package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dankalighieri/crud-api/models"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	DB *sql.DB
}

// Param db: ponteiro para a conexao com o banco de dados
// para garantir que todos os handlers interajam com o mesmo banco
//
// Construtor do UserHandler para inicializar a struct com a conexao com o banco de dados
func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (userHandler *UserHandler) ReadUsers(writer http.ResponseWriter, request *http.Request) {
	// Realiza a consulta nas linhas da tabela de usuarios do banco de dados
	rows, err := userHandler.DB.Query("SELECT * FROM users")
	if err != nil {
		http.Error(writer, "Erro ao consultar usuarios", http.StatusInternalServerError)
		return
	}

	// Realizando a decodificacao dos dados recuperados pela query
	users := make([]models.User, 0)

	for rows.Next() {
		var user models.User

		// mapear os dados da linha para a struct user
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			log.Println("nao foi possivel mapear o usuario", err)
		}

		users = append(users, user)
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(users)
}

func (userHandler *UserHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	user := models.User{}

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, "Erro ao decodificar o usuario", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO users (username, email) VALUES ($1, $2)"
	_, err = userHandler.DB.Exec(query, user.Username, user.Email)
	if err != nil {
		log.Fatal("nao foi possivel inserir o usuario", err)
	}

	log.Println("Usuario inserido com sucesso")
}

func (userHandler *UserHandler) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	user := models.User{}

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, "Erro ao decodificar o usuario", http.StatusBadRequest)
		return
	}

	newUsername := user.Username
	newEmail := user.Email

	id := mux.Vars(request)["id"]

	// verificar se usuario existe
	var exists bool

	err = userHandler.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		log.Fatal("problema com a query no banco de dados", err)
	}

	if !exists {
		log.Fatal("usuario nao encontrado")
	}

	query := "UPDATE users SET username = $1, email = $2 WHERE id = $3"
	_, erro := userHandler.DB.Exec(query, newUsername, newEmail, id)

	if erro != nil {
		http.Error(writer, "Erro ao atualizar usuario", http.StatusInternalServerError)
		log.Fatal("nao foi possivel atualizar o usuario", erro)
	}

	log.Println("Usuario atualizado com sucesso")
}

// TODO
func (userHandler *UserHandler) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"] // recuperando o ID passado na URL

	var exists bool

	err := userHandler.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		log.Println("erro ao verificar se o usuario existe", err)
		return
	}

	if !exists {
		log.Printf("Usuario com ID %s nao encontrado", id)
		return
	}

	query := "DELETE FROM users WHERE id = $1"
	_, error := userHandler.DB.Exec(query, id)
	if error != nil {
		log.Println("nao foi possivel deletar o usuario", err)
		return
	}

	log.Println("Usuario deletado com sucesso")
}
