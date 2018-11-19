package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

var chain = CreateChain()
var LINE_FEED = []byte{0x0A}

func main() {

	chain.Mine()
	chain.Mine()
	chain.Mine()

	router := mux.NewRouter()
	blockrouter := router.PathPrefix("/block").Subrouter().StrictSlash(true)
	blockrouter.HandleFunc("/",GetChain).Methods("GET")
	blockrouter.HandleFunc("/",PostBlock).Methods("POST")
	blockrouter.HandleFunc("/{id:[0-9]+}",GetSpecificBlock).Methods("GET")

	httpsrv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 25 * time.Second,
		ReadTimeout:  25 * time.Second,
	}

	httpsrv.ListenAndServe()
	log.Println("Listening for connections")
}

func PostBlock(writer http.ResponseWriter, request *http.Request) {
	var newblock Block
	json.NewDecoder(request.Body).Decode(&newblock)
	log.Println(newblock)
}

func GetChain(writer http.ResponseWriter, request *http.Request) {
	bytearr,_ := json.Marshal(chain)
	writer.WriteHeader(200)
	writer.Write(bytearr)
	writer.Write(LINE_FEED)
}

func GetSpecificBlock(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(request)["id"])
	if id < 0 || id >= len(chain.Blocklist) {
		writer.WriteHeader(404)
		writer.Write([]byte(fmt.Sprintf("Block %d does not exist\n",id)))
		return
	}

	writer.WriteHeader(200)
	writer.Write(toJson(chain.Blocklist[id]))
	writer.Write(LINE_FEED)
}

func toJson(block Block) []byte {
	bytearr,_ := json.Marshal(block)
	return bytearr
}
