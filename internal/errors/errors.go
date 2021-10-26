package errors

import (
	"fmt"

	"go.uber.org/zap"
)

func LogInitializationError(err error, step string, log *zap.Logger) {
	log.Fatal(
		fmt.Sprintf("an error occurred during %s initialization step: %s", step, err),
		zap.String("step", step),
	)
}
