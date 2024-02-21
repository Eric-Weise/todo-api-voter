package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"drexel.edu/voter/db"
	"github.com/gofiber/fiber/v2"
)

// The api package creates and maintains a reference to the data handler
// this is a good design practice
type VoterAPI struct {
	db *db.VoterList
}

func New() (*VoterAPI, error) {
	dbHandler, err := db.NewVoterList()
	if err != nil {
		return nil, err
	}

	return &VoterAPI{db: dbHandler}, nil
}

//Below we implement the API functions.  Some of the framework
//things you will see include:
//   1) How to extract a parameter from the URL, for example
//	  the id parameter in /voter/:id
//   2) How to extract the body of a POST request
//   3) How to return JSON and a correctly formed HTTP status code
//	  for example, 200 for OK, 404 for not found, etc.  This is done
//	  using the c.JSON() function
//   4) How to return an error code and abort the request.  This is
//	  done using the c.AbortWithStatus() function

func (v *VoterAPI) ListAllVoters(c *fiber.Ctx) error {

	voterList, err := v.db.GetAllVoters()
	if err != nil {
		log.Println("Error Getting All Items: ", err)
		return fiber.NewError(http.StatusNotFound,
			"Error Getting All Items")
	}
	//Note that the database returns a nil slice if there are no items
	//in the database.  We need to convert this to an empty slice
	//so that the JSON marshalling works correctly.  We want to return
	//an empty slice, not a nil slice. This will result in the json being []
	if voterList == nil {
		voterList = make([]db.Voter, 0)
	}

	return c.JSON(voterList)
}

// implementation for GET /voter/:id
// returns a single voter
func (v *VoterAPI) GetVoter(c *fiber.Ctx) error {

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	//Note that ParseInt always returns an int64, so we have to
	//convert it to an int before we can use it.
	voter, err := v.db.GetVoter(uint(id))
	if err != nil {
		log.Println("Item not found: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	//Git will automatically convert the struct to JSON
	//and set the content-type header to application/json
	return c.JSON(voter)
}

func (v *VoterAPI) GetPollHistoryFromVoter(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	voterHistory, err := v.db.GetVoteHistory(uint(id))
	if err != nil {
		log.Println("Item not found:", err)
		return fiber.NewError(http.StatusNotFound)
	}
	return c.JSON(voterHistory)
}

func (v *VoterAPI) GetSinglePollFromVoter(c *fiber.Ctx) error {
	voterIdStr := c.Params("id")
	pollIdStr := c.Params("pollid")

	voterid, err := strconv.Atoi(voterIdStr)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	pollid, err := strconv.Atoi(pollIdStr)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	poll, err := v.db.GetSingleVoteHistory(uint(voterid), uint(pollid))
	if err != nil {
		log.Println("Item not found:", err)
		return fiber.NewError(http.StatusNotFound)
	}
	return c.JSON(poll)
}

func (v *VoterAPI) AddSinglePollToVoter(c *fiber.Ctx) error {
	voterIdStr := c.Params("id")
	voterId, err := strconv.Atoi(voterIdStr)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	var poll db.VoterHistory

	if err := c.BodyParser(&poll); err != nil {
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := v.db.AddPoll(uint(voterId), poll); err != nil {
		log.Println("Failed to add poll to voter:", err)
		return fiber.NewError(http.StatusNotFound, "Failed to add poll to voter")
	}

	return c.JSON(voterId)
}

// implementation for POST /voter
// adds a new voter
func (v *VoterAPI) AddVoter(c *fiber.Ctx) error {
	var voter db.Voter

	//With HTTP based APIs, a POST request will usually
	//have a body that contains the data to be added
	//to the database.  The body is usually JSON, so
	//we need to bind the JSON to a struct that we
	//can use in our code.
	//This framework exposes the raw body via c.Request.Body
	//but it also provides a helper function BodyParser
	//that will extract the body, convert it to JSON and
	//bind it to a struct for us.  It will also report an error
	//if the body is not JSON or if the JSON does not match
	//the struct we are binding to.
	if err := c.BodyParser(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := v.db.AddVoter(voter); err != nil {
		log.Println("Error adding item: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.JSON(voter)
}

// implementation for PUT /voter
// Web api standards use PUT for Updates
func (v *VoterAPI) UpdateVoter(c *fiber.Ctx) error {
	var voter db.Voter
	if err := c.BodyParser(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := v.db.UpdateVoter(voter); err != nil {
		log.Println("Error updating item: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.JSON(voter)
}

// implementation for DELETE /voter/:id
// deletes a voter
func (v *VoterAPI) DeleteVoter(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := v.db.DeleteVoter(uint(id)); err != nil {
		log.Println("Error deleting item: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).SendString("Delete OK")
}

// implementation for DELETE /voter
// deletes all voters
func (v *VoterAPI) DeleteAllVoters(c *fiber.Ctx) error {

	if err := v.db.DeleteAll(); err != nil {
		log.Println("Error deleting all items: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).SendString("Delete All OK")
}

/*   SPECIAL HANDLERS FOR DEMONSTRATION - CRASH SIMULATION AND HEALTH CHECK */

// implementation for GET /crash
// This simulates a crash to show some of the benefits of the
// gin framework
func (v *VoterAPI) CrashSim(c *fiber.Ctx) error {
	//panic() is go's version of throwing an exception
	//note with recover middleware this will not end program
	panic("Simulating an unexpected crash")
}

func (v *VoterAPI) CrashSim2(c *fiber.Ctx) error {
	//A stupid crash simulation example
	i := 0
	j := 1 / i
	jStr := fmt.Sprintf("%d", j)
	return c.Status(http.StatusOK).
		JSON(fiber.Map{
			"val_j": jStr,
		})
}

func (v *VoterAPI) CrashSim3(c *fiber.Ctx) error {
	//A stupid crash simulation example
	os.Exit(10)
	return c.Status(http.StatusOK).
		JSON(fiber.Map{
			"error": "will never get here, nothing you can do about this",
		})
}

// implementation of GET /health. It is a good practice to build in a
// health check for your API.  Below the results are just hard coded
// but in a real API you can provide detailed information about the
// health of your API with a Health Check
func (v *VoterAPI) HealthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).
		JSON(fiber.Map{
			"status":             "ok",
			"version":            "1.0.0",
			"uptime":             100,
			"users_processed":    1000,
			"errors_encountered": 10,
		})
}
