package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"drexel.edu/voter/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Global variables to hold the command line flags to drive the voter CLI
// application
var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {

	//Note some networking lingo, some frameworks start the server on localhost
	//this is a local-only interface and is fine for testing but its not accessible
	//from other machines.  To make the server accessible from other machines, we
	//need to listen on an interface, that could be an IP address, but modern
	//cloud servers may have multiple network interfaces for scale.  With TCP/IP
	//the address 0.0.0.0 instructs the network stack to listen on all interfaces
	//We set this up as a flag so that we can overwrite it on the command line if
	//needed
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

// main is the entry point for our voter API application.  It processes
// the command line flags and then uses the db package to perform the
// requested operation
func main() {
	processCmdLineFlags()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())

	apiHandler, err := api.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//HTTP Standards for "REST" APIS
	//GET - Read/Query
	//POST - Create
	//PUT - Update
	//DELETE - Delete

	app.Get("/voter", apiHandler.ListAllVoters)
	app.Post("/voter", apiHandler.AddVoter)
	app.Put("/voter", apiHandler.UpdateVoter)
	app.Delete("/voter", apiHandler.DeleteAllVoters)
	app.Delete("/voter/:id<int>", apiHandler.DeleteVoter)
	app.Get("/voter/:id<int>", apiHandler.GetVoter)

	app.Get("/voter/:id<int>/polls", apiHandler.GetPollHistoryFromVoter)
	app.Get("/voter/:id<int>/polls/:pollid", apiHandler.GetSinglePollFromVoter)
	app.Post("/voter/:id<int>", apiHandler.AddSinglePollToVoter)

	app.Get("/crash", apiHandler.CrashSim)
	app.Get("/crash2", apiHandler.CrashSim2)
	app.Get("/crash3", apiHandler.CrashSim3)
	app.Get("/health", apiHandler.HealthCheck)

	//We will now show a common way to version an API and add a new
	//version of an API handler under /v2.  This new API will support
	//a path parameter to search for voters based on a status
	// v2 := app.Group("/v2")
	// v2.Get("/voter", apiHandler.ListSelectVoters)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	log.Println("Starting server on ", serverPath)
	app.Listen(serverPath)
}
