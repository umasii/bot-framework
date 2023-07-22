package init

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	taskengine "github.com/umasii/bot-framework/cmd/taskengine"
	errors "github.com/umasii/bot-framework/internal/errors"
	flags "github.com/umasii/bot-framework/internal/flags"
)

func InitSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "",
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
}

func InitCicada() {
	rand.Seed(time.Now().Unix())

	InitSentry()
	defer sentry.Recover()

	log.SetPrefix("framework: ")
	log.SetFlags(0)
	flag.Usage = flags.Usage
	flag.Parse()

	if *flags.FlagVersion {
		fmt.Fprintln(os.Stderr, flags.Version)
		os.Exit(2)
	}
	if *flags.FlagHelp {
		flag.Usage()
		os.Exit(2)
	}
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "does not take any operands")
		flag.Usage()
		os.Exit(2)
	}
	if *flags.FlagRun != "" {

		TEngine := taskengine.TaskEngine{}
		TEngine.InitializeEngine()

		if *flags.FlagRun == "all" {
			TEngine.StartAllTasks()
		} else {
			groups := strings.Split(*flags.FlagRun, ",")
			for _, group := range groups {
				igroup, err := strconv.Atoi(group)
				if err != nil {
					errors.Handler(err)
					log.Fatal(err)
				}
				TEngine.StartTasksInGroup(igroup)
			}
		}
	}
}
