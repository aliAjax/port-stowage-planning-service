FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/planner ./cmd/planner
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/planner /planner
EXPOSE 8080
ENTRYPOINT ["/planner"]
