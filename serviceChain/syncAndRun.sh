#!/bin/bash

git add ../; git commit -m "Auto Sync and Run"; git push
go get https://www.github.com/mmd93ee/ou-tm470/dataPersist
go get https://www.github.com/mmd93ee/ou-tm470/webServer
go build serviceChain.go; ./serviceChain
