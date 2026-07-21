# 适用于occ_time字段转occ_at字段
import pymysql
from datetime import datetime

DB_CFG = {
    "host": "172.19.53.125",
    "port": 38455,
    "user": "ledger",
    "password": "123456",
    "database": "ledgerdb",
    "charset": "utf8mb4"
}
BATCH_SIZE = 1000
COMMIT_INTERVAL = 5  # 每5批提交一次，减少事务开销

def convert_time_to_ms(time_str: str) -> int:
    """解析 2025-08-27T07:51:31+0800 ISO8601 字符串 → 毫秒时间戳"""
    dt = datetime.strptime(time_str, "%Y-%m-%dT%H:%M:%S%z")
    return int(dt.timestamp() * 1000)

def batch_fill_occ_at():
    conn = pymysql.connect(**DB_CFG)
    conn.autocommit(False)
    cursor = conn.cursor()

    last_kid = 0
    total_updated = 0
    batch_count = 0
    error_ids = []

    try:
        while True:
            # 增加条件：只查 occ_at 为null或0的行
            query_sql = """
                SELECT kid, occ_time 
                FROM ledger_transaction 
                WHERE kid > %s 
                  AND (occ_at IS NULL OR occ_at = 0)
                ORDER BY kid 
                LIMIT %s
            """
            print(f"\n===== 查询SQL =====")
            print(query_sql.strip())
            print(f"查询参数：last_kid={last_kid}, batch_size={BATCH_SIZE}")

            cursor.execute(query_sql, (last_kid, BATCH_SIZE))
            rows = cursor.fetchall()

            if not rows:
                print("\n所有待更新数据读取完毕，准备提交剩余事务")
                break

            update_batch = []
            for kid, occ_time in rows:
                try:
                    ms_ts = convert_time_to_ms(occ_time)
                    update_batch.append((ms_ts, kid))
                except Exception as e:
                    print(f"解析失败 | kid={kid}, occ_time={occ_time}, 错误：{str(e)}")
                    error_ids.append(kid)

            if update_batch:
                update_sql = "UPDATE ledger_transaction SET occ_at = %s WHERE kid = %s"
                print(f"\n===== 更新SQL =====")
                print(update_sql)
                print(f"本批待更新数量：{len(update_batch)}")

                cursor.executemany(update_sql, update_batch)
                total_updated += len(update_batch)
                batch_count += 1

                # 达到批次阈值提交
                if batch_count % COMMIT_INTERVAL == 0:
                    conn.commit()
                    print(f"已提交 {batch_count} 批，累计更新 {total_updated} 条")

            # 更新游标为当前最大kid
            last_kid = rows[-1][0]

        # 循环结束提交剩余未提交数据
        conn.commit()
        print(f"\n==================== 任务完成 ====================")
        print(f"成功更新总行数：{total_updated}")

        if error_ids:
            unique_err = list(set(error_ids))
            print(f"解析失败记录总数：{len(unique_err)}")
            print(f"前10条异常kid：{unique_err[:10]}")

    except Exception as err:
        conn.rollback()
        print(f"\n程序异常！所有操作已回滚，错误详情：{err}")
    finally:
        cursor.close()
        conn.close()
        print("数据库连接已关闭")

if __name__ == "__main__":
    batch_fill_occ_at()