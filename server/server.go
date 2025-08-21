package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/matrix-go/bitcoin/core"
	"github.com/matrix-go/bitcoin/utils"
	"github.com/matrix-go/bitcoin/wallet"
	"net/http"
)

var cache = make(map[string]*core.Blockchain)

type Server struct {
	srv  *http.Server
	eg   *gin.Engine
	port int
}

func NewServer(port int) *Server {
	eg := gin.Default()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: eg,
	}
	return &Server{
		srv:  srv,
		eg:   eg,
		port: port,
	}
}

func (s *Server) GetBlockchain() *core.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		miner := wallet.NewWallet()
		bc = core.NewBlockchain(miner.Address(), s.port)
		fmt.Printf("miner private key: %s\n", miner.PrivateKeyStr())
		fmt.Printf("miner public key: %s\n", miner.PublicKeyStr())
		fmt.Printf("miner address: %s\n", miner.Address())
		cache["blockchain"] = bc
	}
	return bc
}

func (s *Server) Start() error {
	s.eg.GET("/blockchain", s.handleGetBlockchain)
	s.eg.POST("/transaction", s.handlePostAddTransaction)
	s.eg.GET("/transactions", s.handleGetTransactions)
	// TODO: just for test
	s.eg.GET("/mine", s.handleMine)
	s.eg.GET("/mining", s.handleStartMining)
	fmt.Println("server listening on", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) handleGetBlockchain(ctx *gin.Context) {
	chain := s.GetBlockchain()
	ctx.JSON(http.StatusOK, gin.H{"chain": chain})
}

type TransactionRequest struct {
	SenderAddress    string `json:"sender_address" binding:"required"`
	RecipientAddress string `json:"recipient_address" binding:"required"`
	Value            int64  `json:"value" binding:"required"`

	SenderPublicKey string `json:"sender_public_key" binding:"required"`
	Signature       string `json:"signature" binding:"required"`
}

func (s *Server) handlePostAddTransaction(ctx *gin.Context) {
	var req TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	publicKey, err := utils.PublicKeyFromString(req.SenderPublicKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	signature, err := utils.SignatureFromString(req.Signature)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bc := s.GetBlockchain()
	isCreated := bc.AddTransaction(req.SenderAddress, req.RecipientAddress, req.Value, publicKey, signature)
	if !isCreated {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "transaction verification failed"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"success": true})
}

func (s *Server) handleGetTransactions(ctx *gin.Context) {
	bc := s.GetBlockchain()
	transactions := bc.GetTransactionPools()
	ctx.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"length":       len(transactions),
	})
}

func (s *Server) handleMine(ctx *gin.Context) {
	bc := s.GetBlockchain()
	isMined := bc.Mining()
	ctx.JSON(http.StatusOK, gin.H{"success": isMined})
}

func (s *Server) handleStartMining(ctx *gin.Context) {
	bc := s.GetBlockchain()
	go bc.StartMining()
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
