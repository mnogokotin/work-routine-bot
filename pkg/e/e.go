package e

import "fmt"

func Wrap(err error, msg string, op string) error {
	if op != "" {
		return fmt.Errorf("%s: %s: %w", op, msg, err)
	} else {
		return fmt.Errorf("%s: %w", msg, err)
	}
}

func WrapIfErr(err error, msg string, op string) error {
	if err == nil {
		return nil
	}

	return Wrap(err, msg, op)
}
