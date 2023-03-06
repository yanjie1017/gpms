package main

import (
	"fmt"
	"net/http"

	// service "src/services"
	db "src/databases"
	handler "src/handlers"

	"github.com/gorilla/mux"
)

const portNumber = ":8000"

func main() {
	fmt.Println("Hello World")
	// service.GeneratePassword("a", "b", "c", "d", 1)
	/* Connect to database */
	dbConnection := db.ConnectDB()

	/* Defer closing database connection */
	defer dbConnection.DB.Close()

	// entryTag := model.PasswordEntryTag{
	// 	ClientId: 10000,
	// 	EntryId:  10000,
	// }
	// // result, err := dbConnection.ListTables()
	// result, err := dbConnection.RetrievePasswordInfo(entryTag)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(result.Length)

	router := mux.NewRouter()
	router.HandleFunc("/v1/health", handler.HandleHealthCheck).Methods(http.MethodPost)
	router.HandleFunc("/v1/password/generate", handler.HandlePasswordGeneration(dbConnection)).Methods((http.MethodPost))
	router.HandleFunc("/v1/entry/create", handler.HandlePasswordEntryCreation(dbConnection)).Methods(http.MethodPost)

	err := http.ListenAndServe(portNumber, router)
	if err != nil {
		fmt.Println(err)
	}

}
