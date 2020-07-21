package main

import (
	"fmt"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/handler"
	"semprojdb/logger"

	"github.com/jackc/pgx/v4"

	"github.com/gorilla/mux"
)

func productsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	response := fmt.Sprintf("Product %s", id)
	fmt.Fprint(w, response)
}

func main() {
	lg := logger.NewLogger("", "", "log.log", pgx.LogLevelDebug)
	db := psql.NewDB(&psql.Config{
		Host:     "localhost",
		Database: "study_ev_table",
		User:     "postgres",
		Password: "docker",
	}, lg.SubLog("PSQL"))
	err := db.Connect()
	if err != nil {
		lg.Fatal(err)
	}
	h := handler.NewHandler(db, lg.SubLog("HANDLER"))
	fmt.Println(db, lg, lg.SubLog("kek"))
	router := mux.NewRouter()

	router.HandleFunc("/faculty", h.HandleFacultyCreate).Methods(http.MethodPost)
	router.HandleFunc("/faculty", h.HandleFacultyRead).Methods(http.MethodGet)
	router.HandleFunc("/faculty", h.HandleFacultyUpdate).Methods(http.MethodPut)
	router.HandleFunc("/faculty", h.HandleFacultyDelete).Methods(http.MethodDelete)

	router.HandleFunc("/products/{id:[0-9]+}", productsHandler).Methods(http.MethodGet)
	http.Handle("/", router)

	lg.Infof("Server is listening on 8181...")
	http.ListenAndServe(":8181", nil)
}
