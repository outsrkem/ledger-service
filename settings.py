# -*- coding=utf-8 -*-
import logging
import os
from logging import handlers

BASE_DIR = os.path.abspath(os.path.dirname(__file__))

# db
DB_HOST = "10.10.10.10"
DB_PORT = 3306
DB_NAME = "ledgerdb"
DB_USER_NAME = "ledger"
DB_PASSWD = "123456"

# 数据库连接池的大小。默认是数据库引擎的默认值 （通常是 5）
SQLALCHEMY_POOL_SIZE = 100
# 控制在连接池达到最大值后可以创建的连接数。当这些额外的连接使用后回收到连接池后将会被断开和抛弃。保证连接池只有设置的大小
SQLALCHEMY_MAX_OVERFLOW = 100
# 自动回收连接的秒数。这对MySQL是必须的，默认情况下MySQL会自动移除闲置8小时或者以上的连接,Flask-SQLAlchemy会自动地设置这个值为 2 小时
SQLALCHEMY_POOL_RECYCLE = 1200
# 如果设置成 True，SQLAlchemy 将会记录所有发到标准输出(stderr)的语句，这对调试很有帮助;默认为false；
SQLALCHEMY_ECHO = False
# 如果设置成 True (默认情况)，Flask-SQLAlchemy 将会追踪对象的修改并且发送信号。这需要额外的内存，如果不必要的可以禁用它。
SQLALCHEMY_TRACK_MODIFICATIONS = False

SECRET_KEY = 'wderqeyJ2Y29kZSI6ImxkbG4ifQ.Xr-Lbg.ojkAcx7BZx7590luvEIvhYASA_8'

# log
LOG_LEVEL = 'DEBUG'
LOG_DIR = os.path.join(BASE_DIR, "logs")
LOG_BACKUP_COUNT = 3
LOG_FORMAT = '[%(asctime)s] [%(levelname)s] [%(filename)s] [line:%(lineno)d] - %(message)s'
LOG_FILE_NAME = 'ledger-service.log'


class Config:
    SQLALCHEMY_DATABASE_URI = f"mysql+pymysql://{DB_USER_NAME}:{DB_PASSWD}@{DB_HOST}:{DB_PORT}/{DB_NAME}?charset=utf8"
    SQLALCHEMY_TRACK_MODIFICATIONS = SQLALCHEMY_TRACK_MODIFICATIONS
    SQLALCHEMY_POOL_SIZE = SQLALCHEMY_POOL_SIZE
    SQLALCHEMY_POOL_RECYCLE = SQLALCHEMY_POOL_RECYCLE
    SQLALCHEMY_MAX_OVERFLOW = SQLALCHEMY_MAX_OVERFLOW
    SECRET_KEY = SECRET_KEY
    SSL_DISABLE = False
    SQLALCHEMY_RECORD_QUERIES = True
    SQLALCHEMY_ECHO = False


config = {
    'default': Config
}


class Logger(object):
    if not os.path.exists(LOG_DIR):
        os.mkdir(LOG_DIR)

    filename = os.path.join(LOG_DIR, LOG_FILE_NAME)
    log_level = {
        "DEBUG": logging.DEBUG,
        "INFO": logging.INFO,
        "WARNING": logging.WARNING,
        "ERROR": logging.ERROR,
        "CRITICAL": logging.CRITICAL,
        "debug": logging.DEBUG,
        "info": logging.INFO,
        "warning": logging.WARNING,
        "error": logging.ERROR,
        "critical": logging.CRITICAL,
    }

    def __init__(self, filename=filename, level=LOG_LEVEL, when='D', back_count=LOG_BACKUP_COUNT, fmt=LOG_FORMAT):
        self.logger = logging.getLogger(filename)
        self.logger.propagate = False
        if not self.logger.handlers:
            '''设置日志格式'''
            format_str = logging.Formatter(fmt)

            '''设置日志级别'''
            self.logger.setLevel(self.log_level.get(level))
            th = handlers.TimedRotatingFileHandler(filename=filename, when=when, backupCount=back_count,
                                                   encoding='utf-8')
            '''设置文件里写入的格式'''
            th.setFormatter(format_str)

            '''把对象加到logger里'''
            self.logger.addHandler(th)

    '''参考博客: https://blog.csdn.net/momoyaoquaoaoao/article/details/87280440'''
    '''
    DEBUG	    详细信息，一般只在调试问题时使用。
    INFO	    证明事情按预期工作。
    WARNING	    某些没有预料到的事件的提示，或者在将来可能会出现的问题提示。例如：磁盘空间不足。但是软件还是会照常运行。
    ERROR	    由于更严重的问题，软件已不能执行一些功能了。
    CRITICAL    严重错误，表明软件已不能继续运行了。
    '''
    '''
    %(name)s                Logger的名字
    %(levelno)s             数字形式的日志级别
    %(levelname)s           本形式的日志级别
    %(pathname)s            调用日志输出函数的模块的完整路径名，可能没有
    %(filename)s            调用日志输出函数的模块的文件名
    %(module)s              调用日志输出函数的模块名
    %(funcName)s            调用日志输出函数的函数名
    %(lineno)d              调用日志输出函数的语句所在的代码行
    %(created)f             当前时间，用UNIX标准的表示时间的浮点数表示
    %(relativeCreated)d     输出日志信息时的，自Logger创建以来的毫秒数
    %(asctime)s             字符串形式的当前时间。默认格式是"2003-07-08 16:49:45,896"。逗号后面的是毫秒
    %(thread)d              线程ID。可能没有
    %(threadName)s          线程名。可能没有
    %(process)d             进程ID。可能没有
    %(message)s             用户输出的消息
    '''
    '''
    # 实例化TimedRotatingFileHandler
    # interval是时间间隔，backupCount是备份文件的个数，如果超过这个个数，就会自动删除，when是间隔的时间单位，单位有以下几种：
    # S 秒
    # M 分
    # H 小时
    # D 天
    # W 每星期（interval==0时代表星期一）
    # midnight 每天凌晨
    '''
