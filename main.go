package main

import (
	"github.com/jrogala/jellyfin-cli/cmd"
	"github.com/jrogala/jellyfin-cli/config"
)

func main() {
	config.Init()
	cmd.Execute()
}
