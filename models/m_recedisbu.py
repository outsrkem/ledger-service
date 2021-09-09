# -*- coding=utf-8 -*-

from models import dbconnect
from sqlalchemy import Table

dbsession, dbmodel, metadata = dbconnect()


class RecedisbuStatement(dbmodel):
    __table__ = Table('ledger_recedisbu_statement', metadata, autoload=True)

    @staticmethod
    def find_by_income_and_daily_statistical(y=0, m=0, d=0):
        sql = f'''SELECT CONCAT(rs.occ_year,'-',IF (rs.occ_month < 10,CONCAT('0', rs.occ_month),rs.occ_month),'-',
    IF (rs.occ_day < 10,CONCAT('0', rs.occ_day),rs.occ_day)) AS record_period,SUM(rs.amount) AS amount_total,
    IF (mt.category = 1,mt.id,mt.category) AS money_type,
     (SELECT title FROM ledger_money_type WHERE id = money_type ) AS money_title FROM ledger_recedisbu_statement rs
    LEFT JOIN ledger_money_type mt ON (rs.mtid = mt.id)
    WHERE rs.occ_year = {y} AND rs.occ_month = {m} AND rs.occ_day = {d} GROUP BY money_type;'''
        results = dbsession.execute(sql)
        # 查询结果是一个list，
        # 在这个list中包含着许多 tuple，
        # 他们看似是tuple，实则并不是
        # 而是一个特殊的类型"<class ‘sqlalchemy.util._collections.result’>"
        # 这是一个 AbstractKeyedTuple 对象
        # 它拥有一个 keys() 方法，可以取到所有的key
        # 我们可以通过这个方法将查询结果转换为字典
        # https://blog.csdn.net/qq_43193386/article/details/115633706
        result = [dict(zip(result.keys(), result)) for result in results]
        if len(result) > 0:
            return result
        return False

    @staticmethod
    def find_by_income_and_monthly_statistical(y=0, m=0):
        sql = f'''SELECT CONCAT(rs.occ_year,'-',IF (rs.occ_month < 10,CONCAT('0', rs.occ_month),rs.occ_month)) 
    AS record_period,SUM(rs.amount) AS amount_total,
    IF (mt.category = 1,mt.id,mt.category) AS money_type,
     (SELECT title FROM ledger_money_type WHERE id = money_type ) AS money_title FROM ledger_recedisbu_statement rs
    LEFT JOIN ledger_money_type mt ON (rs.mtid = mt.id)
    WHERE rs.occ_year = {y} AND rs.occ_month = {m} GROUP BY money_type;'''
        results = dbsession.execute(sql)
        result = [dict(zip(result.keys(), result)) for result in results]
        if len(result) > 0:
            return result
        return False

    @staticmethod
    def find_by_income_and_annual_statistical(y=''):
        sql = f'''SELECT CONCAT(rs.occ_year) AS record_period,SUM(rs.amount) AS amount_total,
    IF (mt.category = 1,mt.id,mt.category) AS money_type,
     (SELECT title FROM ledger_money_type WHERE id = money_type ) AS money_title FROM ledger_recedisbu_statement rs
    LEFT JOIN ledger_money_type mt ON (rs.mtid = mt.id)
    WHERE rs.occ_year = {y} GROUP BY money_type;'''
        results = dbsession.execute(sql)
        result = [dict(zip(result.keys(), result)) for result in results]
        if len(result) > 0:
            return result
        return False

    @staticmethod
    def inst_journal_account(raw):
        # 添加流水
        try:
            for i in raw["raw_data"]:
                occ_year = i["occ_date"].split("-")[0]
                occ_month = i["occ_date"].split("-")[1]
                occ_day = i["occ_date"].split("-")[2]
                dbsession.add(
                    RecedisbuStatement(uid=raw["uid"], mtid=i["mtid"], occ_year=occ_year, occ_month=occ_month,
                                       amount=i["amount"], total=i["total"], describes=i["describes"],
                                       occ_day=occ_day, create_time=raw["create_time"],
                                       update_time=raw["update_time"]))
            dbsession.commit()
            return True
        except Exception as e:
            print(e)
            dbsession.rollback()
            return False

    def find_by_all_month_record(self, y, m):
        # 查询当月流水
        result = dbsession.query(RecedisbuStatement).filter_by(occ_year=y, occ_month=m).group_by("occ_day").all()
        result_list = list()
        for row in result:
            ww = {c.name: getattr(row, c.name) for c in self.__table__.columns}
            result_list.append(ww)

        if len(result) > 0:
            return result_list
        return False

    def find_by_all_annual_record(self, y):
        # 查询当年流水
        result = dbsession.query(RecedisbuStatement).filter_by(occ_year=y).group_by("occ_month").all()
        result_list = list()
        for row in result:
            ww = {c.name: getattr(row, c.name) for c in self.__table__.columns}
            result_list.append(ww)

        if len(result) > 0:
            return result_list
        return False
