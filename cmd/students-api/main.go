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

	//starting server
	go func(){
		err := server.ListenAndServe() //this is blocking
		if err!=nil{
			log.Fatal("failed to start server")
		}
	}()
	
	<-done //once any kind of interrupt signal is received, channel sends and this lines allows for the code to execute ahead i.e shutdown logic, till then it blocks the code
	
	//here we will write the logic of server stop i.e shutdown logic
	slog.Info("shutting down the server")

	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	//shutdown can alone do the job of graceful shutdown, but we use context, why?
	//because : reason written down
	err := server.Shutdown(ctx)
	if err!=nil && err != http.ErrServerClosed{
		slog.Error("failed to shutdown server",slog.String("error",err.Error()))
	}

	slog.Info("server shutdown successfully")
}
//command to start : go run cmd/students-api/main.go -config config/local.yaml


//important concept
// server.Shutdown() by itself already does a graceful shutdown. But the context with a timeout adds an important safety net. Let me break it down:
// 🔑 What server.Shutdown() does
// It stops accepting new connections.
// It waits for existing connections (in-flight requests) to finish.
// It only returns after all requests are done OR the context passed to it is canceled.
// So if you just do:
// server.Shutdown(context.Background())
// then it will wait forever until all requests are done.

// ⚠️ Problem
// What if:
// A client keeps a connection open indefinitely?
// Some request hangs (e.g., DB query stuck, infinite loop)?
// Then your server never exits — it just hangs on shutdown. Not ideal for deployments, CI/CD pipelines, or Dockerized apps.

// ✅ Why use context.WithTimeout
// By wrapping with a timeout:
// You give in-flight requests a grace period (e.g., 5 seconds).
// After that period, the context is canceled, and Shutdown() will forcefully close remaining connections.
// This guarantees that your process will exit eventually, even if some requests misbehave.

// 📝 Mental model
// Think of it like this:
// server.Shutdown() = “I’ll wait as long as you need”.
// server.Shutdown(ctx with timeout) = “You have 5 seconds to finish, or I’m pulling the plug”.

// Shutdown() alone = graceful but potentially infinite wait.
// Shutdown(ctx with timeout) = graceful with a deadline, ensuring the process won’t hang forever.