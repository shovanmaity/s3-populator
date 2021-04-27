FROM docker.io/library/golang:1.14 AS builder
LABEL type=build-container
WORKDIR /go/src/github.com/shovanmaity/s3-populator
COPY . .
RUN make build

FROM gcr.io/distroless/static:latest
ENV PATH=/bin
COPY --from=builder /go/src/github.com/shovanmaity/s3-populator/s3-populator /bin/s3-populator
ENTRYPOINT ["s3-populator"]
