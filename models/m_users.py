# -*- coding=utf-8 -*-

from models import dbconnect
from sqlalchemy import Table
from service.utility import now_timestamp

dbsession, dbmodel, metadata = dbconnect()


class Users(dbmodel):
    __table__ = Table('ledger_user', metadata, autoload=True)

    @staticmethod
    def user_register(username, passwd, mobile, describes):
        now_time = now_timestamp()
        account = mobile
        # noinspection PyBroadException
        try:
            dbsession.add(
                Users(account=account, username=username, password=passwd, mobile=mobile, describes=describes,
                      create_time=now_time,
                      update_time=now_time))
            dbsession.commit()
            return True
        except Exception as e:
            dbsession.rollback()
            return False

    def find_by_userinfo(self, account):
        result = dbsession.query(Users).filter_by(account=account).all()
        if len(result) == 0:
            return False
        return {c.name: getattr(result[0], c.name) for c in self.__table__.columns}
