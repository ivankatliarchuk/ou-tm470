#!/bin/sh

clear

# Remove all the old stuff
rm -rf ./ca 1>/dev/stdout 2>/dev/stderr

# Build the new stuff
echo "Creating folder structure for CA and seeding DB"
mkdir ca
mkdir ./ca/rootca/
mkdir ./ca/rootca/crt
mkdir ./ca/rootca/crl
mkdir ./ca/rootca/private
mkdir ./ca/rootca/db
touch ./ca/rootca/db/rootca.db
touch ./ca/rootca/db/rootca.db.attr
echo "1000" > ./ca/rootca/db/rootca.crt.srl
echo "1000" > ./ca/rootca/db/rootca.crl.srl

echo "Creating intermediate CA folder structure"
mkdir ./ca/intca/
mkdir ./ca/intca/crt
mkdir ./ca/intca/crl
mkdir ./ca/intca/private
mkdir ./ca/intca/db
touch ./ca/intca/db/intca.db
touch ./ca/intca/db/intca.db.attr


# Make the root CA
echo "...Generating root CA key and certificate (i.e. validating entity)" > /dev/stdout
openssl genrsa -aes256 -out ./ca/rootca/private/rootca.key 4096
openssl req -config rootca.conf -sha256 -new -x509 -days 3650 -key ./ca/rootca/private/rootca.key -out ./ca/rootca/crt/rootca.crt

echo "...Generating Intermediate CA key and certificate request (i.e. delegated authority)" > /dev/stdout
openssl genrsa -out ./ca/intca/private/intca.key 2048
openssl req -new -config intca.conf -sha256 -key ./ca/intca/private/intca.key -out ./ca/intca/intca.csr

echo "...Sign the request for the delegated authority" > /dev/stdout
openssl ca -batch -config rootca.conf -notext -in ./ca/intca/intca.csr -out ./ca/intca/crt/intca.crt
openssl ca -config rootca.conf -gencrl -keyfile ./ca/rootca/private/rootca.key -cert ./ca/rootca/crt/rootca.crt -out ./ca/rootca/crl/rootca.crl.pem

echo "...Text output for the certificates" > /dev/stdout
openssl x509 -noout -text -in ./ca/rootca/crt/rootca.crt > rootca.about.txt
openssl x509 -noout -text -in ./ca/intca/crt/intca.crt > intca.about.txt

echo "...Create a certificate chain holding Root and Int. CA"
cat ./ca/intca/crt/intca.crt ./ca/rootca/crt/rootca.crt > ./ca/ca.chain.crt
