package main

import (
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "build")
	cmd.Args = append(
		cmd.Args,

		// `--ldflags`, `-linkmode=external -extldflags="-static"`,
		`-tags`, `osusergo,netgo`,
		`-buildvcs=false`,
		"-o", "bin/app",
		"cmd/app/main.go",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
