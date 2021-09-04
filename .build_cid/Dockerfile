FROM harbor.hub.com/library/python:3.8
ARG version

WORKDIR /home/ledger-service

ADD . .

RUN pip install -r requirements.txt

ENV LEDGER_SERVICE_VERSION $version

CMD ["gunicorn","-c","gunicorn_config.py","main:app"]

