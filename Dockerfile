FROM gcr.io/distroless/static:latest
ADD s3-populator /
ENTRYPOINT ["/s3-populator"]
