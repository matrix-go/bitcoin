package wallet_server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/matrix-go/bitcoin/server"
	"github.com/matrix-go/bitcoin/utils"
	"github.com/matrix-go/bitcoin/wallet"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"strings"
)

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
	eg.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	})
	return &Server{
		srv:  srv,
		eg:   eg,
		port: port,
	}
}

func (s *Server) Start() error {
	_, file, _, _ := runtime.Caller(0)
	s.eg.LoadHTMLGlob(path.Join(path.Dir(file), "templates/*.html"))
	s.eg.GET("/:pathname", func(c *gin.Context) {
		pathname := c.Param("pathname")
		if !strings.HasSuffix(pathname, ".html") {
			pathname = pathname + ".html"
		}
		c.HTML(http.StatusOK, pathname, nil)
	})
	s.eg.POST("/wallet", s.handlePostCreateWallet)
	s.eg.POST("/transaction", s.handlePostCreateTransaction)
	fmt.Println("server listening on", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) Gateway() string {
	return "http://127.0.0.1:5000"
}

func (s *Server) handlePostCreateWallet(ctx *gin.Context) {
	w := wallet.NewWallet()
	fmt.Printf("public key X: %x\n", w.PublicKey().X.Bytes())
	fmt.Printf("public key Y: %x\n", w.PublicKey().X.Bytes())
	fmt.Printf("private key D: %x\n", w.PrivateKey().D.Bytes())
	ctx.JSON(http.StatusOK, w)
}

type TransactionRequest struct {
	SenderPrivateKey string `json:"sender_private_key" binding:"required"`
	SenderPublicKey  string `json:"sender_public_key" binding:"required"`
	SenderAddress    string `json:"sender_address" binding:"required"`
	RecipientAddress string `json:"recipient_address" binding:"required"`
	Amount           string `json:"amount" binding:"required"`
}

func (s *Server) handlePostCreateTransaction(ctx *gin.Context) {
	req := TransactionRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	publicKey, err := utils.PublicKeyFromString(req.SenderPublicKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	privateKey, err := utils.PrivateKeyFromString(req.SenderPrivateKey, *publicKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value, err := strconv.ParseInt(req.Amount, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("public key X: %x\n", publicKey.X.Bytes())
	fmt.Printf("public key Y: %x\n", publicKey.X.Bytes())
	fmt.Printf("private key D: %x\n", privateKey.D.Bytes())
	tx := wallet.NewTransaction(privateKey, publicKey, req.SenderAddress, req.RecipientAddress, value)
	sig := tx.GenerateSignature()

	txReq := server.TransactionRequest{
		SenderAddress:    req.SenderAddress,
		RecipientAddress: req.RecipientAddress,
		Value:            value,
		SenderPublicKey:  req.SenderPublicKey,
		Signature:        sig.String(),
	}

	fmt.Printf("tx: %v\n", tx)
	m, err := json.Marshal(&txReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request, err := http.NewRequest(http.MethodPost, s.Gateway()+"/transaction", bytes.NewBuffer(m))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": resp.Status})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"success": true})
}
