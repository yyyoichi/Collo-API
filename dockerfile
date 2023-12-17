ARG USER_NAME="user"
ARG USER_ID="1000"
ARG USER_GID="1000"
ARG PASSWORD="p@55w0rd"
ARG HOST="http://localhost"
ARG APP_PORT="8765"

FROM node:lts-slim AS app
ARG HOST
ARG APP_PORT
ENV NEXT_PUBLIC_RPC_HOST=${HOST}:${APP_PORT}/rpc

WORKDIR /workspaces
COPY ./web ./
RUN npm install && npm run build


FROM golang:latest AS api
ARG USER_NAME
ARG USER_ID
ARG USER_GID
ARG PASSWORD
ARG HOST
ARG APP_PORT
ENV APP_PORT=${APP_PORT}
ENV CLIENT_HOST=${HOST}:${APP_PORT}

RUN mkdir -p /tmp/mecab 
COPY ./script/mecab-install.sh /tmp/mecab
RUN apt-get update && \
    apt-get install -y wget tar unzip sed build-essential && \
    chmod +x /tmp/mecab/mecab-install.sh && \
    sh /tmp/mecab/mecab-install.sh && \
    rm -rf /tmp/mecab

WORKDIR /tmpspaces
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./ ./
COPY --from=app /workspaces/out ./cmd/server/out

ENV CGO_LDFLAGS="-L/path/to/lib -lmecab -lstdc++"
ENV CGO_CFLAGS="-I/path/to/include"
RUN go build ./cmd/server

WORKDIR /workspaces

RUN groupadd -g ${USER_GID} ${USER_NAME} && \
    useradd -u ${USER_ID} -g ${USER_NAME} -m -s /bin/sh ${USER_NAME} && \
    echo ${USER_NAME}:${PASSWORD} | chpasswd && \
    rm -rf /go && \
    mv /tmpspaces/server ./ && \
    rm -rf /tmpspaces && \
    chown -R ${USER_ID}:${USER_GID} ./ 

USER ${USER_ID}:${USER_GID}
CMD [ "./server" ]

