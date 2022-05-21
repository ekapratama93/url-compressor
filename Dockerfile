# argument for Go version
ARG GO_VERSION=1.17.10

# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./

# Build the executable
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app

# STAGE 2: build the container to run
FROM gcr.io/distroless/static AS final

USER nonroot:nonroot

# copy compiled app
COPY --from=build --chown=nonroot:nonroot /app /app

# run binary; use vector form
CMD ["/app"]
