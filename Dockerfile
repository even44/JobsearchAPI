# Image for building application
FROM golang:1.23-alpine AS build

# Set working directory to intended app location
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy rest of the files then build the app
COPY . .
RUN go build -v -o /usr/local/bin/app/ ./...

FROM scratch

COPY --from=build /usr/local/bin/app/ .

CMD ["./cmd"]