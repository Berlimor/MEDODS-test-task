FROM golang:1.21

WORKDIR /MEDODS-test-task

EXPOSE 8000

COPY . /MEDODS-test-task

RUN go mod download

RUN go build -v

CMD [ "./MEDODS-test-task" ]