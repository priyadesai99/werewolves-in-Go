package data

import (
	"fmt"
	"math"
)

/*
 * Struct to count votes.
 */
type Voters struct {
	user_vote   map[string]int
	voted_users map[string]bool
}

/*
 * New voters struct initialization.
 */
func NewVoters(users []string) *Voters {
	v := &Voters{
		user_vote:   make(map[string]int),
		voted_users: make(map[string]bool),
	}

	for _, user := range users {
		v.user_vote[user] = 0
	}

	return v
}

/*
 * Returns user that has been voted maximum times.
 * Returns a single possible. 2 or more players cannot be killed at the same time.
 */
func (voters *Voters) GetMaxVotedUser() string {
	var dead_users []string
	max_votes := math.MinInt

	for _, vote := range voters.user_vote {
		if vote > max_votes {
			max_votes = vote
		}
	}

	for user, vote := range voters.user_vote {
		if vote == max_votes {
			dead_users = append(dead_users, user)
		}
	}
	if len(dead_users) > 1 || len(dead_users) == 0 {
		return ""
	} else {
		return dead_users[0]
	}
}

// Add vote the votes list.
func (voters *Voters) AddVote(user string, sender string) bool {
	_, ok := voters.user_vote[user]
	if ok && !voters.voted_users[sender] {
		voters.user_vote[user]++
		voters.voted_users[sender] = true
		return true
	} else {
		return false
	}
}

// Clear the votes to reuse voting object.
func (voter *Voters) ClearVotes() {
	for user := range voter.user_vote {
		voter.user_vote[user] = 0
	}

	voter.voted_users = nil
}

// Print user votes.
func (voter *Voters) PrintVotes() {
	fmt.Println("VOTES")
	for user, vote := range voter.user_vote {
		fmt.Printf("%v has %v votes\n", user, vote)
	}
}
