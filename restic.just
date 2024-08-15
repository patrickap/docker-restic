restic_flags := "--repo /srv/restic/data --password-file /run/secrets/restic-password"

backup:
  restic backup /source {{restic_flags}}

prune:
  restic forget {{restic_flags}} \
    --keep-daily 7 \
    --keep-weekly 4 \
    --keep-monthly 12 \
    --keep-yearly 2 \
    --group-by paths \
    --prune

check:
  restic check {{restic_flags}} \
    --read-data
