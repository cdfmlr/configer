package config

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

type Habit struct {
	Name   string
	Rating int
}

type TestConfig struct {
	Name   string
	Age    int
	Habits []Habit
	Birth  time.Time
}

var testConfigInstance = TestConfig{
	Name: "test",
	Age:  18,
	Habits: []Habit{
		{
			Name:   "football",
			Rating: 5,
		}, {
			Name:   "basketball",
			Rating: 4,
		},
	},
	Birth: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
}

var testConfigYaml = `
name: test
age: 18
habits:
  - name: football
    rating: 5
  - name: basketball
    rating: 4
birth: 2000-01-01T00:00:00Z
`

var testConfigJson = `{
	"name": "test",
	"age": 18,
	"habits": [
		{
			"name": "football",
			"rating": 5
		},
		{
			"name": "basketball",
			"rating": 4
		}
	],
	"birth": "2000-01-01T00:00:00Z"
}`

var testConfigToml = `
name = "test"
age = 18
birth = "2000-01-01T00:00:00Z"

[[habits]]
name = "football"
rating = 5

[[habits]]
name = "basketball"
rating = 4
`

func TestConfiger_Read(t *testing.T) {
	for _, test := range []struct {
		encoding encoding
		data     string
	}{
		{YAML, testConfigYaml},
		{JSON, testConfigJson},
		{TOML, testConfigToml},
	} {
		var c TestConfig

		buf := []byte(test.data)

		err := NewConfiger(&c, test.encoding).Read(bytes.NewReader(buf))
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c, testConfigInstance) {
			t.Fatalf("expect %v, got %v", testConfigInstance, c)
		}

		t.Log(c)
	}
}

func TestConfiger_Write(t *testing.T) {
	for _, test := range []struct {
		encoding encoding
		data     string
	}{
		{YAML, testConfigYaml},
		{JSON, testConfigJson},
		{TOML, testConfigToml},
	} {
		var c TestConfig = testConfigInstance

		var buf bytes.Buffer

		err := NewConfiger(&c, test.encoding).Write(&buf)
		if err != nil {
			t.Fatal(err)
		}

		if reflect.DeepEqual(strings.TrimSpace(buf.String()), test.data) {
			t.Fatalf("\n--- expect:\n%v\n--- got:\n%v\n", test.data, buf.String())
		}

		t.Log(buf.String())
	}
}

// A simple example to of reading a TOML configuration.
func ExampleConfiger_Read() {
	var Config struct {
		Version     string
		HTTPService struct {
			ListenAddr string
			Timeout    time.Duration
		}
		Log struct {
			Level int
			File  string
		}
		BackendServices []struct {
			Addr   string
			Labels map[string]string
		}
		DB struct {
			URL  string
			Auth struct{ User, Password string }
		}
	}

	tomlData := `
version = "1.0.0"
HTTPService = { ListenAddr = ":8080", Timeout = "5s" }

[Log]
Level = 2
File = "app.log"

[[BackendServices]]
Addr = "https://backend1"
Labels = { env = "dev" }

[[BackendServices]]
Addr = "https://backend2"
Labels = { env = "prod", region = "us" }

[[BackendServices]]
Addr = "https://backend3"
[BackendServices.Labels]
env = "prod"
region = "eu"
disabled = "true"

[DB]
URL = "mysql://localhost:3306"

[DB.Auth]
User = "root"
Password = "pswd123"
`

	configer := NewConfiger(&Config, TOML)

	if err := configer.Read(strings.NewReader(tomlData)); err != nil {
		panic(err)
	}

	fmt.Println(
		Config.Version,
		Config.HTTPService.Timeout,
		len(Config.BackendServices),
		Config.BackendServices[len(Config.BackendServices)-1].Labels["region"],
		Config.DB.Auth.User)

	// Output:
	// 1.0.0 5s 3 eu root
}

// A simple example to of writing a JSON configuration.
func ExampleConfiger_Write() {
	var Config = struct {
		Version string
		Labels  map[string]string `json:"labels"`
		Comment string            `json:",omitempty"`
	}{
		Version: "1.0.0",
		Labels: map[string]string{
			"env":      "dev",
			"disabled": "true",
		},
	}

	configer := NewConfiger(&Config, JSON)

	var buf bytes.Buffer

	if err := configer.Write(&buf); err != nil {
		panic(err)
	}

	fmt.Println(buf.String())

	// Output:
	// {"Version":"1.0.0","labels":{"disabled":"true","env":"dev"}}
}
