import os

from settings import LOG_DIR, LOG_FILE_NAME

'''
在项目目录里执行命令启动应用
gunicorn -c gunicorn_config.py main:app
'''

# 获取工作目录
BASE_DIR = os.path.abspath(os.path.dirname(__file__))
RUN_DIR = os.path.join(BASE_DIR, "run")

if not os.path.exists(LOG_DIR):
    os.mkdir(LOG_DIR)

if not os.path.exists(RUN_DIR):
    os.mkdir(RUN_DIR)

bind = '0.0.0.0:8080'
# 超过这么多秒后工作将被杀掉，并重新启动。一般设定为30秒
timeout = 120
# 代码更新时将重启工作，默认为False。此设置用于开发，每当应用程序发生更改时，都会导致工作重新启动。
reload = True

# 设置守护进程
daemon = False

# 进程数
workers = 1

# gunicorn要切换到的目的工作目录
chdir = BASE_DIR
# pid文件
pidfile = os.path.join(RUN_DIR, "ledger.pid")

# 日志级别，这个日志级别指的是错误日志的级别，而访问日志的级别无法设置
loglevel = 'info'
access_log_format = '%(t)s %(p)s %(h)s "%(r)s" %(s)s %(L)s %(b)s %(f)s" "%(a)s"'

# 日志
accesslog = os.path.join(LOG_DIR, "access.log")
errorlog = os.path.join(LOG_DIR, LOG_FILE_NAME)
