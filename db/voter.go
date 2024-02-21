package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

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
}

func NewVoterList() (*VoterList, error) {

	voterList := &VoterList{
		Voters: make(map[uint]Voter),
	}

	return voterList, nil
}

func (v *VoterList) AddVoter(voter Voter) error {

	_, ok := v.Voters[voter.VoterId]
	if ok {
		return errors.New("voter already exists")
	}

	v.Voters[voter.VoterId] = voter

	return nil
}

func (v *VoterList) DeleteVoter(id uint) error {

	delete(v.Voters, id)

	return nil
}

func (v *VoterList) DeleteAll() error {

	v.Voters = make(map[uint]Voter)

	return nil
}

func (v *VoterList) UpdateVoter(voter Voter) error {

	_, ok := v.Voters[voter.VoterId]
	if !ok {
		return errors.New("voter does not exist")
	}

	v.Voters[voter.VoterId] = voter

	return nil
}

func (v *VoterList) GetVoter(id uint) (Voter, error) {

	voter, ok := v.Voters[id]
	if !ok {
		return Voter{}, errors.New("voter does not exist")
	}

	return voter, nil
}

func (v *VoterList) GetVoteHistory(id uint) ([]VoterHistory, error) {

	voter, ok := v.Voters[id]
	if !ok {
		return []VoterHistory{}, errors.New("voter history does not exist")
	}

	return voter.VoteHistory, nil
}

func (v *VoterList) GetSingleVoteHistory(voterId uint, pollId uint) (*VoterHistory, error) {

	voter, ok := v.Voters[voterId]
	if !ok {
		return &VoterHistory{}, errors.New("that voter does not exist")
	}

	for _, vote := range voter.VoteHistory {
		if vote.PollId == pollId {
			return &vote, nil
		}
	}

	return nil, errors.New("poll does not exist for the specified voter")
}

func (v *VoterList) AddPoll(voterId uint, poll VoterHistory) error {

	voter, ok := v.Voters[voterId]
	if !ok {
		return errors.New("voter not found")
	}

	voter.VoteHistory = append(voter.VoteHistory, poll)
	v.Voters[voterId] = voter

	return nil

}

func (v *VoterList) GetAllVoters() ([]Voter, error) {

	var voterList []Voter

	for _, voter := range v.Voters {
		voterList = append(voterList, voter)
	}

	return voterList, nil
}

func (v *VoterList) PrintVoter(voter Voter) {
	jsonBytes, _ := json.MarshalIndent(voter, "", "  ")
	fmt.Println(string(jsonBytes))
}

func (v *VoterList) PrintAllVoters(voterList []Voter) {
	for _, voter := range voterList {
		v.PrintVoter(voter)
	}
}

func (v *VoterList) JsonToVoter(jsonString string) (Voter, error) {
	var voter Voter
	err := json.Unmarshal([]byte(jsonString), &voter)
	if err != nil {
		return Voter{}, err
	}

	return voter, nil
}
