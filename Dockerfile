FROM hub.expvent.com.cn:1111/expvent/builder/golang:1.20 as builder

WORKDIR /src

COPY . /src

RUN make build


FROM hub.expvent.com.cn:1111/expvent/base/ubuntu:20.04
WORKDIR /
COPY --from=builder /src/bin/grpctest /grpctest
ENTRYPOINT /grpctest
