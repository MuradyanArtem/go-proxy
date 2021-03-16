FROM postgres:alpine

ADD scripts/init.sql /init.sql

ADD scripts/create.sh /docker-entrypoint-initdb.d/create.sh
RUN chmod +x /docker-entrypoint-initdb.d/create.sh

EXPOSE 5432
