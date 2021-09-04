# -*- coding=utf-8 -*-
import os
from app_fiask import init_app
from router import reg_blueprint

BASE_DIR = os.path.dirname(os.path.abspath(__file__))

app = init_app()
reg_blueprint(app)

if __name__ == '__main__':
    app.debug = True
    app.run()
