#!/bin/bash

PROMPT="127.0.0.1:9851> "
API_ENDPOINT="http://localhost:8001/webhook"

exec_tile38_cli(){
	echo -n "$PROMPT"
	echo "$@"
	echo "$@" | tile38-cli
}
export -f exec_tile38_cli

commands(){
        cat <<EOF | awk 'gsub(/#.*/,"")>=0'
DROP intruder
#
SET intruder point:X POINT 35.6578 139.6971
SETHOOK alert $API_ENDPOINT NEARBY intruder FENCE DETECT inside POINT 35.6580 139.6970 100
SET intruder point:X POINT 35.6581 139.6968
EOF
}

#commands-2 | tile38-cli
commands | while read line
do
	exec_tile38_cli "$line"
done

