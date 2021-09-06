# -*- coding=utf-8 -*-

from models import dbconnect
from sqlalchemy import Table

dbsession, dbmodel, metadata = dbconnect()


class MoneyType(dbmodel):
    __table__ = Table('ledger_money_type', metadata, autoload=True)

    def find_by_money_type(self, tid):
        # 根据id查询单条
        result = dbsession.query(MoneyType).filter_by(id=tid).all()
        if len(result) == 0:
            return False
        return {c.name: getattr(result[0], c.name) for c in self.__table__.columns}

    @staticmethod
    def find_by_all_money_type():
        # 查询所有
        sql = '''SELECT t_mt.id, t_mt.uid, t_mt.category AS category_id, IF (mt.title IS NULL, t_mt.title, mt.title) 
    AS category, t_mt.title, t_mt.status, t_mt.describes, t_mt.create_time, t_mt.update_time FROM ledger_money_type mt 
    RIGHT JOIN ledger_money_type t_mt ON (t_mt.category = mt.id)'''
        results = dbsession.execute(sql)
        result = [dict(zip(result.keys(), result)) for result in results]
        if len(result) > 0:
            return result
        return False

    @staticmethod
    def inst_money_type(raw):
        # 添加
        try:
            for i in raw["raw_data"]:
                dbsession.add(
                    MoneyType(uid=raw["uid"], title=i["title"], category=i["category"], describes=i["category"],
                              create_time=raw["create_time"], update_time=raw["update_time"]))
            dbsession.commit()
            return True
        except Exception as e:
            dbsession.rollback()
            print(e)
            return False

    @staticmethod
    def del_money_type(raws):
        # 删除
        try:
            for type_id in raws["type_id"]:
                dbsession.query(MoneyType).filter_by(id=type_id).delete()
            dbsession.commit()
            return True
        except Exception as e:
            dbsession.rollback()
            print(e)
            return False

    @staticmethod
    def update_money_type(raw):
        # 更新
        try:
            row = dbsession.query(MoneyType).filter_by(id=raw["type_id"]).first()
            if not row:
                print("Type does not exist, type_id: " + str(raw["type_id"]))
                return False
            row.title = raw["title"]
            row.status = raw["status"]
            row.describes = raw["describes"]
            row.update_time = raw["update_time"]
            dbsession.commit()
            return True
        except Exception as e:
            dbsession.rollback()
            print(e)
            return False
