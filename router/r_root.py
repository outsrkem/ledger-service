# -*- coding=utf-8 -*-
from flask import Blueprint, jsonify
import time

root = Blueprint('root', __name__)


@root.route("/", methods=['GET', 'HEAD'])
def index():
    """
    接口检测，健康检测等
    :return:
    """
    request_time = int(round(time.time() * 1000))
    return {"meta_info": {"res_code": 200, "res_msg": "successfully", "request_time": request_time}}


@root.route("/test", methods=['GET'])
def test():
    return {}
