package lock

import (
	"errors"

	"github.com/gofrs/flock"
	"github.com/patrickap/docker-restic/m/v2/internal/env"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
)

var lock = flock.New(env.DOCKER_RESTIC_DIR + "/docker-restic.lock")

func Lock() error {
	locked, err := lock.TryLock()

	if !locked {
		return errors.New("could not acquire lock")
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
	lockErr := lock.Lock()
	if lockErr != nil {
		log.Error().Msg("Could not acquire lock")
		return lockErr
	}

	defer lock.Unlock()

	return f()
}
