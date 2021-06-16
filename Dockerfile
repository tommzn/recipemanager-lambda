FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /go

COPY build_artifact_bin lambdahandler
RUN chmod 775 lambdahandler

ENTRYPOINT ["/go/lambdahandler"]
