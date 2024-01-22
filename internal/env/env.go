package env

import "os"

var DOCKER_RESTIC_DIR = os.Getenv("DOCKER_RESTIC_DIR")
