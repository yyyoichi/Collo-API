sh /workspaces/script/mecab-install.sh

cd /workspaces/web && \
    sudo chown -R vscode:vscode node_modules && \
    npm config set @buf:registry https://buf.build/gen/npm/v1 && \
    npm install -g @buf/connectrpc_eliza.connectrpc_es @connectrpc/connect @connectrpc/connect-web @bufbuild/protoc-gen-es @connectrpc/protoc-gen-connect-es && \
    yarn && \
    sudo chown -R vscode:vscode /go/pkg && \
    go install github.com/bufbuild/buf/cmd/buf@latest && \
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest && \
    sudo mkdir -p /tmp/collo-network