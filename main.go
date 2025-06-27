package main

import (
	"fmt"

	a "github.com/y3933y3933/joker/internal/app"
	"github.com/y3933y3933/joker/internal/routes"
)

func main() {
	app, err := a.NewApplication()
	if err != nil {
		panic(err)
	}

	router := routes.SetupRoutes(app)
	port := fmt.Sprintf(":%d", app.Config.Port)
	router.Run(port)

}
