FROM ubuntu:latest

RUN apt-get update
RUN apt-get -y install wget

RUN wget 'http://nginx.org/download/nginx-1.21.6.tar.gz'
RUN tar -xzvf nginx-1.21.6.tar.gz

RUN apt-get install -y git
RUN apt-get install -y gcc
RUN apt-get install -y make
RUN apt-get install -y libpcre3 libpcre3-dev zlib1g zlib1g-dev libssl-dev

WORKDIR nginx-1.21.6
RUN apt-get install -y mercurial
RUN hg clone http://hg.nginx.org/njs
RUN git clone https://github.com/openresty/redis2-nginx-module
RUN ./configure --prefix=/etc/nginx --modules-path=/etc/nginx/modules --sbin-path=/usr/sbin/nginx --conf-path=/etc/nginx/nginx.conf --with-http_sub_module --add-module=redis2-nginx-module --add-dynamic-module=njs/nginx
RUN make -j2
RUN make install

WORKDIR ../

COPY ./ /etc/nginx

RUN mkdir /var/log/nginx
RUN touch /var/log/nginx/error.log

EXPOSE 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]