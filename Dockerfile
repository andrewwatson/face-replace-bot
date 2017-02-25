FROM golang:1.8

RUN apt-get update -y

RUN apt-get install -y libopencv-dev

RUN go get github.com/zikes/chrisify && go get github.com/lazywei/go-opencv

RUN cd $GOPATH/src/github.com/zikes/chrisify && go build

# RUN mkdir -p /opt/facebot && cp -r $GOPATH/src/github.com/zikes/chrisify/faces /opt/facebot/faces
RUN mkdir -p /opt/facebot

COPY faces /opt/facebot/faces

RUN cp -r $GOPATH/src/github.com/zikes/chrisify/haarcascade_frontalface_alt.xml /opt/facebot/

# RUN go get github.com/andrewwatson/face-replace-bot
COPY . /go/src/github.com/andrewwatson/face-replace-bot

RUN cd /go/src/github.com/andrewwatson/face-replace-bot && go get -d ./... && go install

WORKDIR /opt/facebot

CMD ["/go/bin/face-replace-bot"]
