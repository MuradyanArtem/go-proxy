FROM postgres:alpine

ADD scripts/init.sql /init.sql

ADD scripts/create_db.sh /docker-entrypoint-initdb.d/create_db.sh
RUN chmod +x /docker-entrypoint-initdb.d/create_db.sh

EXPOSE 5432
