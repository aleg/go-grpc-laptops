package service

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func logError(err error, code codes.Code, msg string) error {
	if err != nil {
		log.Print(msg, err)
		return status.Errorf(code, "%s: %v", msg, err)
	}

	log.Print(msg)
	return status.Errorf(code, "%s", msg)
}

func contextError(ctx context.Context) error {
	ctxErr := ctx.Err()
	switch ctxErr {
	case context.Canceled:
		return logError(ctxErr, codes.Canceled, "Request cancelled")
	case context.DeadlineExceeded:
		return logError(ctxErr, codes.DeadlineExceeded, "Deadline is exceeded")
	default:
		return nil
	}
}
