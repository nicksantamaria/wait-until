package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Exit codes adapted from http://tldp.org/LDP/abs/html/exitcodes.html
const (
	// ExitCodeSuccess indicates script success.
	ExitCodeSuccess int = 0
	// ExitCodeCommandNotFound indicates command not found.
	ExitCodeCommandNotFound int = 127
	// ExitCodeTimeoutExceeded indicates timeout exceeded.
	ExitCodeTimeoutExceeded int = 1
	// ExitCodeRetriesExceeded indicates retry limit exceeded.
	ExitCodeRetriesExceeded int = 1
)

var (
	verbose  = kingpin.Flag("verbose", "Enabled verbose output.").Short('v').Bool()
	timeout  = kingpin.Flag("timeout", "Timeout before aborting pipeline. Omit for no limit.").Short('t').Duration()
	retries  = kingpin.Flag("retries", "Number of attempts before aborting pipeline. -1 for no limit.").Short('r').Default("-1").Int()
	sleep    = kingpin.Flag("sleep", "Sleep time between each execution.").Short('s').Default("1s").Duration()
	exitCode = kingpin.Flag("exit-code", "Desired exit code before allowing pipeline to proceed.").Short('e').Default("0").Int()
	command  = kingpin.Arg("command", "Command to repeatedly execute until exit code met, timeout exceeded, or retry limit exceeded.").Required().Strings()
)

func main() {
	kingpin.Parse()
	args := *command

	printf("waiting for command to return exit code %d: '%s'\n", *exitCode, strings.Join(*command, " "))

	var endTime *time.Time
	if timeout != nil {
		t := time.Now().Add(*timeout)
		endTime = &t
	}

	i := 0
	for {
		if timeout != nil && time.Now().After(*endTime) {
			fmt.Println("timeout exceeded")
			os.Exit(ExitCodeTimeoutExceeded)
		}
		if *retries > -1 && i >= *retries {
			fmt.Println("retry limit reached")
			os.Exit(ExitCodeRetriesExceeded)
		}

		cmd := exec.Command(args[0], args[1:]...)
		_ = cmd.Run()
		if cmd.ProcessState.ExitCode() == *exitCode {
			os.Exit(ExitCodeSuccess)
		}

		time.Sleep(*sleep)
		i++
	}
}

func printf(format string, args ...interface{}) {
	if *verbose {
		fmt.Printf(format, args...)
	}
}
