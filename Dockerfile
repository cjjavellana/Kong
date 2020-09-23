FROM centos:centos7

RUN yum install -y git patch make zlib-devel libyaml libyaml-devel
RUN yum groups install -y "Development Tools"

RUN cd /tmp && \
git clone https://github.com/Kong/kong-build-tools.git && \
./kong-build-tools/openresty-build-tools/kong-ngx-build \
    --prefix /usr/local/kong \
    --work work \
    --openresty 1.15.8.3 \
    --openssl 1.1.1g \
    --kong-nginx-module master \
    --luarocks 3.3.1 \
    --pcre 8.44 \
    --jobs 6 \
    --force

ENV OPENSSL_DIR=/usr/local/kong/openssl
ENV PATH=/usr/local/kong/openresty/bin:$PATH
ENV PATH=/usr/local/kong/openresty/nginx/sbin:$PATH
ENV PATH=/usr/local/kong/openssl/bin:$PATH
ENV PATH=/usr/local/kong/luarocks/bin:$PATH

RUN cd /tmp && \
git clone https://github.com/Kong/kong.git --depth 1 && \
cd kong && \
git fetch --tags --all && \
git checkout tags/2.1.4 -b 2.1.4 && \
eval $(luarocks path --bin) && \
make install

ADD start-kong.sh /

