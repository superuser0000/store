FROM golang

RUN mkdir -p /go/src/qxklmrhx7qkzais6.onion/Tochka/tochka-free-market
WORKDIR /go/src/qxklmrhx7qkzais6.onion/Tochka/tochka-free-market

ADD . /go/src/qxklmrhx7qkzais6.onion/Tochka/tochka-free-market

RUN ./scripts/build.sh