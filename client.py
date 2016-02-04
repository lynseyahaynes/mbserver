import bmemcached
import time
client = bmemcached.Client(('127.0.0.1:9955', ), '', '')
client.set("this_key_has_taken_its_toll_on_me", "i_cant_say_goodbye_anymore")
client = bmemcached.Client(('127.0.0.1:9955', ), '', '')
print client.get("this_key_has_taken_its_toll_on_me")

