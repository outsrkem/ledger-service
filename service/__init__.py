from service.utility import now_timestamp


def response_body(code=200, msg='', res=''):
    """
    response 统一的response返回体模板
    :param code: 状态码
    :param msg:  返回消息
    :param res:  返回的内容实体
    :return: response
    """
    now_time = now_timestamp()
    _response = dict()
    if not msg:
        msg = 'successfully'
    meta_info = {"res_code": code, "res_msg": msg, "request_time": now_time}
    _response["meta_info"] = meta_info

    if res:
        _response["payload"] = res

    return _response
