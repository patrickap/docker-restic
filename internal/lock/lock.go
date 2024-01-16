package lock

import (
	"github.com/gofrs/flock"
	"github.com/patrickap/docker-restic/m/v2/internal/env"
)

var lock = flock.New(env.DOCKER_RESTIC_DIR + "/docker-restic.lock")

func Lock() (bool, error) {
	return lock.TryLock()
}

func Unlock() error {
	return lock.Unlock()
}
