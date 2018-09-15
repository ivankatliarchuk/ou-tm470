#!/bin/bash

git add ../; git commit -m "Auto Sync and Run"; git push
#go get github.com/mmd93ee/ou-tm470/dataPersist
#go get github.com/mmd93ee/ou-tm470/webServer

go build serviceChain.go; ./serviceChain
