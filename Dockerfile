FROM golang:1.22-alpine as builder
RUN apk --update add build-base

WORKDIR /src/app
ADD go.mod .
RUN go mod download

ADD . .

# Building TailwindCSS with tailo
RUN go run github.com/paganotoni/tailo/cmd/build@a4899cd
RUN go build -tags osusergo,netgo -buildvcs=false -o bin/app ./cmd/app

# Building the migrate command to run migrations
RUN go build -tags osusergo,netgo -buildvcs=false -o bin/migrate ./cmd/migrate

FROM alpine
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /bin/

# Copying binaries
COPY --from=builder /src/app/bin/app .
COPY --from=builder /src/app/bin/migrate .

# Running migrations and then the app
CMD migrate && app
