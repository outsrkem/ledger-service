# Step 1
FROM outsrkem/alpine:3.19.1-golng1.22.0-v1
ARG GO_VERSION="go1.22.0"
ARG APP_NAME
ARG APP_VERSION
ARG APP_REVISION
WORKDIR /opt/$APP_NAME


COPY . /opt/$APP_NAME

# -trimpath 移除源代码中的文件路径信息
# -ldflags -s：不生成符号表 -w：不生成DWARF调试信息
ARG LD_PATH="$APP_NAME/src/config"
ARG LD_FLAGS="-X $LD_PATH.Version=${APP_VERSION} -X $LD_PATH.GoVersion=${GO_VERSION} -X $LD_PATH.GitCommit=${APP_REVISION}"
RUN go build -trimpath  -ldflags "-s -w $LD_FLAGS" -o output/$APP_NAME src/main/main.go

# RUN upx -9 output/app
RUN output/$APP_NAME -version
RUN cp ledger.yaml output
RUN cp docker-entrypoint.sh output/entrypoint.sh

# Step 2
FROM alpine:3.19.1
ARG APP_NAME
ARG APP_VERSION
ARG APP_REVISION

COPY --from=0 /opt/$APP_NAME/output/ /usr/local/bin
RUN mkdir "/etc/ledger" && mv /usr/local/bin/ledger.yaml /etc/ledger/ledger.yaml

ENV APP_NAME=$APP_NAME APP_VERSION=$APP_VERSION APP_REVISION=$APP_REVISION

ENTRYPOINT ["entrypoint.sh"]
