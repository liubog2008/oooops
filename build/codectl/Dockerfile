FROM alpine:3.9
LABEL maintainer="LiuBo <iubog2008@gmail.com>"

RUN apk add --no-cache git
ADD _output/codectl /app/bin/codectl

ENTRYPOINT ["/app/bin/codectl"]
