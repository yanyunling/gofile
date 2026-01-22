#!/bin/sh
rm -rf ./web/dist
cd ./web
npm install
npm run build
cd ../
go build
echo build finished
