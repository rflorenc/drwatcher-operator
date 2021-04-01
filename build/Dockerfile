FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR_DEST_DIR=/usr/local/bin/drwatcher-operator \
    USER_UID=1337

# Assumes `make` previously ran sucessfully from the root dir
COPY drwatcher-operator ${OPERATOR_DEST_DIR}

USER ${USER_UID}

CMD ["/usr/local/bin/drwatcher-operator"]
