# -*- coding=utf-8 -*-
from flask import session
from service import response_body
from werkzeug.security import check_password_hash, generate_password_hash
from models.m_users import Users
from settings import Logger
from service.utility import to_json

_log = Logger()


def user_register(data):
    # 用户注册
    username = data["username"].strip()
    password = data["password"].strip()
    mobile = data["mobile"].strip()
    account = mobile
    describes = data["describes"]

    """重复注册查询校验"""
    is_register = Users().find_by_userinfo(account)
    if is_register:
        return response_body(409, 'A user with a mobile phone number is registered repeatedly.')

    """用户注册"""
    passwd = generate_password_hash(password, method='pbkdf2:sha256', salt_length=16)
    row = Users().user_register(account, username, passwd, mobile, describes)
    if row:
        _log.logger.info("User registration succeeded,account: %s" % account)
        return response_body(201)
    return response_body(500, 'User registration failure!')


def user_login(data):
    """处理登陆参数异常"""
    data = to_json(data)
    if not data:
        return response_body(400, 'The json formatting fails or the parameter is abnormal')

    try:
        account = data["account"].strip()
    except KeyError as e:
        account = False
        _log.logger.error("The login account parameter is abnormal: %s" % e)

    try:
        password = data["password"].strip()
    except KeyError as e:
        password = False
        _log.logger.error("The login password parameter is abnormal: %s" % e)

    """无账号或密码"""
    if not account:
        _log.logger.error("Login to terminate.")
        return response_body(400, '''User login parameter error,No [account] number or [password].''')
    elif not password:
        return response_body(400, '''User login parameter error,No [account] number or [password].''')

    """登陆校验"""
    result = Users().find_by_userinfo(account)

    if result and check_password_hash(result["password"], password):
        session['is_login'] = True
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
