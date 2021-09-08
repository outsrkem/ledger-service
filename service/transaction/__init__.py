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
            i["amount_total"] = i["amount_total"]
            period = i["record_period"]
            details.append(i)

        for i in details:
            row[i["money_type"]] = i["amount_total"]
        rows.append(row)

        for i in amount_total:  # 计算总金额
            amount_sum = amount_sum + i
        '''四舍五入，保留3位小数，实际计算按4位精度'''
        amount_sum = amount_sum

        response["details"], response["period"], response["rows"], response["amount_sum"] = \
            details, period, rows, amount_sum

    return response


def formatting_type_decimal(rows, k):
    """
    格式化decimal类型为字符串
    :param rows: 需要格式化的列表数据：[{"a": Decimal('24.1000')}, {"a": Decimal('24.1000')}]
    :param k: 转换类型的键
    :return: 返回和入参一样的列表
    """
    _rows = list()
    for i in rows:
        i[k] = str(i[k])
        _rows.append(i)
    return _rows


def statistical_type_decimal_to_str(data):
    """
    用于decimal 转 str ，只用于 日统计，月统计，年统计
    :param data:
    :return: 原格式
    """
    _data = dict()
    _details = list()
    _rows = list()
    _amount_sum = ''
    for i in data["details"]:
        i["amount_total"] = str(i["amount_total"])
        _details.append(i)

    for i in data["rows"]:
        _rows_dict = dict()
        for key, value in i.items():
            _rows_dict[key] = str(value.quantize(Decimal('0.000')))
        _rows.append(_rows_dict)

    _amount_sum = str(data["amount_sum"].quantize(Decimal('0.000')))
    _data["details"] = _details
    _data["period"] = data["period"]
    _data["rows"] = _rows
    _data["amount_sum"] = _amount_sum
    return _data
