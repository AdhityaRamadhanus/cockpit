FROM golang:1.12.9

WORKDIR /go/src/github.com/AdhityaRamadhanus/cockpit
COPY . .

RUN make build-cron

EXPOSE 8000
ENTRYPOINT ["./cron"]