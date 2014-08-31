package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/satoshun/gitcd/git"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitcd"
	app.Version = "0.0.1"
	app.Usage = "For git commit driven development"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		{
			Name:      "init",
			ShortName: "i",
			Usage:     "initialize commit driven list",
			Action: func(c *cli.Context) {
				wd, _ := os.Getwd()
				if !git.Exists(wd) {
					log.Println("no git repository", wd)
					return
				}

				config := &Config{WorkingDirectory: wd, Current: 0}
				fmt.Println("input commit message / if empty then done")
				for {
					var message string
					fmt.Scanf("%s", &message)
					message = strings.TrimSpace(message)
					if message == "" {
						break
					}

					commit := Commit{Message: message}
					config.Commits = append(config.Commits, commit)
				}

				config.Save()
				fmt.Println(`initialize ok`)
			},
		},
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "list commit driven list",
			Action: func(c *cli.Context) {
				wd, _ := os.Getwd()
				if !git.Exists(wd) {
					log.Println("no git repository", wd)
					return
				}

				config, err := Read(wd)
				if err != nil {
					log.Println("yet initialize", err)
					return
				}
				for i, commit := range config.Commits {
					prefix := strconv.Itoa(i + 1)
					message := commit.Message
					if commit.Done {
						prefix = "done " + prefix
					}
					if i == config.Current {
						prefix = "current " + prefix
					}
					fmt.Printf("%10s: %s\n", prefix, message)
				}
			},
		},
		{
			Name:      "next",
			ShortName: "n",
			Usage:     "next commit driven list",
			Action: func(c *cli.Context) {
				wd, _ := os.Getwd()
				if !git.Exists(wd) {
					log.Println("no git repository", wd)
					return
				}

				config, err := Read(wd)
				if err != nil {
					log.Println("yet initialize", err)
					return
				}

				commit := config.NextCommit()
				if commit == nil {
					log.Println("already done")
					return
				}
				cmd := git.CommitCmd(commit.Message)
				err = cmd.Run()
				if err != nil {
					log.Println("fail commit", err)
					return
				}
				commit.Done = true
				config.Current += 1
				config.Save()

				commit = config.NextCommit()
				fmt.Printf("next: %s", commit.Message)
			},
		},
		{
			Name:      "prev",
			ShortName: "p",
			Usage:     "prev commit driven list",
			Action: func(c *cli.Context) {
				wd, _ := os.Getwd()
				if !git.Exists(wd) {
					log.Println("no git repository", wd)
					return
				}

				config, err := Read(wd)
				if err != nil {
					log.Println("yet initialize", err)
					return
				}

				if !config.IsPrev() {
					log.Println("yet commit")
					return
				}

				commit := config.PrevCommit()

				cmd := git.ResetCmd(false)
				err = cmd.Run()
				if err != nil {
					log.Println("fail prev command", err)
					return
				}

				commit.Done = false
				config.Current -= 1
				config.Save()

				commit = config.NextCommit()
				fmt.Printf("next: %s", commit.Message)
			},
		},
		{
			Name:      "fix",
			ShortName: "f",
			Usage:     "fix commit driven list",
			Action: func(c *cli.Context) {
				wd, _ := os.Getwd()
				if !git.Exists(wd) {
					log.Println("no git repository", wd)
					return
				}
				config, err := Read(wd)
				if err != nil {
					log.Println("yet initialize", err)
					return
				}

				if !config.IsPrev() {
					log.Println("yet commit")
					return
				}

				commit := config.PrevCommit()
				cmd := git.AmendCmd(commit.Message)
				err = cmd.Run()
				if err != nil {
					log.Println("fail amend command", err)
					return
				}
			},
		},
		{
			Name:      "abort",
			ShortName: "a",
			Usage:     "abort commit driven list",
			Action: func(c *cli.Context) {
				wd, _ := os.Getwd()
				if !git.Exists(wd) {
					log.Println("no git repository", wd)
					return
				}
				if !Exists(ConfigPath(wd)) {
					log.Println("yet initalize")
					return
				}

				err := os.Remove(ConfigPath(wd))
				if err != nil {
					log.Println("fail remove file")
					return
				}
				fmt.Println("success: remove file", ConfigPath(wd))
			},
		},
	}

	app.Run(os.Args)
}
