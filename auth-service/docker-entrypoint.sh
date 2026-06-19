#!/bin/sh
set -e

/bin/auth-migrate
exec /bin/auth-service
