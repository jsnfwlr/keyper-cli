#!/bin/sh
USER="$1"
FP="$2"
HOST=`hostname -s`
KEYPER_HOST=keyper-cli-test-keyper
CURL_ARGS="-s -q -f -m 7"
CURL_ARGS="${CURL_ARGS} --data-urlencode username=${USER}"
CURL_ARGS="${CURL_ARGS} --data-urlencode host=${HOST}"

[ -z ${FP} ] || CURL_ARGS="${CURL_ARGS} --data-urlencode fingerprint=${FP}"

## Use this if you want to get public keys using HTTP GET
# curl -G ${CURL_ARGS} https://${KEYPER_HOST}/api/authkeys

## Use this if you want get public keys using HTTP POST
curl ${CURL_ARGS} http://${KEYPER_HOST}/api/authkeys

## Ensure a new line is added to the end of the output
EXITCODE=$?
echo ""

exit $EXITCODE