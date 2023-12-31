package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/yankeguo/activate-toolchain"
	"log"
	"os"
	"strings"

	_ "github.com/yankeguo/activate-toolchain/toolchains/jdk"
	_ "github.com/yankeguo/activate-toolchain/toolchains/kubectl"
	_ "github.com/yankeguo/activate-toolchain/toolchains/maven"
	_ "github.com/yankeguo/activate-toolchain/toolchains/node"
	_ "github.com/yankeguo/activate-toolchain/toolchains/ossutil"
	_ "github.com/yankeguo/activate-toolchain/toolchains/pnpm"
)

func main() {
	var err error
	defer func() {
		if err == nil {
			return
		}
		log.Println("exited with error:", err.Error())
		os.Exit(1)
	}()

	ctx := context.Background()

	args := os.Args[1:]

	// support toolchains.txt
	{
		buf, _ := os.ReadFile("toolchains.txt")
		for _, line := range bytes.Split(buf, []byte("\n")) {
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			args = append(args, string(line))
		}
	}

	var scripts []string

	for _, arg := range args {
		var spec activate_toolchain.Spec

		if spec, err = activate_toolchain.ParseSpec(arg); err != nil {
			return
		}

		toolchain, ok := activate_toolchain.FindToolchain(spec)

		if !ok {
			err = errors.New("no supported toolchain: " + arg)
			return
		}

		var script string
		if script, err = toolchain.Activate(ctx, spec); err != nil {
			return
		}

		head := fmt.Sprintf("# activate %s\n", arg)

		scripts = append(scripts, head+script)
	}

	if _, err = os.Stdout.WriteString(strings.Join(scripts, "\n\n")); err != nil {
		return
	}
	_ = os.Stdout.Sync()
}
