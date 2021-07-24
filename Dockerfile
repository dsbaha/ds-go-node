FROM scratch
COPY ds-go-node /
ENTRYPOINT [ "/ds-go-node" ]
