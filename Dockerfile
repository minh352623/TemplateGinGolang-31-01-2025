FROM golang:1.23-alpine AS build

# Cài đặt các công cụ cần thiết
RUN apk add --no-cache git curl make gcc musl-dev


# Cài đặt migrate
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.0/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin/ && \
#     chmod + x /usr/local/bin/migrate

WORKDIR /build

# Copy và cài đặt các dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy toàn bộ mã nguồn
COPY . .
EXPOSE 8001
# Thực hiện migrate
# RUN migrate -database "cockroachdb://dev:3986c9107742f01daa7cf3c291@cockroachdb.railway.internal:26257/defaultdb?sslmode=disable" -path ./db/migrations up

# Chạy ứng dụng
CMD ["go", "run", "cmd/server/main.go"]
