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
	"regexp"
	"time"
)

// TODO: Store blocklist in ComputeHash Map and use prev string
var currentHead string
var blocklist = make(map[string]blockchain.Block)
// I need those sorted by fee to always incorporate max fees into mined blocklist
var unclaimedTransactions = treeset.NewWith(compareTxByCollectableFee)
var utxoList = make(map[string]blockchain.Txoutput)
var LINE_FEED = []byte{0x0A}
var REGEX_VALID_HASH = regexp.MustCompile(`[a-fA-F0-9]{32}`)

var peerList = make([]networking.Peer, 0, 5)

func main() {
	var head blockchain.Block = blockchain.CreateGenesisBlock()
	currentHead = head.ComputeHash()
	blocklist[currentHead] = head

	key1 := blockchain.CreateKeypair()
	key2 := blockchain.CreateKeypair()
	outputlist := []blockchain.Txoutput{blockchain.CreateTxOutput(0, key1.PublicKey), blockchain.CreateTxOutput(1, key2.PublicKey)}
	inputlist := []blockchain.Txinput{blockchain.CreateTxInput(&outputlist[0], key1), blockchain.CreateTxInput(&outputlist[1], key2)}
	transactions := []blockchain.Transaction{blockchain.Transaction{Outputs: outputlist, Inputs: inputlist}}
	newblock := blockchain.Mine(transactions, currentHead)
	blocklist[newblock.Hash] = newblock
	currentHead = newblock.Hash

	router := mux.NewRouter()
	router.HandleFunc("/pending_transaction", PostTransaction).Methods("POST")
	router.HandleFunc("/pending_transaction", GetTransactions).Methods("GET")
	router.HandleFunc("/peers", GetPeers).Methods("GET")
	router.HandleFunc("/ping", GetPing).Methods("GET")

	blockrouter := router.PathPrefix("/block").Subrouter().StrictSlash(true)
	blockrouter.HandleFunc("/", GetAllBlocks).Methods("GET")
	blockrouter.HandleFunc("/", PostBlock).Methods("POST")
	blockrouter.HandleFunc("/{hash:[a-fA-F0-9]+}", GetSpecificBlock).Methods("GET")

	httpsrv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
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
	if newtx != nil && newtx.Validate() {
		unclaimedTransactions.Add(newtx)
	} else if newtx == nil {
		writer.WriteHeader(400)
		writer.Write([]byte(fmt.Sprintf("JSON is invalid\n")))
	} else {
		writer.WriteHeader(400)
		writer.Write([]byte(fmt.Sprintf("Transaction %s is invalid\n", newtx.ComputeHash())))
	}
}

func GetTransactions(writer http.ResponseWriter, request *http.Request) {
	writeJson(unclaimedTransactions.Values(), writer)
}

func PostBlock(writer http.ResponseWriter, request *http.Request) {
	var newblock *blockchain.Block
	json.NewDecoder(request.Body).Decode(newblock)
	if newblock != nil {
		newblock.Hash = newblock.ComputeHash()
		if !newblock.Validate() {
			writer.WriteHeader(400)
			writer.Write([]byte(fmt.Sprintf("Block %s is invalid\n", newblock.Hash)))
			return
		}

		blocklist[newblock.Hash] = *newblock
		if (blockchain.ComputeBlockHeight(*newblock, &blocklist) > blockchain.ComputeBlockHeight(blocklist[currentHead], &blocklist)) {
			currentHead = newblock.ComputeHash()
		}

		if newblock != nil {
			unclaimedTransactions.Remove(newblock.Transactions)
		}
	} else {
		writer.WriteHeader(400)
		writer.Write([]byte(fmt.Sprintf("JSON is invalid\n")))
	}
}

func GetAllBlocks(writer http.ResponseWriter, request *http.Request) {
	var result = make([]blockchain.Block, 0, len(blocklist))
	for _, elem := range blocklist {
		result = append(result, elem)
	}
	writeJson(result, writer)
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
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(bytearr)
	writer.Write(LINE_FEED)
}
