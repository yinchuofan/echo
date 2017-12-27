# -*- coding: utf-8 -*-

import json
import os
import unittest
import md5
import hashlib
import urllib
import time
from datetime import *

import requests

host_api = 'http://localhost:9000'

session = requests.session()

class PushMessageTest(unittest.TestCase):

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_push_message_alarm(self):
        print "\n----------test_push_message_alarm----------"

        url = host_api + '/publish'

        message = dict()
        message['code'] = 'MsgCodeAlarm'
        message['content'] = 'hello world'
        message['publisher'] = 'echo'
        message['time'] = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        message_str = json.dumps(message)

        channel = 'echo_alarm'
        body = 'channel:' + channel + '\r\n' + 'message:' + message_str

        print 'url: \n', url
        print 'body: \n', body

        response = session.post(url, data=body)
        print(response.text)

if __name__ == '__main__':
    unittest.main()
