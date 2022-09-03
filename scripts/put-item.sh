#!/usr/bin/env sh
AWS_ACCESS_KEY_ID=DUMMY AWS_SECRET_ACCESS_KEY=DUMMY \
aws dynamodb put-item \
    --table-name Invoice \
    --item '{"Id": {"S": "123"}, "BillTo": {"S": "pay store"}, "CreatedAt": {"S": "2022-09-02 16:46:12.667233 -0300 -03 m=+44.419170897"}, "UpdatedAt": {"S": "2022-09-02 16:46:12.667233 -0300 -03 m=+44.419170897"}}' \
    --condition-expression "attribute_not_exists(Id)" \
    --endpoint-url http://localhost:8000 \
    --region us-est-1