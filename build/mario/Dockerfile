FROM alpine/git:v2.24.3

RUN mkdir /app
WORKDIR /app

COPY _output/mario mario
RUN chmod +x mario

CMD ["/app/mario"]

