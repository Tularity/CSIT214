package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

// openBrowser saves someone who double-clicked the binary from having to know
// which address to type. Failing to open one is harmless: the address is
// printed as well.
func openBrowser(url string) {
	if releaseBuild != "true" || os.Getenv("CSIT214_OPEN_BROWSER") == "0" {
		return
	}

	var command *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		command = exec.Command("open", url)
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		command = exec.Command("xdg-open", url)
	}
	if err := command.Start(); err != nil {
		log.Printf("could not open a browser (%v); open %s yourself", err, url)
		return
	}
	go command.Wait()
}
