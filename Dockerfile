FROM golang:1.21

WORKDIR /MEDODS-test-task

COPY . /MEDODS-test-task

RUN go mod download

RUN go build -v

CMD [ "./MEDODS-test-task" ]