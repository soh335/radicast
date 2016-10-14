FROM ubuntu:trusty
MAINTAINER soh335

RUN echo "Asia/Tokyo\n" > /etc/timezone && dpkg-reconfigure --frontend noninteractive tzdata

RUN apt-get update && apt-get install -y \
        ntp \
        curl \
        libav-tools \
        rtmpdump \
        swftools \
        git

# http://blog.gopheracademy.com/advent-2014/easy-deployment/
RUN mkdir /goroot && curl https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1

ENV GOROOT /goroot
ENV GOPATH /gopath
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

RUN go get -v github.com/soh335/radicast

ENTRYPOINT ["radicast"]
CMD ["--help"]
