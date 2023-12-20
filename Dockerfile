FROM alpine

WORKDIR "/opt"
ADD ./build/saibtc.config saibtc.config
ADD ./build/main main
RUN chmod +x /opt/main

ENTRYPOINT ["/opt/main"]