# -*- coding=utf-8 -*-
from .r_common import common as comm
from .r_root import root
from .r_recedisbus import recedisbu


def reg_blueprint(app):
    """
    注册蓝图,所有接口前缀都在这里
    :param app:
    """
    app.register_blueprint(root, url_prefix='/')
    app.register_blueprint(comm, url_prefix='/api/v1/common')
    app.register_blueprint(recedisbu, url_prefix='/api/v1/recedisbu')
