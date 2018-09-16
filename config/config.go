package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// Config - Application configuration
type Config struct {
	workinkDir string

	Template string `json:"template"`
	Input    string `json:"input"`
	Output   string `json:"output"`

	BlogTitle       string `json:"blog_title"`
	BlogDescription string `json:"blog_description"`
	Author          string `json:"author"`
	Lang            string `json:"lang"`
	ServerPath      string `json:"server_path"`
	Origin          string `json:"origin"`
}

// New - creates new Config instance
func New(configPathRel string) *Config {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	c := &Config{}
	c.workinkDir = dir

	configPathAbs := path.Join(dir, configPathRel)

	file, err := ioutil.ReadFile(configPathAbs)

	if err != nil {
		panic(fmt.Errorf("Config not found in path: %s", configPathAbs))
	}

	e := json.Unmarshal(file, c)

	if e != nil {
		panic(fmt.Errorf("Cannot read config: %s", err))
	}

	return c
}

// GetAbsPath - returns absolute path based on working directory and passed relative path
func (c *Config) GetAbsPath(relPath string) string {
	return path.Join(c.workinkDir, relPath)
}
