package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	Name string
	C    Config
	M    sync.Mutex
)

type Config struct {
	Users    []User   `json:"users"`
	Sources  []string `json:"sources"`
	Port     string   `json:"port"`
	Database string   `json:"database"`
}

type User struct {
	Name  string `json:"name"`
	Hash  string `json:"hash"`
	Admin bool   `json:"admin"`
}

func passwordHash(p string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(h)
}

func UseDefault() error {
	wd, _ := os.Getwd()
	b, err := json.Marshal(Config{
		Users: []User{User{
			Name:  "admin",
			Hash:  passwordHash("admin"),
			Admin: true,
		}},
		Port:     ":8000",
		Database: path.Join(wd, "audiolib.db"),
	})
	if err != nil {
		return err
	}
	return os.WriteFile(Name, b, 0600)
}

func Parse() error {
	b, err := os.ReadFile(Name)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &C)
	if err != nil {
		return err
	}
	return nil
}

func AddUser(name, password string, admin bool) error {
	M.Lock()
	for _, u := range C.Users {
		if strings.ToLower(u.Name) == strings.ToLower(name) {
			return fmt.Errorf("user '%v' already exists", name)
		}
	}
	C.Users = append(C.Users, User{
		Name:  name,
		Hash:  passwordHash(password),
		Admin: admin,
	})
	M.Unlock()
	return nil
}

func LoginSuccess(name, password string) bool {
	var h string
	M.Lock()
	for _, u := range C.Users {
		if strings.ToLower(name) == strings.ToLower(u.Name) {
			h = u.Hash
			break
		}
	}
	M.Unlock()
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(password)) == nil
}

func Save() error {
	M.Lock()
	b, err := json.Marshal(C)
	if err != nil {
		return err
	}
	M.Unlock()
	return os.WriteFile(Name, b, 0600)
}
