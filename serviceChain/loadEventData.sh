#!/bin/#!/usr/bin/env bash

curl -X POST --data @./md5589_service_test_data.json -H "Content-Type: application/json" http://localhost:8000/serviceevent/add/
