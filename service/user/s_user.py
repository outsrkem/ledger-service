# -*- coding=utf-8 -*-
from flask import session
from service import response_body
from werkzeug.security import check_password_hash, generate_password_hash
from models.m_users import Users
from settings import Logger

_log = Logger()


def user_register(data):
    # 用户注册
    username = data["username"].strip()
    password = data["password"].strip()
    mobile = data["mobile"].strip()
    describes = data["describes"]

    passwd = generate_password_hash(password, method='pbkdf2:sha256', salt_length=16)
    row = Users().user_register(username, passwd, mobile, describes)
    if row:
        return response_body(201)
    return response_body(500, 'User registration failure!')


def user_login(data):
    account = data["account"].strip()
    password = data["password"].strip()
    result = Users().find_by_userinfo(account)
    if result and check_password_hash(result["password"], password):
        session['islogin'] = True
        session['user_id'] = result["id"]
        session['account'] = result["account"]
        session['user_name'] = result["username"]
        session['mobile'] = result["mobile"]
        session['status'] = result["status"]
        payload = {"user_id": result["id"], "account": result["account"], "user_name": result["username"],
                   "describes": result["describes"],
                   "update_time": result["update_time"], "status": result["status"], "token": ""}

        _log.logger.info("Successful user login. userinfo: %s" % payload)
        return response_body(200, '', payload)

    _log.logger.warning("User login failure. login account: %s" % data["account"])
    return response_body(403, 'Login Error')
