# -*- coding=utf-8 -*-
import time


def now_timestamp():
    _now_time = int(round(time.time() * 1000))
    return _now_time
