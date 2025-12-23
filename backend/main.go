package main

import (
	"log"
	"yangdongju/gtd_todo/internal/config"
	"yangdongju/gtd_todo/internal/db"
	"yangdongju/gtd_todo/internal/server"
)

func main() {
	cfg := config.Load()
	pool, err := db.NewConnectionPool(cfg)
	if err != nil {
		log.Fatalf("Database connection failed.\n%v", err)
	}
	server.SetupRouter(pool)
}

