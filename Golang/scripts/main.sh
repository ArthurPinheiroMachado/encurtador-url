
# !/usr/bin/sh

export DB_TYPE=postgres
export DB_NAME=encurtador
export DB_PASS=postgres
export DB_PORT=5432
export DB_HOST=0.0.0.0
export DB_USER=postgres
export HTTP_PORT=6060
export HTTP_BASE="/api/"
export TIMEOUT_TIME=3 #Segundos

go run cmd/server/main.go
 # go build -ldflags "-s -w" cmd/server/main.go

 # ./main
