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
	bankstatementhandler "github.com/faisalhardin/amartha-reconciliation-service/internal/http/bankstatement"
	reconciliationhandler "github.com/faisalhardin/amartha-reconciliation-service/internal/http/reconciliation"
	transactionhandler "github.com/faisalhardin/amartha-reconciliation-service/internal/http/transaction"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/db/xorm"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/messaging"
	bankstatementrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/repo/bankstatement"
	bankstatementfilerepo "github.com/faisalhardin/amartha-reconciliation-service/internal/repo/bankstatementfile"
	reconciliationrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/repo/reconciliation"
	transactionrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/repo/transaction"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/server"
	bankstatementuc "github.com/faisalhardin/amartha-reconciliation-service/internal/usecase/bankstatement"
	bankstatementfileuc "github.com/faisalhardin/amartha-reconciliation-service/internal/usecase/bankstatementfile"
	bankstatementprocessuc "github.com/faisalhardin/amartha-reconciliation-service/internal/usecase/bankstatementprocess"
	reconciliationuc "github.com/faisalhardin/amartha-reconciliation-service/internal/usecase/reconciliation"
	transactionuc "github.com/faisalhardin/amartha-reconciliation-service/internal/usecase/transaction"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/worker"
)

func gracefulShutdown(
	apiServer *http.Server,
	workerCancel context.CancelFunc,
	processQueue *messaging.BankStatementProcessQueue,
	reconciliationQueue *messaging.ReconciliationQueue,
	done chan bool,
) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	workerCancel()
	processQueue.Close()
	reconciliationQueue.Close()

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

	processQueue := messaging.NewBankStatementProcessQueue(100)
	reconciliationQueue := messaging.NewReconciliationQueue(100)
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	transactionDB := transactionrepo.NewTransactionDB(conn)
	bankStatementFileDB := bankstatementfilerepo.NewBankStatementFileDB(conn)
	bankStatementDB := bankstatementrepo.NewBankStatementDB(conn)
	reconciliationDB := reconciliationrepo.NewReconciliationDB(conn)

	reconciliationUC := reconciliationuc.NewReconciliationUC(
		bankStatementFileDB,
		bankStatementDB,
		transactionDB,
		reconciliationDB,
	)

	processUC := bankstatementprocessuc.NewBankStatementProcessUC(
		bankStatementFileDB,
		bankStatementDB,
		reconciliationQueue,
	)
	worker.StartBankStatementWorker(workerCtx, processQueue, processUC)
	worker.StartReconciliationWorker(workerCtx, reconciliationQueue, reconciliationUC)

	transactionUC := transactionuc.NewTransactionUC(transactionDB)
	bankStatementFileUC := bankstatementfileuc.NewBankStatementFileUC(bankStatementFileDB, processQueue)
	bankStatementUC := bankstatementuc.NewBankStatementUC(bankStatementDB, bankStatementFileDB)

	handlers := &entityhttp.Handlers{
		TransactionHandler:    transactionhandler.New(transactionUC),
		BankStatementHandler:  bankstatementhandler.New(bankStatementFileUC, bankStatementUC),
		ReconciliationHandler: reconciliationhandler.New(reconciliationUC),
	}

	apiServer := server.NewServer(handlers)

	done := make(chan bool, 1)
	go gracefulShutdown(apiServer, workerCancel, processQueue, reconciliationQueue, done)

	err = apiServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
