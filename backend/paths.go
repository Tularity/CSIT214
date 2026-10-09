package main

import (
	"log"
	"os"
	"path/filepath"
)

// releaseBuild is set to "true" by the release build script. A release binary
// is usually started by double-clicking it, and then the working directory is
// wherever the desktop chose (the home folder on macOS), so its files are
// placed relative to the executable instead.
var releaseBuild = "false"

type paths struct {
	database string
	seed     string
}

func resolvePaths() paths {
	resolved := paths{
		database: os.Getenv("CSIT214_DB_PATH"),
		seed:     os.Getenv("CSIT214_SEED_PATH"),
	}
	if releaseBuild != "true" {
		return resolved
	}

	executable, err := os.Executable()
	if err != nil {
		log.Printf("cannot locate the executable, using the working directory: %v", err)
		return resolved
	}
	if dir, err := filepath.EvalSymlinks(filepath.Dir(executable)); err == nil {
		executable = filepath.Join(dir, filepath.Base(executable))
	}
	binDir := filepath.Dir(executable)

	if resolved.database == "" {
		resolved.database = filepath.Join(binDir, "data.db")
	}
	if resolved.seed == "" {
		bundled := filepath.Join(binDir, "..", "seed", "seed.sql")
		if _, err := os.Stat(bundled); err == nil {
			resolved.seed = bundled
		}
	}
	return resolved
}
