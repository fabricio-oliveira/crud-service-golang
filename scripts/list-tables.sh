#!/usr/bin/env sh
AWS_ACCESS_KEY_ID=DUMMY AWS_SECRET_ACCESS_KEY=DUMMY aws dynamodb list-tables --endpoint-url http://localhost:8000 --region us-est-1