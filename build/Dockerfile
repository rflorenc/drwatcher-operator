FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR_BIN=/usr/local/bin/drwatcher-operator \
    USER=1337:1337

# Assumes `make` previously ran sucessfully from the root dir
COPY drwatcher-operator ${OPERATOR_BIN}

USER ${USER}

CMD ["${OPERATOR_BIN}"]
