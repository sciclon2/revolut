FROM golang:1.19 AS builder

# Peprare container to compile the GO code
WORKDIR /app
COPY . .
RUN go mod init github.com/scilcon2/revolut_homework || true
RUN go mod tidy
RUN apt-get update && apt-get install -y gcc g++ make
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./api.go

# Prepare container to run the compiled code
FROM amazonlinux:2023
RUN yum install -y ca-certificates sqlite
WORKDIR /root/
COPY --from=builder /app/main .
RUN chmod +x ./main
EXPOSE 8080

# Run the app
CMD ["./main"]
