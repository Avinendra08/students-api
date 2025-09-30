package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avinendra08/students-api/internal/config"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("welcome to golang backend"))
	})

	//setup server
	server := http.Server {
		Addr: cfg.HTTPServer.Addr,
		Handler: router,
	}
	
	slog.Info("server started ", slog.String("address",cfg.HTTPServer.Addr))
	//fmt.Printf("server started %s",cfg.Addr)

	//making gracefull shutdown : whenever we do ctrl+c ot any interruption, it should complete the ongoing request but should not take new requests
	done := make(chan os.Signal ,1)
	signal.Notify(done, os.Interrupt,syscall.SIGINT,syscall.SIGTERM)

	go func(){
		err := server.ListenAndServe() //this is blocking
		if err!=nil{
			log.Fatal("failed to start server")
		}
	}()
	
	<-done 
	
	//here we will write the logic of server stop
	slog.Info("shutting down the server")

	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err!=nil {
		slog.Error("failed to shutdown server",slog.String("error",err.Error()))
	}

	slog.Info("server shutdown successfully")
}
//command to start : go run cmd/students-api/main.go -config config/local.yaml