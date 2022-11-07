package main

import (
	"challenge/internal/api"
	"challenge/internal/cache"
	"challenge/internal/database"
	gorm "challenge/internal/database/gorm"
	"challenge/internal/logger"
	"challenge/internal/middleware"
	"challenge/internal/repository"
	"net/http"
	"os"
)
const(
	DbType = "DATABASE_TYPE"
	Gorm = "gorm"
)

func main() {
	var repo repository.Repository
	var err error

	_ = logger.InitLogger()
	dBIMP, ok := os.LookupEnv(DbType)

	if ok && dBIMP == Gorm {
		repo, err = gorm.New()
		if err != nil {
			logger.Log().Panic(err)
		}

	} else {
		repo, err = database.New()
		if err != nil {
			logger.Log().Panic(err)
		}
	}

	cache, err := cache.NewClient()
	if err != nil{
		logger.Log().Panic(err)
	}
	
	router := api.New(repo, cache)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", router.HandlerList)
	mux.HandleFunc("/tasks/", router.HandlerGetTask)

	wrappedMux := middleware.NewContentMiddleware(mux) 

	logger.Log().Fatal(http.ListenAndServe(":8080", wrappedMux))

}
