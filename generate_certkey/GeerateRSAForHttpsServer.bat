cd /d %~dp0
set OPENSSL_CONF=openssl.cnf
openssl genrsa -out server.key 1024
openssl req -new -key server.key -out server.csr
openssl rsa -in server.key -out server_nopwd.key
openssl x509 -req -days 365 -in server.csr -signkey server_nopwd.key -out server.crt


