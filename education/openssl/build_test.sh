#!/bin/sh

# Remove all the old stuff
rm -rf ./ca 1>/dev/stdout 2>/dev/stderr

# Build the new stuff
echo "Creating folder structure for CA and seeding DB"
mkdir ca
mkdir ./ca/rootca/
mkdir ./ca/rootca/private
mkdir ./ca/rootca/public
mkdir ./ca/rootca/db
touch ./ca/rootca/db/rootca.db
touch ./ca/rootca/db/rootca.db.attr
echo "1000" > ./ca/rootca/db/rootca.crt.srl
echo "1000" > ./ca/rootca/db/rootca.crl.srl

echo "Creating intermediate CA folder structure"
mkdir ./ca/intca/
mkdir ./ca/intca/private
mkdir ./ca/intca/public
mkdir ./ca/intca/db
touch ./ca/intca/db/intca.db
touch ./ca/intca/db/intca.db.attr


# Make the root CA
echo "...Generating root CA (i.e. controlling entity), put in passcode when prompted" > /dev/stdout
openssl genrsa -aes256 -out rootca.key 4096
openssl req -config root_ca.conf -sha256 -new -x509 -days 3650 -key rootca.key -out rootca.crt

echo "...Generating Intermediate CA (i.e. delegated authority)" > /dev/stdout
openssl genrsa -out intca.key 2048
openssl req -new -config root_ca.conf -sha256 -key intca.key -out intca.csr

echo "...Moving key files and certs to correct locations"
mv ./rootca.key ./ca/rootca/private/
mv ./rootca.crt ./ca/
mv ./intca.key ./ca/intca/private/

echo "...Sign the request for the delegated authority" > /dev/stdout
openssl ca -batch -config root_ca.conf -notext -in intca.csr -out ./ca/intca.crt
openssl ca -config root_ca.conf -gencrl -keyfile ./ca/rootca/private/rootca.key -cert ./ca/rootca.crt -out rootca.crl.pem
