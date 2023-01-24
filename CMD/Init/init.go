package Init

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cicadaaio/LVBot/CMD/TaskEngine"
	"github.com/cicadaaio/LVBot/Internal/Errors"
	"github.com/cicadaaio/LVBot/Internal/Flags"
	"github.com/getsentry/sentry-go"
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
	flag.Usage = Flags.Usage
	flag.Parse()

	if *Flags.FlagVersion {
		fmt.Fprintln(os.Stderr, Flags.Version)
		os.Exit(2)
	}
	if *Flags.FlagHelp {
		flag.Usage()
		os.Exit(2)
	}
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "does not take any operands")
		flag.Usage()
		os.Exit(2)
	}
	if *Flags.FlagRun != "" {

		TEngine := TaskEngine.TaskEngine{}
		TEngine.InitializeEngine()

		if *Flags.FlagRun == "all" {
			TEngine.StartAllTasks()
		} else {
			groups := strings.Split(*Flags.FlagRun, ",")
			for _, group := range groups {
				igroup, err := strconv.Atoi(group)
				if err != nil {
					Errors.Handler(err)
					log.Fatal(err)
				}
				TEngine.StartTasksInGroup(igroup)
			}
		}
	}
}
