# Test Client to localhost:9955 mbserver

import time
import logging
import sys
import mymemcached
import myclient as bmemcached

# logger = logging.getLogger()
# logger.setLevel(logging.DEBUG)
# ch = logging.StreamHandler(sys.stdout)
# ch.setLevel(logging.DEBUG)
# logger.addHandler(ch)

client = bmemcached.Client(('127.0.0.1:9955', ), '', '')
client.set('Blue Wonderful', 'Blue Wonderful to me')
client = bmemcached.Client(('127.0.0.1:9955', ), '', '')
print client.get('Blue Wonderful')

client = bmemcached.Client(('127.0.0.1:9955', ), '', '')
client.set('Looking Up', 5)
client = bmemcached.Client(('127.0.0.1:9955', ), '', '')
print client.get('Looking Up')

array = ['hi', 'hello', 'good morning']
client = bmemcached.Client(('127.0.0.1:9955', ), '', '')
client.set('Elton', array)
client = bmemcached.Client(('127.0.0.1:9955', ), '', '')
print client.get('Elton')



