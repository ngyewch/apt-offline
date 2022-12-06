FROM debian:stretch
WORKDIR /workspace
RUN dpkg --add-architecture armhf
RUN apt update -o APT::Architecture="armhf" -o APT::Architectures="armhf"
COPY docker-entrypoint.sh .
RUN chmod a+x docker-entrypoint.sh
VOLUME /workspace/packages
ENTRYPOINT ["./docker-entrypoint.sh"]
