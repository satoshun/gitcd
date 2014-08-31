package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

const configName = ".gitcd"

type Config struct {
	WorkingDirectory string   `json:"working_directory"`
	Current          int      `json:"current"`
	Commits          []Commit `json:"messages"`
}

type Commit struct {
	Message string `json:"message"`
	Done    bool   `json:"done"`
}

func (c *Config) Save() (err error) {
	b, _ := json.MarshalIndent(c, "", "  ")
	err = ioutil.WriteFile(ConfigPath(c.WorkingDirectory), b, 0644)

	return
}

func (c *Config) NextCommit() *Commit {
	if c.IsNext() {
		return &c.Commits[c.Current]
	}

	return nil
}

func (c *Config) PrevCommit() *Commit {
	return &c.Commits[c.Current-1]
}

func (c *Config) IsNext() bool {
	return len(c.Commits) > c.Current
}

func (c *Config) IsPrev() bool {
	return c.Current > 0
}

func Read(wd string) (config *Config, err error) {
	b, err := ioutil.ReadFile(path.Join(wd, ".git", configName))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &config)
	return
}

func ConfigPath(wd string) string {
	return path.Join(wd, ".git", configName)
}

func Exists(p string) bool {
	if _, err := os.Stat(p); err == nil {
		return true
	}

	return false
}
