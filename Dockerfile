# Sử dụng image golang để build ứng dụng
FROM golang:1.22-alpine

WORKDIR /app

# Copy go.mod và go.sum và tải các dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy mã nguồn vào container
COPY . .

# Build ứng dụng
RUN go build -o main .

# Chạy ứng dụng
CMD ["./main"]