FROM golang:1.19 AS builder
RUN mkdir /app
WORKDIR /app
ADD . /app
COPY . .
COPY go.mod .
RUN go mod download
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -o godocker -a -ldflags '-linkmode external -extldflags "-static"' .
RUN chmod +x /app/godocker

FROM scratch
COPY --from=builder /app /app
COPY --from=builder /app/static /app/static
COPY --from=builder /app/edMember.html /app/edMember.html
COPY --from=builder /app/sqlite3msr.db /app/sqlite3msr.db
EXPOSE 8084

CMD ["/app/godocker"]
