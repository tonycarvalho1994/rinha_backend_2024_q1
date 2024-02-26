package main

import (
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/api/handler"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/service"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/infra/database"
	"log"
	"net/http"
)

func SetupCustomer(service *service.CustomerService) {
	_ = service.Repository.CreateCustomer("1", 100000)
	_ = service.Repository.CreateCustomer("2", 80000)
	_ = service.Repository.CreateCustomer("3", 1000000)
	_ = service.Repository.CreateCustomer("4", 10000000)
	_ = service.Repository.CreateCustomer("5", 500000)
}

func main() {
	//mongoUri := os.Getenv("MONGO_URI")
	//databaseName := os.Getenv("DATABASE_NAME")
	//collectionName := os.Getenv("COLLECTION_NAME")
	//ctx := context.Background()
	//collection, client, err := database.CreateConnection("mongodb://"+mongoUri, databaseName, collectionName, ctx)
	//defer func() {
	//	if err = client.Disconnect(context.Background()); err != nil {
	//		log.Fatal(err)
	//	}
	//}()
	//
	//repository := database.CustomerRepositoryMongo{Collection: collection, Ctx: ctx}

	psqlDb, err := database.CreateDatabasePSQL()
	if err != nil {
		log.Fatalf(err.Error())
	}
	repository := database.CustomerRepositoryPSQL{Db: psqlDb}

	customerService := service.CustomerService{
		Repository: &repository,
	}

	SetupCustomer(&customerService)
	log.Println("Customers created.")

	server := handler.NewServer(&customerService)
	server.SetupRouter()

	log.Println("Starting server...")
	err = http.ListenAndServe("0.0.0.0:8080", server.Mux)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
