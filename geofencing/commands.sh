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
SETHOOK alert $API_ENDPOINT NEARBY intruder FENCE DETECT inside POINT 35.6581 139.6975 100
SET intruder point:X POINT 35.65949 139.69963
SET intruder point:X POINT 35.65946 139.69947
SET intruder point:X POINT 35.65944 139.69932
SET intruder point:X POINT 35.65936 139.69909
SET intruder point:X POINT 35.65932 139.69893
SET intruder point:X POINT 35.65929 139.69879
SET intruder point:X POINT 35.65923 139.69860
SET intruder point:X POINT 35.65917 139.69836
SET intruder point:X POINT 35.65914 139.69821
SET intruder point:X POINT 35.65908 139.69801
SET intruder point:X POINT 35.65899 139.69769
SET intruder point:X POINT 35.65891 139.69755
SET intruder point:X POINT 35.65889 139.69750
SET intruder point:X POINT 35.65882 139.69740
SET intruder point:X POINT 35.65873 139.69731
SET intruder point:X POINT 35.65868 139.69723
SET intruder point:X POINT 35.65860 139.69718
SET intruder point:X POINT 35.65849 139.69705
SET intruder point:X POINT 35.65841 139.69689
SET intruder point:X POINT 35.65828 139.69679
SET intruder point:X POINT 35.65817 139.69670
SET intruder point:X POINT 35.65809 139.69664
SET intruder point:X POINT 35.65794 139.69653
SET intruder point:X POINT 35.65778 139.69642
SET intruder point:X POINT 35.65764 139.69629
SET intruder point:X POINT 35.65753 139.69619
SET intruder point:X POINT 35.65745 139.69615
SET intruder point:X POINT 35.65733 139.69605
SET intruder point:X POINT 35.65725 139.69600
SET intruder point:X POINT 35.65716 139.69596
SET intruder point:X POINT 35.65703 139.69588
SET intruder point:X POINT 35.65691 139.69576
EOF
}

#commands-2 | tile38-cli
commands | while read line
do
	exec_tile38_cli "$line"
	sleep 0.5
done

