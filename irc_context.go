package main

import (
	"fmt"
	"log"
	"net"

	"github.com/nlopes/slack"
)

// Channel represents an IRC channel. It maps to Slack's groups and channels.
// Private messages are handled differently.
type Channel struct {
	Members []string
	Topic   string
	// Slack groups are different from channels. Here I try to uniform them for
	// IRC, but I still need to know which is which to use the right API calls.
	IsGroup bool
}

// MembersDiff compares the members of this channel with another members list
// and return a slice of members who joined and a slice of members who left.
func (c Channel) MembersDiff(otherMembers []string) ([]string, []string) {
	var membersMap = map[string]bool{}
	for _, m := range c.Members {
		membersMap[m] = true
	}
	var otherMembersMap = map[string]bool{}
	for _, m := range otherMembers {
		otherMembersMap[m] = true
	}

	added := make([]string, 0)
	for _, m := range otherMembers {
		if _, ok := membersMap[m]; !ok {
			added = append(added, m)
		}
	}

	removed := make([]string, 0)
	for _, m := range c.Members {
		if _, ok := otherMembersMap[m]; !ok {
			removed = append(removed, m)
		}
	}
	return added, removed
}

// IrcContext holds the client context information
type IrcContext struct {
	Conn           *net.TCPConn
	Nick           string
	UserName       string
	RealName       string
	SlackClient    *slack.Client
	SlackAPIKey    string
	SlackConnected bool
	ServerName     string
	Channels       map[string]Channel
	Users          []slack.User
}

// GetUsers returns a list of users of the Slack team the context is connected
// to
func (ic *IrcContext) GetUsers(refresh bool) []slack.User {
	if refresh || ic.Users == nil || len(ic.Users) == 0 {
		users, err := ic.SlackClient.GetUsers()
		if err != nil {
			log.Printf("Failed to get users: %v", err)
			return nil
		}
		ic.Users = users
		log.Printf("Fetched %v users", len(users))
	}
	return ic.Users
}

// GetUserInfo returns a slack.User instance from a given user ID, or nil if
// no user with that ID was found
func (ic *IrcContext) GetUserInfo(userID string) *slack.User {
	users := ic.GetUsers(false)
	if users == nil || len(users) == 0 {
		return nil
	}
	// XXX this may be slow, convert user list to map?
	for _, user := range users {
		if user.ID == userID {
			return &user
		}
	}
	return nil
}

// GetUserInfoByName returns a slack.User instance from a given user name, or
// nil if no user with that name was found
func (ic *IrcContext) GetUserInfoByName(username string) *slack.User {
	users := ic.GetUsers(false)
	if users == nil || len(users) == 0 {
		return nil
	}
	for _, user := range users {
		if user.Name == username {
			return &user
		}
	}
	return nil
}

// Mask returns the IRC mask for the current user
func (ic IrcContext) Mask() string {
	var username string
	if ic.UserName == "" {
		user := ic.GetUserInfo(ic.Nick)
		if user == nil {
			username = "unknown"
		} else {
			username = user.ID
		}
	}
	return fmt.Sprintf("%v!%v@%v", ic.Nick, username, ic.Conn.RemoteAddr().(*net.TCPAddr).IP)
}

// UserIDsToNames returns a list of user names corresponding to a list of user
// IDs. If an ID is unknown, it is returned unmodified in the output list
func (ic IrcContext) UserIDsToNames(userIDs ...string) []string {
	var names []string
	// TODO implement using ic.GetUsers() instead
	for _, uid := range userIDs {
		user, err := ic.SlackClient.GetUserInfo(uid)
		if err != nil {
			names = append(names, uid)
		} else {
			names = append(names, user.Name)
		}
	}
	return names
}

// Maps of user contexts and nicknames
var (
	UserContexts  = map[net.Addr]*IrcContext{}
	UserNicknames = map[string]*IrcContext{}
)