#!/bin/bash
set -eu

source upload.env

rm -f ./bootstrap ./upload.zip
go build -o ./bootstrap
zip upload.zip bootstrap
aws --output=text --no-paginate --region=${REGION} lambda update-function-code \
    --function-name "${ARN}" \
    --zip-file "fileb://upload.zip"
rm -f ./bootstrap ./upload.zip
aws --output=text --no-paginate --region=${REGION} lambda wait function-updated-v2 \
    --function-name "${ARN}"

