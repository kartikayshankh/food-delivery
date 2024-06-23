package main

import (
	configurations "github.com/kartikayshankh/food-delivery/config"
	"github.com/kartikayshankh/food-delivery/db"
	"github.com/kartikayshankh/food-delivery/handler"
	"github.com/kartikayshankh/food-delivery/service/health"
	restaurant "github.com/kartikayshankh/food-delivery/service/restaurant"
	"github.com/kartikayshankh/food-delivery/service/rider"

	"context"
	"log"

	"github.com/kartikayshankh/food-delivery/service/user"
)

func main() {
	config := configurations.Init("config")
	port := ":" + config.GetString("service.port")

	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	//mongodb
	client, errClient := db.InitMongoClient(config)
	if errClient != nil {
		log.Println("errClient:", errClient)
	}
	defer client.Disconnect(context.Background())

	healthService := health.NewService()
	userService := user.NewUserService(config, client)
	RiderService := rider.NewRiderService(config, client)
	RestaurantService := restaurant.NewRestaurantService(config, client)

	// Start server
	e := handler.MakeHTTPHandler(config, healthService, userService, RiderService, RestaurantService)
	e.Logger.Fatal(e.Start(port))
}
