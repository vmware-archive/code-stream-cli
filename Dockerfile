FROM scratch
COPY cs-cli /
ENTRYPOINT [ "/cs-cli" ]