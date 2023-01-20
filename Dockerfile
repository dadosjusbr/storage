# DockerFile responsável por pegar o script SQL e gerar as tabelas no banco postgres

FROM postgres

WORKDIR /database

COPY ./init_db.sql /docker-entrypoint-initdb.d/

USER root