# -*- coding=utf-8 -*-

from models import dbconnect
from sqlalchemy import Table

dbsession, dbmodel, metadata = dbconnect()


class Menus(dbmodel):
    __table__ = Table('ledger_menus', metadata, autoload=True)

    def find_by_menus_level1(self):
        results = dbsession.query(Menus).filter_by(m_level=1).order_by(Menus.seq_sort).all()
        results_list = list()
        for i in results:
            results_list.append({c.name: getattr(i, c.name) for c in self.__table__.columns})

        if len(results_list) == 0:
            return False
        return results_list

    def find_by_menus_level2(self, m_id):
        results = dbsession.query(Menus).filter_by(parent_menu_id=m_id).order_by(Menus.seq_sort).all()
        results_list = list()
        for i in results:
            results_list.append({c.name: getattr(i, c.name) for c in self.__table__.columns})

        if len(results_list) == 0:
            return False
        return results_list
