FROM library/ubuntu:xenial

ENV GOVERSION 1.9.2
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN apt-get update && apt-get install -y wget curl git postgresql-client \
    software-properties-common && \
    add-apt-repository 'deb http://apt.postgresql.org/pub/repos/apt/ xenial-pgdg main' && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | \
    apt-key add - && apt-get update && apt-get install -y postgresql-client-10


RUN cd /usr/local && wget https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz && \
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/bin/

RUN mkdir /go /go/bin /go/pkg /go/src /go/src/rosella
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/rosella