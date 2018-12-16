# Runs my blockchain implementation
FROM ubuntu:18.04
LABEL maintainer="Martin Haug <martinhaug@smart-mail.de>"

#RUN apt-get -y update && apt-get -y upgrade && apt-get -y install net-tools

WORKDIR /app
COPY awesomeProject .
COPY app.sh .


ENTRYPOINT ["./app.sh"]
CMD [""]
