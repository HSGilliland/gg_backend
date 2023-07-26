#!/bin/bash

# create default GG database
set -e 
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL 
   CREATE DATABASE $GG_DB; 
   GRANT ALL PRIVILEGES ON DATABASE $GG_DB TO $POSTGRES_USER; 
EOSQL

# setup GG tables
psql --dbname "$GG_DB" -f /sql/gg.sql