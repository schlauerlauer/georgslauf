FROM docker.io/library/golang:1.16.6 AS builder
WORKDIR /app/
COPY models ./models
COPY controllers ./controllers
COPY main.go go.mod go.sum .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api .

FROM scratch
WORKDIR /app/
COPY uploads ./uploads
COPY keymatch_model.conf .
COPY config_default.yaml config.yaml
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]