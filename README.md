# mbserver


Instructions
-----------------
I've written it in go, so in order to run it, you'll need to have go downloaded on your machine. You can run the code with:

go run mbserver.go <IP>:<PORT>

You will need to set your GOPATH environment variable to your workspace.

Performance
----------------

This server handles each connected client in separate goroutines, 
allowing multiple readers and writers to try to access data at once.
The data is protected using a Read Write lock, which allows multiple
readers but blocking readers when a thread wants to write.

A potential deadlock that could arise is I put a blocking send call inside the lock. If a connected client dies and the server is unable to complete
the send, the server will block inside the lock and not free the lock. A better design would be to implement non-blocking sends inside the locked areas.

The set calls might be expectected to happen sequentially on the client side. However, they may arrive out of order on the server side. My server does not handle out-of-order requests and deals with requests first come, first serve.

A potential bottleneck is having many writers requesting the lock -- only one writer will be allowed to write at a time, and they will also lock out any readers. The read write lock puts a preference on writers, so a steady stream of writers will push out the many requests of readers.    

I used a third party for an external library, gomemcached. I used it to
encode and decode my packets as per the protocol. This meant I did not have to reinvent the wheel in sending and receiving packets, however I did need to test it carefully with the python library. Happily I found out they were compatible, however they had different levels of strictness. My go client happily accepted GET responses without flags, whereas the python client crashed when receiving GET responses without flags. After inspecting the source code, I was able to fix this error.

Testing
------------
mbclient.go and client.py are the two files I used for testing
