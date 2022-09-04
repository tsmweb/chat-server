FROM postgres:11
LABEL maintainer "Tiago S. Martins (tiago.tsmweb@gmail.com)"
ENV PG_MAX_WAL_SENDERS 3
ENV PG_WAL_REPLICATION_SLOTS 3
ENV PG_WAL_KEEP_SEGMENTS 8

COPY setup-replication.sh /docker-entrypoint-initdb.d/
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint-initdb.d/setup-replication.sh /docker-entrypoint.sh
