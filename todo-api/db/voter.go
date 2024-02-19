package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// // ToDoItem is the struct that represents a single ToDo item
// type VoterItem struct {
// 	Id     int    `json:"id"`
// 	Title  string `json:"title"`
// 	IsDone bool   `json:"done"`
// }

// // DbMap is a type alias for a map of ToDoItems.  The key
// // will be the ToDoItem.Id and the value will be the ToDoItem
//type DbMap map[int]Voter

// ToDo is the struct that represents the main object of our
// todo app.  It contains a map of ToDoItems and the name of
// the file that is used to store the items.
//
// This is just a mock, so we will only be managing an in memory
// map
// type Votermap struct {
// 	voterMap DbMap
// 	//more things would be included in a real implementation
// }

type VoterHistory struct {
	PollId   uint
	VoteId   uint
	VoteDate time.Time
}

type Voter struct {
	VoterId     uint
	Name        string
	Email       string
	VoteHistory []VoterHistory
}

type VoterList struct {
	Voters map[uint]Voter
	//more things would be included in a real implementation
}

// New is a constructor function that returns a pointer to a new
// ToDo struct.  It takes a single string argument that is the
// name of the file that will be used to store the ToDo items.
// If the file doesn't exist, it will be created.  If the file
// does exist, it will be loaded into the ToDo struct.
func NewVoterList() (*VoterList, error) {

	//Now that we know the file exists, at at the minimum we have
	//a valid empty DB, lets create the ToDo struct
	voterList := &VoterList{
		Voters: make(map[uint]Voter),
	}

	// We should be all set here, the ToDo struct is ready to go
	// so we can support the public database operations
	return voterList, nil
}

//------------------------------------------------------------
// THESE ARE THE PUBLIC FUNCTIONS THAT SUPPORT OUR TODO APP
//------------------------------------------------------------

// AddItem accepts a ToDoItem and adds it to the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must not already exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if so, return an error
//
// Postconditions:
//
//	 (1) The item will be added to the DB
//		(2) The DB file will be saved with the item added
//		(3) If there is an error, it will be returned
func (v *VoterList) AddVoter(voter Voter) error {

	//Before we add an item to the DB, lets make sure
	//it does not exist, if it does, return an error
	_, ok := v.Voters[voter.VoterId]
	if ok {
		return errors.New("voter already exists")
	}

	//Now that we know the item doesn't exist, lets add it to our map
	v.Voters[voter.VoterId] = voter

	//If everything is ok, return nil for the error
	return nil
}

// DeleteItem accepts an item id and removes it from the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be removed from the DB
//		(2) The DB file will be saved with the item removed
//		(3) If there is an error, it will be returned
func (v *VoterList) DeleteVoter(id uint) error {

	// we should if item exists before trying to delete it
	// this is a good practice, return an error if the
	// item does not exist

	//Now lets use the built-in go delete() function to remove
	//the item from our map
	delete(v.Voters, id)

	return nil
}

// DeleteAll removes all items from the DB.
// It will be exposed via a DELETE /todo endpoint
func (v *VoterList) DeleteAll() error {
	//To delete everything, we can just create a new map
	//and assign it to our existing map.  The garbage collector
	//will clean up the old map for us
	v.Voters = make(map[uint]Voter)

	return nil
}

// UpdateItem accepts a ToDoItem and updates it in the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be updated in the DB
//		(2) The DB file will be saved with the item updated
//		(3) If there is an error, it will be returned
func (v *VoterList) UpdateVoter(voter Voter) error {

	// Check if item exists before trying to update it
	// this is a good practice, return an error if the
	// item does not exist
	_, ok := v.Voters[voter.VoterId]
	if !ok {
		return errors.New("voter does not exist")
	}

	//Now that we know the item exists, lets update it
	v.Voters[voter.VoterId] = voter

	return nil
}

// GetItem accepts an item id and returns the item from the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be returned, if it exists
//		(2) If there is an error, it will be returned
//			along with an empty ToDoItem
//		(3) The database file will not be modified
func (v *VoterList) GetVoter(id uint) (Voter, error) {

	// Check if item exists before trying to get it
	// this is a good practice, return an error if the
	// item does not exist
	voter, ok := v.Voters[id]
	if !ok {
		return Voter{}, errors.New("voter does not exist")
	}

	return voter, nil
}

// ChangeItemDoneStatus accepts an item id and a boolean status.
// It returns an error if the status could not be updated for any
// reason.  For example, the item itself does not exist, or an
// IO error trying to save the updated status.

// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The items status in the database will be updated
//		(2) If there is an error, it will be returned.
//		(3) This function MUST use existing functionality for most of its
//			work.  For example, it should call GetItem() to get the item
//			from the DB, then it should call UpdateItem() to update the
//			item in the DB (after the status is changed).
func (v *VoterList) ChangeVoterDoneStatus(id uint, value bool) error {

	//update was successful
	return errors.New("not implemented")
}

// GetAllItems returns all items from the DB.  If successful it
// returns a slice of all of the items to the caller
// Preconditions:   (1) The database file must exist and be a valid
//
// Postconditions:
//
//	 (1) All items will be returned, if any exist
//		(2) If there is an error, it will be returned
//			along with an empty slice
//		(3) The database file will not be modified
func (v *VoterList) GetAllVoters() ([]Voter, error) {

	//Now that we have the DB loaded, lets crate a slice
	var voterList []Voter

	//Now lets iterate over our map and add each item to our slice
	for _, voter := range v.Voters {
		voterList = append(voterList, voter)
	}

	//Now that we have all of our items in a slice, return it
	return voterList, nil
}

// PrintItem accepts a ToDoItem and prints it to the console
// in a JSON pretty format. As some help, look at the
// json.MarshalIndent() function from our in class go tutorial.
func (v *VoterList) PrintVoter(voter Voter) {
	jsonBytes, _ := json.MarshalIndent(voter, "", "  ")
	fmt.Println(string(jsonBytes))
}

// PrintAllItems accepts a slice of ToDoItems and prints them to the console
// in a JSON pretty format.  It should call PrintItem() to print each item
// versus repeating the code.
func (v *VoterList) PrintAllVoters(voterList []Voter) {
	for _, voter := range voterList {
		v.PrintVoter(voter)
	}
}

// JsonToItem accepts a json string and returns a ToDoItem
// This is helpful because the CLI accepts todo items for insertion
// and updates in JSON format.  We need to convert it to a ToDoItem
// struct to perform any operations on it.
func (v *VoterList) JsonToVoter(jsonString string) (Voter, error) {
	var voter Voter
	err := json.Unmarshal([]byte(jsonString), &voter)
	if err != nil {
		return Voter{}, err
	}

	return voter, nil
}
