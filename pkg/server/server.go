package server

import (
	"log"
	"fmt"
	"time"
	"strconv"
	"net/http"

	"github.com/gorilla/mux"
)

func Serve() {
	LoadData()

	r := mux.NewRouter()

  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	r.HandleFunc("/start", startHandler)
	r.HandleFunc("/walk/{choice}/{nodeid}", walkHandler)

	srv := &http.Server{
    Handler:      r,
    Addr:         "127.0.0.1:3000",
    WriteTimeout: 15 * time.Second,
    ReadTimeout:  15 * time.Second,
  }

	fmt.Println("Listening at port 3000...")
  log.Fatal(srv.ListenAndServe())
}

func startHandler(res http.ResponseWriter, req *http.Request) {
	state := GetStartState();
	err := WriteHtml(res, state)
	if err != nil {
		fmt.Fprintf(res, "ERROR: %v\n", err)
	}
}

func walkHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	nid, err := strconv.Atoi(vars["nodeid"])
	if err != nil {
		fmt.Fprintf(res, "ERROR: %v\n", err)
	}
	state := GetWalkState(vars["choice"], uint32(nid));
	err = WriteHtml(res, state)
	if err != nil {
		fmt.Fprintf(res, "ERROR: %v\n", err)
	}
}





