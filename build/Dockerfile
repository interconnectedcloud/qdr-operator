FROM alpine:3.8

ENV OPERATOR=/usr/local/bin/qdrouterd-operator \
    USER_UID=1001 \
    USER_NAME=qdrouterd-operator \
    QDROUTERD_IMAGE=quay.io/interconnectedcloud/qdrouterd:1.6.0

# install operator binary
COPY build/_output/bin/qdrouterd-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
