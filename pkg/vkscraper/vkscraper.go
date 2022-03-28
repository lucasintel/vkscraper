package vkscraper

import (
	"log"
	"net/http"
	"os"

	"github.com/kandayo/vkscraper/pkg/vk"
)

type Config struct {
	BaseDir    string
	FastUpdate bool
}

type Instance struct {
	Config    Config
	Vk        *vk.Client
	Log       *log.Logger
	transport *http.Client
}

func New(client *vk.Client, config Config) *Instance {
	log := log.New(os.Stdout, "", 0)
	return &Instance{
		Config: config,
		Vk:     client,
		Log:    log,
		transport: &http.Client{
			Timeout: timeout,
		},
	}
}
