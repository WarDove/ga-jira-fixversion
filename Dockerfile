FROM golang:1.21.4-alpine3.18 as builder
WORKDIR /app
COPY go.* main.go ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o jira-release-automation .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/jira-release-automation .
RUN addgroup -g 1001 appuser && \
    adduser -D -H -s /sbin/nologin -u 1001 -G appuser appuser && \
    chown -R appuser:appuser /app
RUN chmod +x /app/jira-release-automation
USER appuser
RUN rm -rf /bin/*
CMD ["/app/jira-release-automation"]