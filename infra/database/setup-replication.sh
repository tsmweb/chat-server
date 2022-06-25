#!/bin/bash

if [ "$REPLICATE_MASTER" = "on" ]; then

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
CREATE USER $REPLICATE_USER REPLICATION LOGIN CONNECTION LIMIT 100 ENCRYPTED PASSWORD '$REPLICATE_PASSWORD';
EOSQL

cat >> ${PGDATA}/postgresql.conf <<-EOREP
wal_level = logical
max_wal_senders = $PG_MAX_WAL_SENDERS
wal_keep_segments = $PG_WAL_KEEP_SEGMENTS
max_replication_slots = 1

primary_conninfo = 'host=${REPLICATE_FROM} port=5432 user=${REPLICATE_USER} password=${REPLICATE_PASSWORD} sslmode=prefer sslcompression=0 gssencmode=prefer krbsrvname=postgres target_session_attrs=any'
promote_trigger_file = '/tmp/promote_to_master'
hot_standby = on
EOREP

fi