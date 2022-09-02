from models import dbconnect
from sqlalchemy import Table

dbsession, dbmodel, metadata = dbconnect()


class Options(dbmodel):
    __table__ = Table('ledger_options', metadata, autoload=True)

    def find_by_options(self, user_id):
        """查询用户配置"""
        results = dbsession.query(Options).filter_by(uid=user_id).order_by(Options.id).all()
        results_list = list()
        if len(results) == 0:
            return False
        for i in results:
            results_list.append({c.name: getattr(i, c.name) for c in self.__table__.columns})
        return results_list
