-- -
-- 请手工创建数据库
-- -
CREATE DATABASE /*!32312 IF NOT EXISTS*/ ledgerdb /*!40100 DEFAULT CHARACTER SET utf8mb4 */;
GRANT ALL ON ledgerdb.* TO ledger@'%' IDENTIFIED BY '123456';
GRANT ALL ON ledgerdb.* TO ledger@'localhost' IDENTIFIED BY '123456';