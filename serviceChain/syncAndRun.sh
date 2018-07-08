#!/bin/bash

git add ../; git commit -m "Debugging persistence"; git push; go get github.com/mmd93ee/ou-tm470/dataPersist
go build serviceChain.go; ./serviceChain
