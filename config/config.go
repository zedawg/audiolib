package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	C Config
	M sync.Mutex
)

type Config struct {
	Name     string           `json:"-"`
	Port     string           `json:"port"`
	Database string           `json:"database"`
	Sources  []string         `json:"sources,omitempty"`
	Users    map[string]*User `json:"users,omitempty"`
}

func (c Config) String() string {
	return fmt.Sprintf("config=%v\ndatabase=%v\nport=%v", c.Name, c.Database, c.Port)
}

type User struct {
	Hash  string `json:"hash"`
	Admin bool   `json:"admin"`
}

func passwordHash(pw string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(h)
}

func Parse(name string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if !path.IsAbs(name) {
		name = path.Join(wd, name)
	}
	C.Name = name

	if b, err := os.ReadFile(C.Name); err == nil {
		if err = json.Unmarshal(b, &C); err != nil {
			log.Fatal(err)
		}
	}
	if C.Users == nil {
		C.Users = map[string]*User{}
	}
	if len(C.Port) == 0 {
		C.Port = ":8000"
	}
	if len(C.Database) == 0 {
		C.Database = path.Join(wd, "audiolib.db")
	}
	if len(C.Users) == 0 {
		AddUser("admin", "", true)
	}
	if err := save(); err != nil {
		log.Fatal(err)
	}
}

func AddUser(name, password string, admin bool) error {
	M.Lock()
	defer M.Unlock()
	name = strings.ToLower(name)
	if _, ok := C.Users[name]; ok {
		return fmt.Errorf("user '%v' already exists", name)
	}
	C.Users[name] = &User{Hash: passwordHash(password), Admin: admin}
	return save()
}

func UpdatePassword(name, password string, admin bool) error {
	name = strings.ToLower(name)
	M.Lock()
	defer M.Unlock()
	u, ok := C.Users[name]
	if !ok {
		return fmt.Errorf("user '%v' doesn't exist", name)
	}
	u.Hash = passwordHash(password)
	u.Admin = admin
	return save()
}

func Login(name, password string) bool {
	name = strings.ToLower(name)
	M.Lock()
	defer M.Unlock()
	u, ok := C.Users[name]
	if !ok {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(password)) == nil
}

func save() error {
	M.Lock()
	defer M.Unlock()
	b, err := json.Marshal(C)
	if err != nil {
		return err
	}
	return os.WriteFile(C.Name, b, 0600)
}
