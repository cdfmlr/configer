package config

import (
	"bytes"
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
