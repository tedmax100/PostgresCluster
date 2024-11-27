#!/bin/bash
set -e

echo "Setting up master node..."

# 等待 PostgreSQL 完全啟動
until pg_isready -U admin; do
  echo "Waiting for PostgreSQL master to start..."
  sleep 2
done

# 創建複製用戶
psql -U admin -d mydb <<-EOSQL
  CREATE ROLE replica WITH REPLICATION LOGIN PASSWORD 'replica_password';
  SELECT * FROM pg_create_physical_replication_slot('replica_slot');
EOSQL

echo "Master setup complete."
