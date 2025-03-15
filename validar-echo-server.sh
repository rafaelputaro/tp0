#!/bin/bash
OK_MESSAGE="action: test_echo_server | result: success"
FAIL_MESSAGE="action: test_echo_server | result: fail"
ERROR_MESSAGE="ERROR ON VALIDATE"
# run netcat-cli and run script validar
docker compose -f docker-compose-validar-dev.yaml up -d --build > /dev/null 2> /dev/null
# get netcat-cli log
LOG=$(docker compose -f docker-compose-validar-dev.yaml logs -f)
# remove netcat-cli container
docker rm container netcat-cli > /dev/null 2> /dev/null
# explore log
if echo "$LOG" | grep -q "$OK_MESSAGE"; then
  echo "$OK_MESSAGE"
elif echo "$LOG" | grep -q "$FAIL_MESSAGE"; then
  echo "$FAIL_MESSAGE"
else
  echo "$ERROR_MESSAGE"
fi