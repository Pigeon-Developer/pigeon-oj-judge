FROM ubuntu:22.04



RUN apt-get -y update && apt install -y ca-certificates

RUN apt-get -y update  && \
    apt-get -y upgrade && \
    DEBIAN_FRONTEND=noninteractive \
    apt-get -y install --no-install-recommends \
    nginx \
    libmysqlclient-dev \
    libmysql++-dev \
    php-common \
    php-fpm \
    php-mysql \
    php-gd \
    php-zip \
    php-mbstring \
    php-xml \
    php-yaml \
    make \
    flex \
    gcc \
    g++ \
    git


RUN mkdir /app && cd /app && git clone -b 24.12.25 https://github.com/zhblue/hustoj.git
RUN mv /app/hustoj/trunk /trunk
COPY hustoj.setup.sh /app/hustoj.setup.sh
RUN bash /app/hustoj.setup.sh
COPY hustoj.entrypoint.sh /hustoj.entrypoint.sh

EXPOSE 80

ENTRYPOINT [ "/bin/bash", "/hustoj.entrypoint.sh" ]
