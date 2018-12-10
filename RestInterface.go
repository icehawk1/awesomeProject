package main

import (
	"awesomeProject/blockchain"
	"awesomeProject/networking"
	"awesomeProject/wallet"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// The head of the currently longest chain
var currentHead string
// The first block in the chain
var genesis string

// All known valid blocks: Blockhash -> Block
var blocklist = make(map[string]blockchain.Block)
// I need those sorted by fee to always incorporate max fees into mined blocklist
var unclaimedTransactions = treeset.NewWith(compareTxByCollectableFee)
// Hash -> UTXO
var utxoList = make(map[string]blockchain.Txoutput)
var LINE_FEED = []byte{0x0A}

func main() {
	host := *flag.String("host", "localhost", "Host to listen on")
	port := *flag.Int("port", 8000, "Port to listen on")
	flag.Parse()
	networking.SelfAddr = networking.CreatePeer(fmt.Sprintf("http://%s:%d",host,port))

	var head = blockchain.CreateGenesisBlock()
	genesis = head.ComputeHash()
	currentHead = genesis
	blocklist[currentHead] = head

	go exchangePeersContinously(1000)
	go mineContinously(200)
	go createTxContinously(1000)

	router := defineRoutingRules()

	httpsrv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("%s:%d", host, port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 25 * time.Second,
		ReadTimeout:  25 * time.Second,
	}

	log.Println("Listening for connections on", httpsrv.Addr)
	httpsrv.ListenAndServe()
}

func defineRoutingRules() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/pending_transaction", PostTransaction).Methods("POST")
	router.HandleFunc("/pending_transaction", GetTransactions).Methods("GET")
	router.HandleFunc("/peers", GetPeers).Methods("GET")
	router.HandleFunc("/ping", GetPing).Methods("GET")

	blockrouter := router.PathPrefix("/block").Subrouter().StrictSlash(true)
	blockrouter.HandleFunc("/", GetAllBlocks).Methods("GET")
	blockrouter.HandleFunc("/", PostBlock).Methods("POST")
	blockrouter.HandleFunc("/genesis", GetGenesisBlock).Methods("GET")
	blockrouter.HandleFunc("/head", GetHead).Methods("GET")
	blockrouter.HandleFunc("/{hash:[a-fA-F0-9]+}", GetSpecificBlocks).Methods("GET")

	return router
}

func exchangePeersContinously(delay int) {
	for {
		for i:=0; i<len(networking.PeerList); i++ {
			networking.ContactPeer(i)
		}
		time.Sleep(time.Duration(delay)*time.Millisecond)
	}
}

/* Simulates people using the chain */
func createTxContinously(maxdelay int) {
	keypair := blockchain.CreateKeypair()

	for {
		newtx := blockchain.CreateRandomTransaction(utxoList, keypair)
		if newtx != nil {
			networking.BroadcastTransaction(*newtx)
			unclaimedTransactions.Add(newtx)
		}
		time.Sleep(time.Duration(rand.Intn(maxdelay)) * time.Millisecond)
	}
}

/* Simulates continous mining activity. The random delay is to make mining harder. */
func mineContinously(maxdelay int) {
	keypair := blockchain.CreateKeypair()

	for {
		txToInclude := blockchain.SelectTransactionsForNextBlock(unclaimedTransactions)
		txToInclude, _ = blockchain.ClaimFees(txToInclude, keypair)
		txToInclude = append(txToInclude, blockchain.CreateCoinbaseTransaction(keypair.PublicKey))

		valid := false
		var newblock blockchain.Block
		for !valid {
			newblock, valid = blockchain.MineAttempt(txToInclude, currentHead)
			time.Sleep(time.Duration(rand.Intn(maxdelay)) * time.Millisecond)
		}

		if blockchain.ComputeBlockHeight(newblock, &blocklist) > blockchain.ComputeBlockHeight(blocklist[currentHead], &blocklist) {
			networking.BroadcastBlock(newblock)
			for _, tx := range txToInclude {
				unclaimedTransactions.Remove(tx)
			}
			blocklist[newblock.Hash] = newblock
			currentHead = newblock.Hash
		}
	}
}

func GetHead(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	writer.Write([]byte(currentHead))
}
func GetGenesisBlock(writer http.ResponseWriter, request *http.Request) {
	writeJson(blocklist[genesis], writer)
}

func PostTransaction(writer http.ResponseWriter, request *http.Request) {
	var newtx *blockchain.Transaction
	err := json.NewDecoder(request.Body).Decode(&newtx)
	if newtx != nil && newtx.Validate() {
		unclaimedTransactions.Add(newtx)
	} else if newtx == nil {
		writer.WriteHeader(400)
		writer.Write([]byte(fmt.Sprintf("JSON is invalid: %s\n", err)))
	} else {
		writer.WriteHeader(400)
		writer.Write([]byte(fmt.Sprintf("Transaction %s is invalid\n", newtx.ComputeHash())))
	}
}

func GetTransactions(writer http.ResponseWriter, request *http.Request) {

	vals := unclaimedTransactions.Values()
	txlist := make([]blockchain.Transaction, 0, len(vals))
	for _, elem := range vals {
		tx, _ := elem.(blockchain.Transaction)
		txlist = append(txlist, tx)
	}

	writeJson(vals, writer)
}

func PostBlock(writer http.ResponseWriter, request *http.Request) {
	var newblock *blockchain.Block
	json.NewDecoder(request.Body).Decode(&newblock)
	if newblock != nil {
		newblock.Hash = newblock.ComputeHash()
		if !newblock.Validate() {
			writer.WriteHeader(400)
			writer.Write([]byte(fmt.Sprintf("Block %s is invalid\n", newblock.Hash)))
			return
		}

		blocklist[newblock.Hash] = *newblock
		if blockchain.ComputeBlockHeight(*newblock, &blocklist) > blockchain.ComputeBlockHeight(blocklist[currentHead], &blocklist) {
			currentHead = newblock.ComputeHash()
		}

		if newblock != nil {
			unclaimedTransactions.Remove(newblock.Transactions)
			wallet.TxHasBeenPublished(newblock.Transactions.GetElements())
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
	peeraddr := request.FormValue("url")
	if peeraddr != "" {
		networking.AddPeer(networking.CreatePeer(peeraddr))
	}
	writeJson(networking.PeerList, writer)
}
func GetSpecificBlocks(writer http.ResponseWriter, request *http.Request) {
	// Return more than one block iff requested
	tmp := request.FormValue("num")
	num, err := strconv.Atoi(tmp)
	if err != nil || int(num) <= 0 {
		writer.WriteHeader(400)
		writer.Write([]byte(fmt.Sprintf("Incorrect number of blocks requested: %d , %s", num, tmp)))
		return
	}

	result := make([]blockchain.Block, 0, num)
	current, ok := mux.Vars(request)["hash"]
	for i := 0; i < num && ok && current != ""; i++ {
		block, ok := blocklist[current]
		if ok {
			result = append(result, block)
			current = block.Prev
		}
	}

	writeJson(result, writer)
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
