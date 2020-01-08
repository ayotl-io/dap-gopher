FROM golang:1.13.5 as builder

WORKDIR /app
COPY . .

#RUN go get "github.com/dgrijalva/jwt-go"
RUN go get "github.com/fatih/color"
RUN CGO_ENABLED=0 GOOS=linux go build -v -o dap-gopher

FROM alpine
run apk add --no-cache bash curl openssl

COPY --from=builder /app/dap-gopher /dap-gopher
COPY --from=builder /app/retrieve.sh /retrieve.sh

#EXPOSE 8181
