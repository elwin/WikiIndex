FROM golang:1.23 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY app ./app
COPY assets ./assets
COPY database ./database
COPY mock ./mock
RUN CGO_ENABLED=0 go build


FROM golang:alpine
WORKDIR /app
COPY simplewiki-20211001-pages-articles-multistream.xml.bz2 wiki.xml.bz2
COPY --from=builder /app/WikiIndex ./WikiIndex
COPY view ./view
ENTRYPOINT ["/app/WikiIndex", "-i", "wiki.xml.bz2"]