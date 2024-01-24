package main

import (
	"fmt"
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"lk-service/internal/api_db/reindexer_db"
	"lk-service/internal/handler"
	"lk-service/internal/service"
	"log"
	"resenje.org/logging"
)

func main() {
	viper.SetConfigFile("config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		logging.Info(fmt.Errorf("config: %v", err))
	}

	db := reindexer.NewReindex(fmt.Sprintf("%s://%s:%s/%s",
		viper.GetString("db.scheme"),
		viper.GetString("db.hostname"),
		viper.GetString("db.port"),
		viper.GetString("db.path"),
	))

	if err := db.Status().Err; err != nil {
		logging.Info(fmt.Errorf("reindexer connection: %v", err))
	}

	logging.Info("Connection to reindexer DB successful!")

	userApiDB := reindexer_db.NewUserApiDB(db)
	generalApiDB := reindexer_db.NewGeneralApiDB(db)
	emailCodeApiDB := reindexer_db.NewEmailCodeApiDB(db)

	s := service.NewService(userApiDB, generalApiDB, emailCodeApiDB)

	h := handler.NewHandler(s)

	server := fasthttp.Server{
		Handler: h.InitRoutes,
	}

	if err := server.ListenAndServe(viper.GetString("http.port")); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
