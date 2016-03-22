FROM golang:1.5
MAINTAINER Eric Rasche <esr@tamu.edu>
EXPOSE 80

ENV GIT_REV 542e61c8762779be48690d9e6a42523143ab9277

RUN wget https://github.com/erasche/barcode-server/archive/${GIT_REV}.tar.gz && \
    tar xvfz ${GIT_REV}.tar.gz && \
    mv barcode-server-${GIT_REV}/ /app/ && \
    go get -v github.com/codegangsta/cli && \
    go get -v github.com/boombuler/barcode && \
    go get -v github.com/gorilla/mux && \
    go get -v github.com/gorilla/handlers

CMD ["go", "run", "/app/main.go", "-l", "0.0.0.0:80"]

