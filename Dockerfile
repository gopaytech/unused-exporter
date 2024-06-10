#################
# Builder image
#################
FROM --platform=$BUILDPLATFORM golang:1.22-bullseye AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build main.go

#################
# Final image
#################
FROM --platform=$BUILDPLATFORM gcr.io/distroless/base

COPY --from=builder /app/main /

CMD ["/main"]
