from service.utility import now_time_timestamp


def response_body(code=200, msg='', res=''):
    now_time = now_time_timestamp()
    _response = dict()
    if not msg:
        msg = 'successfully'
    meta_info = {"res_code": code, "res_msg": msg, "request_time": now_time}
    _response["meta_info"] = meta_info

    if res:
        _response["response"] = res

    return _response
