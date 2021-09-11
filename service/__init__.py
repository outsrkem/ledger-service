from service.utility import now_timestamp


def response_body(code=200, msg='', payload=''):
    """
    response 统一的response返回体模板
    :param code: 状态码
    :param msg:  返回消息
    :param payload:  返回的内容实体
    :return: response
    """
    now_time = now_timestamp()
    _response = dict()
    if not msg:
        meta_info = {"res_code": code, "request_time": now_time}
    else:
        meta_info = {"res_code": code, "res_msg": msg, "request_time": now_time}

    _response["meta_info"] = meta_info

    if payload:
        _response["payload"] = payload

    return _response
