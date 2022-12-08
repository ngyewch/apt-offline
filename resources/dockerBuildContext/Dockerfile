FROM debian:stretch
WORKDIR /workspace
COPY docker-entrypoint.sh .
RUN chmod a+x docker-entrypoint.sh
VOLUME /workspace/packages
ENTRYPOINT ["./docker-entrypoint.sh"]
