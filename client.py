# Test Client to localhost:9955 mbserver

import bmemcached
import time
client = bmemcached.Client(('127.0.0.1:9955', ), '', '')
client.set('Blue Wonderful', 'Blue Wonderful to me')
client = bmemcached.Client(('127.0.0.1:9955', ), '', '')
print client.get('Blue Wonderful')

array = ['hi', 'hello', 'good morning']
client = bmemcached.Client(('127.0.0.1:9955', ), '', '')
client.set('Elton', array)
client = bmemcached.Client(('127.0.0.1:9955', ), '', '')
print client.get('Elton')


