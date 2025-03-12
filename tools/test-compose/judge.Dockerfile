FROM golang:1.24-bookworm

RUN mkdir /app
ADD . /source
RUN cd /source && go build -o /app/judged . 

CMD [ "bash", "-c", "cd /app && ./judged" ]
