package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SmartMeshFoundation/Atmosphere/rest"
	"github.com/SmartMeshFoundation/Atmosphere/service"
	"github.com/urfave/cli"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "sm-client",
			Usage: "The smartraiden restful host, Default http://localhost:5001",
			Value: "http://localhost:5001",
		},
		cli.StringFlag{
			Name:  "lnd-client",
			Usage: "The lnd rpc host, Default localhost:10001",
			Value: "localhost:10009",
		},
		cli.StringFlag{
			Name:  "listen",
			Usage: "listen at, Default localhost:7001",
			Value: "localhost:7001",
		},
	}
	app.Action = mainCtx
	app.Name = "atmosphere"
	app.Version = "0.1"
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("exit with err :", err.Error())
	}
}

func mainCtx(ctx *cli.Context) (err error) {
	log.Printf("Welcom to atomsphere,version %s\n", ctx.App.Version)
	service.InitAPI(ctx.String("sm-client"), ctx.String("lnd-client"))
	rest.Start(ctx.String("listen"))
	return
}
