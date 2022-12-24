#!/bin/bash

openssl genrsa -out ca.key 2048 # 生成ca私钥
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt -subj "/CN=mingshen.com" #自签名证书

openssl genrsa -out server.key 2048
#生成服务端证书申请文件
openssl req -new -key server.key -out server.csr -subj "/CN=myserver.com" -reqexts SAN \
-config <(cat /usr/local/openssl/ssl/openssl.cnf <(printf "\n[SAN]\nsubjectAltName=DNS:myserver.com"))

#CA签名的服务端证书
openssl x509 -req -days 3650 -in server.csr -out server.crt -CA ca.crt -CAkey ca.key -CAcreateserial \
-extensions SAN -extfile <(cat /usr/local/openssl/ssl/openssl.cnf <(printf "\n[SAN]\nsubjectAltName=DNS:myserver.com"))


openssl genrsa -out client.key 2048
#生成客户端证书申请文件
openssl req -new -key client.key -out client.csr -subj "/CN=myclient.com" -reqexts SAN \
-config <(cat /usr/local/openssl/ssl/openssl.cnf <(printf "\n[SAN]\nsubjectAltName=DNS:myserver.com"))
#CA签名的客户端证书
openssl x509 -req -days 3650 -in client.csr -out client.crt -CA ca.crt -CAkey ca.key -CAcreateserial \
-extensions SAN -extfile <(cat /usr/local/openssl/ssl/openssl.cnf <(printf "\n[SAN]\nsubjectAltName=DNS:myserver.com"))

