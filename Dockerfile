FROM dockerfile/ubuntu
MAINTAINER soh335

RUN echo "Asia/Tokyo\n" > /etc/timezone && dpkg-reconfigure --frontend noninteractive tzdata

RUN apt-get update && apt-get install -y \
        ntp \
        curl \
        libav-tools \
        rtmpdump \
        swftools

RUN curl -LO https://github.com/soh335/radicast/releases/download/0.0.1/linux_amd64.zip
RUN unzip linux_amd64.zip
RUN mv radicast /usr/local/bin/

ENTRYPOINT ["radicast"]
CMD ["--help"]
