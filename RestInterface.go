package main

import (
	"awesomeProject/blockchain"
	"awesomeProject/networking"
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// TODO: Store blocklist in Hash Map and use prev string
var currentHead = blockchain.CreateGenesisBlock()
var blocklist = make(map[string]blockchain.Block)
// I need those sorted by fee to always incorporate max fees into mined blocklist
var unclaimedTransactions = treeset.NewWith(compareTxByCollectableFee)
var utxoList = make(map[string]blockchain.Txoutput)
var LINE_FEED = []byte{0x0A}

var peerList = make([]networking.Peer, 0, 5)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/transaction", PostTransaction).Methods("POST")
	router.HandleFunc("/peers", GetPeers).Methods("GET")
	router.HandleFunc("/ping", GetPing).Methods("GET")

	blockrouter := router.PathPrefix("/block").Subrouter().StrictSlash(true)
	blockrouter.HandleFunc("/", GetAllBlocks).Methods("GET")
	blockrouter.HandleFunc("/", PostBlock).Methods("POST")
	blockrouter.HandleFunc("/{id:[0-9]+}", GetSpecificBlock).Methods("GET")

	httpsrv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 25 * time.Second,
		ReadTimeout:  25 * time.Second,
	}

	httpsrv.ListenAndServe()
	log.Println("Listening for connections")
}

func PostTransaction(writer http.ResponseWriter, request *http.Request) {
	var newtx *blockchain.Transaction
	json.NewDecoder(request.Body).Decode(newtx)
	if newtx != nil {

	}
}
func PostBlock(writer http.ResponseWriter, request *http.Request) {
	var newblock *blockchain.Block
	json.NewDecoder(request.Body).Decode(newblock)
	log.Println(newblock)

	if newblock != nil {
		unclaimedTransactions.Remove(newblock.Transactions)
	}
}

func GetAllBlocks(writer http.ResponseWriter, request *http.Request) {
	writeJson(blocklist, writer)
}
func GetPeers(writer http.ResponseWriter, request *http.Request) {
	writeJson(peerList, writer)
}
func GetSpecificBlock(writer http.ResponseWriter, request *http.Request) {
	blockhash, _ := mux.Vars(request)["hash"]
	block, ok := blocklist[blockhash]
	if (ok) {
		writeJson(block, writer)
	} else {
		writer.WriteHeader(404)
		writer.Write([]byte(fmt.Sprintf("Block %s does not exist\n", blockhash)))
	}
}
func GetPing(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("pong\n"))
}

func compareTxByCollectableFee(a, b interface{}) int {
	txA, _ := a.(blockchain.Transaction)
	txB, _ := b.(blockchain.Transaction)
	feeA := txA.ComputePossibleFee()
	feeB := txB.ComputePossibleFee()

	switch {
	case feeA > feeB:
		return 1
	case feeA < feeB:
		return -1
	default:
		return 0
	}
}
func writeJson(obj interface{}, writer http.ResponseWriter) {
	bytearr, _ := json.Marshal(obj)
	writer.WriteHeader(200)
	writer.Write(bytearr)
	writer.Write(LINE_FEED)
}
func insertBlock(block blockchain.Block) {
	blocklist[block.ComputeHash()] = block
	if (blockchain.ComputeBlockHeight(block, &blocklist) > blockchain.ComputeBlockHeight(currentHead, &blocklist)) {
		currentHead = block
	}
}
