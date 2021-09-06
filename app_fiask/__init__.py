# -*- coding=utf-8 -*-
from urllib import request
from flask import Flask, request, session, current_app
from flask_sqlalchemy import SQLAlchemy
from service import response_body
from settings import config
import time

__app = Flask(__name__)
__app.config.from_object(config["default"])

db = SQLAlchemy(__app)


@__app.errorhandler(404)
def page_not_found(e):
    request_time = int(round(time.time() * 1000))
    msg = f"{str(e)} URL: {request.path} "
    return {"meta_info": {"code": 404, "msg": msg, "request_time": request_time}}, 404


@__app.errorhandler(500)
def server_error(e):
    request_time = int(round(time.time() * 1000))
    current_app.logger.error("Server internal error")
    return {"meta_info": {"code": 500, "msg": "Internal server error", "request_time": request_time}}, 500


@__app.before_request
def before():
    url = request.path
    pass_list = [
        '/',
        '/api/v1/ledger/common/user/login',
        '/api/v1/ledger/common/user/logout',
        '/api/v1/ledger/common/user/register',
    ]
    if url not in pass_list:
        if not session.get('islogin'):
            return response_body(401, 'Invalid login status'), 401


def init_app():
    return __app
