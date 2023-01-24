package ActivityApi

import (
	"fmt"
	"strings"
	"testing"
)

func TestLoggerInitialization(t *testing.T) {
	logger := CreateEmptyLogger()

	logger.AddCheckout()

	if logger.CheckoutCounter != 1 {
		t.Errorf("Add Checkout does not add properly")
	}

	logger.AddDecline()

	if logger.DeclineCounter != 1 {
		t.Errorf("AddDecline does not add properly")
	}

	logger.PushTaskStatusUpdate(1, "DAB", LogLevel)

	if logger.taskStatuses[0].Id != 1 && logger.taskStatuses[0].Status != "DAB" {
		t.Errorf("status was not pushed correctly")
	}

	out := logger.GetTaskStatusArray()

	if out[0] != fmt.Sprintf("[%d] %s", 1, "DAB") {
		t.Errorf("GetTaskStatusArray is missing a pushed status for arr: [%s]", strings.Join(out, ","))
	}
}
