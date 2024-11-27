#!/bin/bash
set -e

echo "Setting up replica node..."

# 等待主節點完全啟動
until pg_isready -h postgres_master -U replica; do
  echo "Waiting for PostgreSQL master to start..."
  sleep 2
done

# 初始化複製，使用已創建的插槽
pg_basebackup -h postgres_master -U replica -D /var/lib/postgresql/data -P -R --slot=replica_slot -X stream

echo "Replica setup complete."
