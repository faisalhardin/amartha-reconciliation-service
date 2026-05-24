package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/config"
	entityhttp "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/http"
	transactionhandler "github.com/faisalhardin/amartha-reconciliation-service/internal/http/transaction"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/db/xorm"
	bankstatementrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/repo/bankstatement"
	transactionrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/repo/transaction"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/server"
	transactionuc "github.com/faisalhardin/amartha-reconciliation-service/internal/usecase/transaction"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

func main() {
	cfg, err := config.New("amartha-reconciliation-service")
	if err != nil {
		panic(fmt.Sprintf("config error: %s", err))
	}

	conn, err := xorm.NewDBConnection(cfg)
	if err != nil {
		panic(fmt.Sprintf("database connection error: %s", err))
	}
	defer func() {
		if err := conn.CloseDBConnection(); err != nil {
			log.Printf("close database connection error: %v", err)
		}
	}()

	transactionDB := transactionrepo.NewTransactionDB(conn)
	_ = bankstatementrepo.NewBankStatementDB(conn)

	transactionUC := transactionuc.NewTransactionUC(transactionDB)
	handlers := &entityhttp.Handlers{
		TransactionHandler: transactionhandler.New(transactionUC),
	}

	apiServer := server.NewServer(handlers)

	done := make(chan bool, 1)
	go gracefulShutdown(apiServer, done)

	err = apiServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
