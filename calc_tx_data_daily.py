#! /usr/bin/python3
# coding: utf-8

import time
from datetime import datetime, timedelta
import sys
import pymysql
import requests

DATA_URL = "http://ip:port/api/custom/current/price_to_currency"
SQL_HOST = ""
SQL_USER = ""
SQL_PASSWD = ""
SQL_DB = ""

conn = pymysql.connect(
    host= SQL_HOST, 
    user=SQL_USER, 
    password=SQL_PASSWD, 
    database=SQL_DB, 
    charset="utf8")
cursor = conn.cursor()

def calc_tx_amount_count(start_time, end_time):
    """计算昨日交易总金额"""
    while start_time < end_time:
        time1 = start_time.timestamp()
        # every 10 Min
        time2 = (start_time + timedelta(minutes=10)).timestamp()
        print("calc time duration:", time1, time2)
        sql1 = "select distinct token_symbol from dgateway_tx where time >= {} and time < {}".format(time1, time2)

        cursor.execute(sql1)
        ret =  cursor.fetchall()
        for i in ret:
            if i[0] == '':
                continue
            #获取amount
            sql2 = "select sum(amount) from dgateway_tx where time >= {} and time < {} and token_symbol = '{}'".format(time1, time2, i[0])
            cursor.execute(sql2)
            amount  = cursor.fetchone()[0]

            #获取count
            sql3 = "select count(*) from dgateway_tx where time >= {} and time < {} and token_symbol = '{}'".format(time1, time2, i[0])
            cursor.execute(sql3)
            count = cursor.fetchone()[0]

            #获取currency_amount
            sql5 = "select chain, symbol, decimals, relate_chain, relate_token_code from dgateway_token_info where symbol = '{}'".format(i[0])
            cursor.execute(sql5)
            token_info  = cursor.fetchone()
            symbol = token_info[1]
            decimals = token_info[2]

            if token_info[0] == "eth":
                sql6 = "select symbol, decimals from dgateway_token_info where chain = '{}' and token_code = '{}'".format(token_info[3], token_info[4])
                cursor.execute(sql6)
                token_info = cursor.fetchone()
                symbol = token_info[0]
                decimals = token_info[1]

            if token_info[0] == "eos" and token_info[1] == "XIN":
                usd_price = get_currency_price("BTC", "USD")
                cny_price = get_currency_price("BTC", "CNY")
                # eos私链的XIN币和美元的兑换比例是1000：1
                currency_amount = int(float(amount) * (cny_price / usd_price) / 1000)
            else: 
                price = get_currency_price(symbol, "CNY")
                currency_amount = int(float(price) * float(amount) / (10 ** float(decimals)))
        
            #新增数据
            time = start_time.strftime("%Y-%m-%d %H:%M:%S")
            print("time", time, "amount:", amount, "count:", count, "symbol:", i[0], "currency_amount:", currency_amount)
            sql4 = "insert into dgateway_tx_statistics (time, amount, count, symbol, currency_amount) values ('{}', {}, {}, '{}', {})".format(
                time, amount, count, i[0], currency_amount)
            try:
                cursor.execute(sql4)
                conn.commit()
            except Exception as err:
                print("insert into mysql, err:{}".format(err))
                conn.rollback()
        start_time = start_time + timedelta(minutes=10)


def get_currency_price(symbol, unit):
    print(symbol, unit)
    params = {
        "target": symbol,
        "unit": unit
    }
    #数据服务地址
    url =  DATA_URL
    try:
        res = requests.post(url=url, json=params)
        price = res.json()['price']
        return price
    except Exception as err:
        print(err)
        return 0


if __name__ == "__main__":
    start_time = datetime.strptime(sys.argv[1], "%Y-%m-%d")
    end_time = datetime.strptime(sys.argv[2], "%Y-%m-%d")
    calc_tx_amount_count(start_time, end_time)
    cursor.close()
    conn.close()

