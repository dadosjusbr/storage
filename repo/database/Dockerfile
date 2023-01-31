FROM postgres

WORKDIR /

ENV POSTGRES_USER root
ENV POSTGRES_PASSWORD dadosjusbr_test
ENV POSTGRES_DB dadosjusbr_test

COPY ./init_db.sql /docker-entrypoint-initdb.d/


USER root

