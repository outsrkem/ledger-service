-- occ_time对应毫秒时间戳
ALTER TABLE ledger_transaction ADD COLUMN occ_at BIGINT NULL AFTER occ_time;