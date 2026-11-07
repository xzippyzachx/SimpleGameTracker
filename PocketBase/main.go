package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		registerIgdbRoutes(se)
		registerSyncRoutes(se)
		return se.Next()
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		if publicDir := os.Getenv("PB_PUBLIC_DIR"); publicDir != "" {
			se.Router.GET("/{path...}", apis.Static(os.DirFS(publicDir), true))
		}
		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
