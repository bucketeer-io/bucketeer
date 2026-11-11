# Atlas configuration for the Docker Compose migration service.
#
# The complete database URL (user, password, host, port and database name) is
# read from the migration_db_url Docker secret, so no part of it is written in
# compose.yml or in this file. `make docker-compose-setup` creates
# docker-compose/secrets/migration_db_url.txt for the bundled MySQL when the
# file is missing; edit that file to change any part of the URL.
# Format: mysql://user:password@host:port/db_name, with special characters in
# the user name and password percent-encoded.

env "docker-compose" {
  # A URL has no leading or trailing whitespace of its own (a space in the
  # password is percent-encoded), so trimming only removes the newline an
  # editor leaves at the end of the file.
  url = trimspace(file("/run/secrets/migration_db_url"))
}
