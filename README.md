gotunnel
-------

`gotunnel` is a simple proof of concept using the go.crypto ssh package to tunnel HTTP PROXY traffic over SSH.

`gotunnel` also support reverse proxy through SSH tunnel.

Fight the man, eat dogfood, amaze your friends.

Usage
=====

     ./gotunnel -host=123.123.123.123 -pass=123456 -local_addr=127.0.0.1:8080 -user=root

Then configure your web browser to use a HTTP proxy on localhost:8888
