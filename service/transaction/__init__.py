# -*- coding=utf-8 -*-

def statistical_trim(data=list()):
    # data 固定内容如下
    # data = [{"amount_total": "10000005.0000", "money_title": "餐饮", "money_type": 1001, "record_period": "2021-10"}]
    response = dict()
    if data:
        details, period, row = list(), '', dict()
        for i in data:
            i["amount_total"] = str(i["amount_total"])
            period = i["record_period"]
            details.append(i)

        for i in details:
            row[i["money_type"]] = i["amount_total"]
        response["details"], response["period"], response["row"] = details, period, row
    return response
