package utils

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"slices"
	"werewolves-go/data"
	"werewolves-go/types"

	"github.com/anthdm/hollywood/actor"
)

/*
 * Checks if werewolves are alive.
 */
func AreWerewolvesAlive(users map[string]*data.Client) bool {
	for _, user := range users {
		if user.Status && user.Role == "werewolf" {
			return true
		}
	}

	return false
}

/*
 * Checks if witch is alive.
 */
func IsWitchAlive(users map[string]*data.Client) bool {
	for _, user := range users {
		if user.Status && user.Role == "witch" {
			return true
		}
	}

	return false
}

/*
 * Checks if a particular user is alive.
 */
func IsUserAlive(user *data.Client) bool {
	return user.Status
}

/*
 * Returns count of alive werewolves.
 */
func CountWerewolvesAlive(users map[string]*data.Client) int {
	count := 0
	for _, user := range users {
		if user.Status && user.Role == "werewolf" {
			count++
		}
	}

	return count
}

/*
 * Checks if townpersons are alive.
 */
func AreTownspersonAlive(users map[string]*data.Client) bool {
	for _, user := range users {
		if user.Status && (user.Role == "townsperson" || user.Role == "witch") {
			return true
		}
	}

	return false
}

/*
 * Returns count of how many townpersons are alive.
 */
func CountUsersAlive(users map[string]*data.Client) int {
	count := 0
	for _, user := range users {
		if user.Status {
			count++
		}
	}

	return count
}

/*
 * Returns pidList of alive werewolves for communication.
 */
func GetAliveWerewolves(users map[string]*data.Client, clients map[string]*actor.PID) []*actor.PID {

	var pidList []*actor.PID
	for cAddr, data := range users {
		if data.Role == "werewolf" && data.Status {
			pidList = append(pidList, clients[cAddr])
		}
	}

	return pidList
}

/*
 * Returns pid of witch (if alive) for communication.
 */
func GetAliveWitch(users map[string]*data.Client, clients map[string]*actor.PID) *actor.PID {

	var pid *actor.PID
	for cAddr, data := range users {
		if data.Role == "witch" && data.Status {
			pid = clients[cAddr]
		}
	}

	return pid
}

/*
 * Send identities of clients (werewolve / townsperson).
 */
func SendIdentities(users map[string]*data.Client, clients map[string]*actor.PID, ctx *actor.Context) {
	for caddr, pid := range clients {
		role := users[caddr].Role
		ctx.Send(pid, FormatMessageResponseFromServer(fmt.Sprintf("========== You are a %v =========", role)))
	}
}

/*
 * Get PID list of alive users.
 */
func GetAliveTownperson(users map[string]*data.Client, clients map[string]*actor.PID) []*actor.PID {

	var pidList []*actor.PID
	for cAddr, data := range users {
		if data.Status {
			pidList = append(pidList, clients[cAddr])
		}
	}

	return pidList
}

/*
 * Returns list of usernames that are alive.
 */
func GetListofUsernames(users map[string]*data.Client) []string {
	var userList []string
	for _, data := range users {
		if data.Status {
			userList = append(userList, data.Name)
		}
	}

	return userList
}

// Get address of the client based on the username for the client.
func GetCAddrFromUsername(users map[string]*data.Client, username string) string {
	for caddr, user := range users {
		if user.Name == username {
			return caddr
		}
	}

	return ""
}

// Return messgae formatted as being sent by server.
func FormatMessageResponseFromServer(message string) *types.Message {
	msgResponse := &types.Message{
		Username: "server/primary",
		Msg:      message,
	}

	return msgResponse
}

// Check if username is allowed to be used.
func IsUsernameAllowed(username string, users map[string]*data.Client) bool {
	for _, user := range users {
		if user.Name == username {
			return true
		}
	}

	return false
}


// Set up roles before initiating the game.
func SetUpRoles(users map[string]*data.Client, witches map[string]*data.Client, werewolves map[string]*data.Client, number_of_werewolves int) {
	//create a list of unique random numbers. length of list is equal to the number of werewolves you want in the
	//game.
	var listRand []int
	for i := 0; i < number_of_werewolves; {
		randNum := rand.IntN(10000) % len(users)
		if slices.Contains(listRand, randNum) {
			//randNum already present in our list - so re-run the random number generation again
		} else {
			listRand = append(listRand, randNum)
			i++
		}
	}

	// assign the witch role to a player - checks whether that player is already assigned to be a werewolf.
	var witchRand int = rand.IntN(10000) % len(users)
	for {
		if slices.Contains(listRand, witchRand) {
			//witchRand is a part of the werewolves list - so get a different random number
			witchRand = rand.IntN(10000) % len(users)
		} else {
			break
		}
	}

	user_names := GetListofUsernames(users)
	slog.Info(fmt.Sprintf("listRand = %v :These indices in %v will be the werewolves.\n", listRand, user_names))

	for i, userName := range user_names {
		var assignWerewolf = false
		var assignWitch = false
		for _, randNum := range listRand {
			if randNum == i {
				assignWerewolf = true
			}
		}
		if witchRand == i {
			assignWitch = true
		}

		caddr := GetCAddrFromUsername(users, user_names[i])
		if assignWerewolf {
			if users[caddr].Role == "" { //performing an additional sanity check with this line
				if entry, ok := users[caddr]; ok {
					entry.Role = "werewolf"
					users[caddr] = entry
					werewolves[caddr] = entry
				}
				slog.Info(userName + " has been assigned to be a werewolf")
			}
		} else if assignWitch {
			if users[caddr].Role == "" { //performing an additional sanity check with this line
				if entry, ok := users[caddr]; ok {
					entry.Role = "witch"
					users[caddr] = entry
					witches[caddr] = entry
				}
				slog.Info(userName + " has been assigned to be a witch")
			}
		} else {
			if users[caddr].Role == "" { //performing an additional sanity check with this line
				if entry, ok := users[caddr]; ok {
					entry.Role = "townsperson"
					users[caddr] = entry
				}
				slog.Info(userName + " has been assigned to be a townsperson")
			}
		}
	}
}

// Print all users.
func PrintUsers(users map[string]*data.Client) {
	fmt.Println("Print Users")
	for user, data := range users {
		fmt.Printf(
			"%v has been assigned %v username, has role %v with alive status %v\n",
			user,
			data.Name,
			data.Role,
			data.Status)
	}
}
