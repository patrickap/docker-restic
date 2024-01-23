package lock

import (
	"errors"

	"github.com/gofrs/flock"
	"github.com/patrickap/docker-restic/m/v2/internal/env"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
)

var lock = flock.New(env.DOCKER_RESTIC_DIR + "/tmp/docker-restic.lock")

func Lock() error {
	locked, err := lock.TryLock()

	if !locked {
		return errors.New("failed to acquire lock")
	}

	if err != nil {
		return err
	}

	return nil
}

func Unlock() error {
	return lock.Unlock()
}

func RunWithLock(f func() error) error {
	err := lock.Lock()
	if err != nil {
		log.Error().Msg("Failed to acquire lock")
		return err
	}

	defer lock.Unlock()

	return f()
}
