# Atlas configuration for the Docker Compose migration service.
#
# The database password is read from the Docker secret that
# `make docker-compose-setup` creates (docker-compose/secrets/mysql_password.txt),
# the same file the mysql service reads, so the database URL is not written in
# compose.yml.

locals {
  # chomp() removes only the trailing newlines, like the mysql image does when
  # it reads the same file, so spaces that belong to the password are kept.
  mysql_password = chomp(file("/run/secrets/mysql_password"))
  # urlescape() encodes a space as "+", which is a literal plus sign in the
  # user-info part of a URL. A real plus sign is already encoded as %2B, so
  # every "+" left after urlescape() stands for a space.
  mysql_password_escaped = replace(urlescape(local.mysql_password), "+", "%20")
}

env "docker-compose" {
  url = "mysql://bucketeer:${local.mysql_password_escaped}@mysql:3306/bucketeer"
}
