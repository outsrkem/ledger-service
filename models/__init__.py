# -*- coding=utf-8 -*-
from sqlalchemy import MetaData


def dbconnect():
    from app_flask import db
    db.get_engine()
    dbsession = db.session
    dbmodel = db.Model
    metadata = MetaData(bind=db.engine)
    return dbsession, dbmodel, metadata


