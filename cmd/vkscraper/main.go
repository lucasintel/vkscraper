package main

import (
	"bufio"
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
		accessToken string
		fastUpdate  bool
		batchFile   string
		users       []string
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

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "access-token",
			Aliases:     []string{"t"},
			Usage:       "VK API Token. Can also be passed as an environment variable. See: https://vk.com/dev/access_token.",
			Destination: &accessToken,
			EnvVars:     []string{"VKSCRAPER_ACCESS_TOKEN"},
			Required:    true,
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

		client := vk.NewClient(accessToken)
		currentDir := "."
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
		panic(err)
	}
}
