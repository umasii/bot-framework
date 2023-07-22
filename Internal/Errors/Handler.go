package errors

import (
	"fmt"

	"github.com/getsentry/sentry-go"
)

func Handler(err error) error {
	if err != nil {
		sentry.CaptureException(err)
		// if Args.Args.IsDebug {
			fmt.Println(err)
		// }
	}
	return err
}
