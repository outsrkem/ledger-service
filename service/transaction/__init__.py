# -*- coding=utf-8 -*-
from decimal import Decimal


def statistical_trim(data=list()):
    # data 固定内容如下
    # data = [{"amount_total": "10000005.0000", "money_title": "餐饮", "money_type": 1001, "record_period": "2021-10"}]
    response = dict()
    rows = list()
    amount_total = list()
    amount_sum = Decimal()
    if data:
        details, period, row = list(), '', dict()
        for i in data:
            amount_total.append(i["amount_total"])
            i["amount_total"] = str(i["amount_total"].quantize(Decimal('0.000')))
            period = i["record_period"]
            details.append(i)

        for i in details:
            row[i["money_type"]] = i["amount_total"]
        rows.append(row)

        for i in amount_total:  # 计算总金额
            amount_sum = amount_sum + i
        '''四舍五入，保留3位小数，实际计算按4位精度'''
        amount_sum = amount_sum.quantize(Decimal('0.000'))

        response["details"], response["period"], response["rows"], response["amount_sum"] = \
            details, period, rows, str(amount_sum)

    return response
