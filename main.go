package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-trip/trip"
	"log"
	"os"
)

func main() {

	log.SetFlags(log.Llongfile | log.LstdFlags)

	env := "./config/env.ini"

	if len(os.Args) > 1 {
		env = os.Args[1]
	}

	a := trip.TripApp{}

	err := app.Load(&a, "./app.ini")

	if err != nil {
		log.Panicln(err)
	}

	err = app.Load(&a, env)

	if err != nil {
		log.Panicln(err)
	}

	app.Obtain(&a)

	_, err = a.GetDB()

	if err != nil {
		log.Println(err)
	}

	log.Println(a.RouteTable)

	app.Handle(&a, &app.InitTask{})

	if a.Runloop {
		app.Handle(&a, &app.RunloopTask{})
	}

	kk.DispatchMain()

	app.Recycle(&a)

}
