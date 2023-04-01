package flags

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGetBotKey(t *testing.T) {
	var tests = []struct {
		testname                 string
		botKeyFlag               string
		hasBotFile               bool
		botKeyFileContents       string
		botKeyEnvVar             string
		hasBotEnvFile            bool
		botKeyFileEnvVarContents string
		expected                 string
	}{
		{
			testname:                 "Bot key flag only",
			botKeyFlag:               "foo",
			hasBotFile:               false,
			botKeyFileContents:       "",
			botKeyEnvVar:             "",
			hasBotEnvFile:            false,
			botKeyFileEnvVarContents: "",
			expected:                 "foo",
		},
		{
			testname:                 "Bot key file only",
			botKeyFlag:               "",
			hasBotFile:               true,
			botKeyFileContents:       "foo",
			botKeyEnvVar:             "",
			hasBotEnvFile:            false,
			botKeyFileEnvVarContents: "",
			expected:                 "foo",
		},
		{
			testname:                 "Bot key env var only",
			botKeyFlag:               "",
			hasBotFile:               false,
			botKeyFileContents:       "",
			botKeyEnvVar:             "foo",
			hasBotEnvFile:            false,
			botKeyFileEnvVarContents: "",
			expected:                 "foo",
		},
		{
			testname:                 "Bot key env file only",
			botKeyFlag:               "",
			hasBotFile:               false,
			botKeyFileContents:       "",
			botKeyEnvVar:             "",
			hasBotEnvFile:            true,
			botKeyFileEnvVarContents: "foo",
			expected:                 "foo",
		},
		{
			testname:                 "Both bot token flag and env var, match",
			botKeyFlag:               "foo",
			hasBotFile:               false,
			botKeyFileContents:       "",
			botKeyEnvVar:             "foo",
			hasBotEnvFile:            false,
			botKeyFileEnvVarContents: "",
			expected:                 "foo",
		},
	}

	for _, testcase := range tests {
		t.Logf("Running test: %s", testcase.testname)

		os.Setenv("BOT_TOKEN", "")
		os.Setenv("BOT_TOKEN_FILE", "")

		botFile := ""
		if testcase.hasBotFile {
			f, err := ioutil.TempFile("/tmp", "test")
			if err != nil {
				t.Errorf("Error creating temp file: %v", err)
			}
			defer os.Remove(f.Name())
			f.Write([]byte(testcase.botKeyFileContents))
			botFile = f.Name()
		}

		os.Setenv("BOT_TOKEN", testcase.botKeyEnvVar)
		if testcase.hasBotEnvFile {
			f, err := ioutil.TempFile("/tmp", "test")
			if err != nil {
				t.Errorf("Error creating temp file: %v", err)
			}
			defer os.Remove(f.Name())
			f.Write([]byte(testcase.botKeyFileEnvVarContents))
			os.Setenv("BOT_TOKEN_FILE", f.Name())
		}

		actual, err := GetBotKey(testcase.botKeyFlag, botFile)
		if err != nil {
			t.Errorf("Got unexpected error %v", err)
		}
		if actual != testcase.expected {
			t.Errorf("Expected %v, got %v", testcase.expected, actual)
		}
	}
}

func TestGetBotKeyErrors(t *testing.T) {
	var tests = []struct {
		testname                 string
		botKeyFlag               string
		hasBotFile               bool
		botKeyFileContents       string
		botKeyEnvVar             string
		hasBotEnvFile            bool
		botKeyFileEnvVarContents string
	}{
		{
			testname:                 "Bot key flag and env var both specified but differ",
			botKeyFlag:               "foo",
			hasBotFile:               false,
			botKeyFileContents:       "",
			botKeyEnvVar:             "bar",
			hasBotEnvFile:            false,
			botKeyFileEnvVarContents: "",
		},
		{
			testname:                 "Bot token file flag and env var both specified but differ",
			botKeyFlag:               "",
			hasBotFile:               true,
			botKeyFileContents:       "foo",
			botKeyEnvVar:             "",
			hasBotEnvFile:            true,
			botKeyFileEnvVarContents: "foo",
		},
		{
			testname:                 "No bot token or file",
			botKeyFlag:               "",
			hasBotFile:               false,
			botKeyFileContents:       "",
			botKeyEnvVar:             "",
			hasBotEnvFile:            false,
			botKeyFileEnvVarContents: "",
		},
	}

	for _, testcase := range tests {
		t.Logf("Running test: %s", testcase.testname)
		os.Setenv("BOT_TOKEN", "")
		os.Setenv("BOT_TOKEN_FILE", "")

		botFile := ""
		if testcase.hasBotFile {
			f, err := ioutil.TempFile("/tmp", "test")
			if err != nil {
				t.Errorf("Error creating temp file: %v", err)
			}
			defer os.Remove(f.Name())
			f.Write([]byte(testcase.botKeyFileContents))
			botFile = f.Name()
		}

		os.Setenv("BOT_TOKEN", testcase.botKeyEnvVar)
		if testcase.hasBotEnvFile {
			f, err := ioutil.TempFile("/tmp", "test")
			if err != nil {
				t.Errorf("Error creating temp file: %v", err)
			}
			defer os.Remove(f.Name())
			f.Write([]byte(testcase.botKeyFileEnvVarContents))
			os.Setenv("BOT_TOKEN_FILE", f.Name())
		}

		actual, err := GetBotKey(testcase.botKeyFlag, botFile)
		if err == nil {
			t.Errorf("GetBotKey succeeded unexpectedly, returned %s", actual)
		}
	}
}
