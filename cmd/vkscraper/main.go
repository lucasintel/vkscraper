package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kandayo/vkscraper/pkg/vk"
	"github.com/kandayo/vkscraper/pkg/vkscraper"
	"github.com/urfave/cli/v2"
)

func main() {
	var (
		login      string
		password   string
		fastUpdate bool
		batchFile  string
		users      []string
	)

	app := cli.App{
		Name:     "vkscraper",
		Version:  "v0.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{Name: "Lucas M.D.", Email: "lucas@absolab.xyz"},
		},
		Usage:     "Download posts, photos, and videos along with their captions and other metadata from ВКонтакте.",
		UsageText: "vkscraper [--fast-update] [--batch-file] profile",
	}

	app.UseShortOptionHandling = true

	// TODO: remove this library.
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "login",
			Aliases:     []string{"l"},
			Usage:       "VK username",
			Destination: &login,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "password",
			Aliases:     []string{"p"},
			Usage:       "VK password",
			Destination: &password,
			Required:    false,
		},
		&cli.BoolFlag{
			Name:        "fast-update",
			Aliases:     []string{"f"},
			Usage:       "For each target, stop when encountering the first already-downloaded resource. This option is recommended when you use vkscraper to update your personal archive.",
			Destination: &fastUpdate,
			Required:    false,
			Value:       false,
		},
		&cli.StringFlag{
			Name:        "batch-file",
			Aliases:     []string{"a"},
			Usage:       "Read users and communities from file. Empty lines or lines starting with '#' are considered as comments and ignored.",
			Destination: &batchFile,
			Required:    false,
			TakesFile:   true,
		},
	}

	app.Action = func(c *cli.Context) error {
		if batchFile != "" {
			file, err := os.Open(batchFile)
			if err != nil {
				return err
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				userCandidate := scanner.Text()
				if strings.HasPrefix(userCandidate, "#") {
					continue
				} else {
					inlineCommentRegexp := regexp.MustCompile("#.*$")
					userCandidate = inlineCommentRegexp.ReplaceAllString(userCandidate, "")
					userCandidate = strings.TrimSpace(userCandidate)
					users = append(users, userCandidate)
				}
			}
		}

		requestedUsersFromArgs := c.Args().Slice()
		for _, user := range requestedUsersFromArgs {
			users = append(users, user)
		}

		client := vk.NewClient()

		// TODO: Refactor
		sessionFilename := login + ".vksession"
		sessionFileBuffer, _ := ioutil.ReadFile(sessionFilename)
		sessionFileData := strings.TrimSpace(string(sessionFileBuffer))
		if sessionFileData == "" {
			if password == "" {
				fmt.Fprintf(os.Stderr, "Session file not found. Please provide a password\n")
				os.Exit(1)
			}
			err := client.Login(login, password)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not login: %v\n", err)
				os.Exit(1)
			}
			sessionFile, err := os.Create(sessionFilename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could create session file: %v\n", err)
				os.Exit(1)
			}
			sessionFile.WriteString(client.AccessToken)
			sessionFile.Close()
			fmt.Printf("Saved session to %s\n", sessionFilename)
		} else {
			fmt.Printf("Loaded session from %s\n", sessionFilename)
			client.SetAccessToken(sessionFileData)
		}
		fmt.Println("Successfully logged in")

		currentDir := "."
		// TODO: Refactor
		config := vkscraper.Config{
			FastUpdate: fastUpdate,
			BaseDir:    currentDir,
		}
		instance := vkscraper.New(client, config)
		instance.DownloadProfiles(users)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nERROR: %s\n", err)
		os.Exit(1)
	}
}
