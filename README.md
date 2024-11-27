# PostgresReplication
-----

## postgres_primary
配置 PostgreSQL 啟動參數，包括：
wal_level=replica: 啟用複製的 WAL 日誌級別。
hot_standby=on: 允許從節點以唯讀模式處理查詢。
max_wal_senders=10: 允許最多 10 個 WAL 發送器進行複製。
max_replication_slots=10: 設置最多 10 個複製插槽。
hot_standby_feedback=on: 啟用熱備反饋以避免因 GC 而導致查詢取消。

## postgres_replica
pg_basebackup: 用於從主節點進行基礎資料備份並啟用 Streaming Replication：
--pgdata: 指定數據目錄。
-R: 自動生成從節點的複製配置文件。
--slot: 使用 replication_slot 複製插槽。
--host: 指定主節點的位置。
--port: 主節點的Port為 5432。
使用 bash 腳本保證在主節點可用之前循環等待連線。
一旦備份完成，啟動從節點的 postgres 服務。
