package tests

import (
	"log"
	"os"
	"testing"

	"drexel.edu/todo/db"
	fake "github.com/brianvoe/gofakeit/v6" //aliasing package name
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

var (
	BASE_API = "http://localhost:1080"

	cli = resty.New()
)

func TestMain(m *testing.M) {

	//SETUP GOES FIRST
	rsp, err := cli.R().Delete(BASE_API + "/voter")

	if rsp.StatusCode() != 200 {
		log.Printf("error clearing database, %v", err)
		os.Exit(1)
	}

	code := m.Run()

	//CLEANUP

	//Now Exit
	os.Exit(code)
}

func newRandVoter(id uint) db.Voter {
	return db.Voter{
		VoterId: id,
		Name:    fake.Name(),
		Email:   fake.Email(),
	}
}

func Test_LoadDB(t *testing.T) {
	numLoad := 3
	for i := 0; i < numLoad; i++ {
		item := newRandVoter(uint(i))
		rsp, err := cli.R().
			SetBody(item).
			Post(BASE_API + "/voter")

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())
	}
}

func Test_GetAllVoters(t *testing.T) {
	var items []db.Voter

	rsp, err := cli.R().SetResult(&items).Get(BASE_API + "/voter")

	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode())

	assert.Equal(t, 3, len(items))
}

func Test_DeleteVoter(t *testing.T) {
	var item db.Voter

	rsp, err := cli.R().SetResult(&item).Get(BASE_API + "/voter/2")
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "voter #2 expected")

	rsp, err = cli.R().Delete(BASE_API + "/voter/2")
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "voter not deleted expected")

	rsp, err = cli.R().SetResult(item).Get(BASE_API + "/voter/2")
	assert.Nil(t, err)
	assert.Equal(t, 404, rsp.StatusCode(), "expected not found error code")
}
