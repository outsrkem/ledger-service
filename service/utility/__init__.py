# -*- coding=utf-8 -*-
import time
import json
from settings import Logger

_log = Logger()


def now_timestamp():
    _now_time = int(round(time.time() * 1000))
    return _now_time


def to_json(t_str):
    """
    字符串转json并校验
    :param t_str:
    :return:
    """
    try:
        _ = json.loads(t_str)
    except Exception as e:
        _log.logger.error("The json formatting fails. Procedure, %s, str: %s" % (e, t_str))
        return False
    return json.loads(t_str)
