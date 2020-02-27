package main

import (
	"errors"
	"github.com/rivo/tview"
	"sort"
	"strconv"
	"sync"
)

type UserList struct {
	Container *tview.Table
	Users map[string]*UserListItem
	sync.Mutex
}

type UserListItem struct {
	Online      bool
	NewMessages int
	Group       bool
}

func NewUserList(ulContainer *tview.Table) *UserList {
	return &UserList{
		Users: make(map[string]*UserListItem),
		Mutex: sync.Mutex{},
		Container: ulContainer,
	}
}

func (u *UserList) PopulateFromDB() {
	clientUsername, err := client.Cache.Get("username")
	if err != nil {
		return
	}
	users, err := dbUserList()
	if err != nil {
		app.Stop()
	}
	// List all of our contacts in our local DB
	for _, user := range users {
		if user.Username != clientUsername {
			u.AddUser(user.Username)
		}
	}
}

func (u *UserList) AddUser(username string) {
	u.Lock()
	defer u.Unlock()
	u.Users[username] = &UserListItem{
		Online:      false,
		NewMessages: 0,
		Group:       false,
	}
	u.Draw()
}

func (u *UserList) RemoveUser(username string) {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.Users[username]; ok {
		delete(u.Users, username)
	}
	u.Draw()
}

func (u *UserList) NewMessage(username string) {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.Users[username]; ok {
		if inboxSelectedUsername != username {
			u.Users[username].NewMessages += 1
		}
	}
	u.Draw()
}

func (u *UserList) ClearNewMessage(username string) {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.Users[username]; ok {
		u.Users[username].NewMessages = 0
	}
	u.Draw()
}

func (u *UserList) SetUserOnline(username string) {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.Users[username]; ok {
		u.Users[username].Online = true
	}
	u.Draw()
}

func (u *UserList) SetUserOffline(username string) {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.Users[username]; ok {
		u.Users[username].Online = false
	}
	u.Draw()
}

func (u *UserList) Get(username string) (string, error) {
	if _, ok := u.Users[username]; !ok {
		return "", errors.New("user does not exist in user list")
	}
	// Define our base var
	fUsername := ""

	// Selected status
	if username == inboxSelectedUsername {
		fUsername += "-"
	}

	// Handle online status
	if u.Users[username].Online {
		fUsername += "[green]"
	} else {
		fUsername += "[red]"
	}

	// Add our username string
	fUsername += username

	// Add number of missed messages
	if u.Users[username].NewMessages > 0 {
		total := strconv.Itoa(u.Users[username].NewMessages)
		fUsername += " (" + total + ")"
	}

	return fUsername, nil
}

func (u *UserList) Draw() {
	app.QueueUpdateDraw(func() {
		// Clear our current list
		userListContainer.Clear()
		// Create sort array
		sortList := make([]string, 0)
		for username, _ := range u.Users {
			sortList = append(sortList, username)
		}
		// Sort the array
		sort.Strings(sortList)
		// Draw
		count := 0
		for _, username := range sortList {
			u, err := u.Get(username)
			if err != nil {
				continue
			}
			userListContainer.SetCell(count, 0, tview.NewTableCell(u))
			count++
		}
	})
}
