# -*- coding=utf-8 -*-
from flask import session

from models import dbconnect
from sqlalchemy import Table, func
from service.utility import now_timestamp

dbsession, dbmodel, metadata = dbconnect()


class Users(dbmodel):
    __table__ = Table('ledger_user', metadata, autoload=True)

    @staticmethod
    def user_register(account, username, passwd, mobile, describes):
        """注册用户，插入数据库"""
        now_time = now_timestamp()
        # noinspection PyBroadException
        try:
            dbsession.add(
                Users(account=account, username=username, password=passwd, mobile=mobile, describes=describes,
                      create_time=now_time,
                      update_time=now_time))
            dbsession.commit()
            return True
        except Exception as e:
            print(e)
            dbsession.rollback()
            return False

    def find_by_userinfo(self, account):
        result = dbsession.query(Users).filter_by(account=account).all()
        if len(result) == 0:
            return False
        return {c.name: getattr(result[0], c.name) for c in self.__table__.columns}

    @staticmethod
    def find_by_users_count():
        """查询用户总数目"""
        count = dbsession.query(func.count(Users.id)).scalar()
        return count

    def find_by_users(self, page, per_page):
        """查用用户列表信息"""
        # results = dbsession.query(Users).paginate(page=page, per_page=per_page, error_out=False).items
        results = dbsession.query(Users).order_by(Users.id.desc()).limit(per_page).offset(
            (page - 1) * per_page).all()
        results_list = list()
        if len(results) == 0:
            return False
        for i in results:
            results_list.append({c.name: getattr(i, c.name) for c in self.__table__.columns})
        # 删除结果中的密码信息
        for i in results_list:
            i["is_edited"] = 1
            if i["id"] == session.get('user_id'):
                i["is_edited"] = 0
            i.pop("password")
        return results_list
