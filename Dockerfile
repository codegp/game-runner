FROM alpine

RUN mkdir /app
COPY gamerunner/gamerunner /app/gamerunner
WORKDIR /app
RUN ls
EXPOSE 9000

CMD ["./gamerunner"]
