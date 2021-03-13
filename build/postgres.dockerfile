FROM postgres:alpine

COPY scripts/init.sql /init.sql

COPY scripts/initdb.sh /docker-entrypoint-initdb.d/initdb.sh
RUN chmod +x /docker-entrypoint-initdb.d/initdb.sh

EXPOSE 5432
