#!/bin/bash

commands(){
        cat <<EOF | awk 'gsub(/#.*/,"")>=0'
DROP location
DROP lunch
#
GET location me POINT 35.6581 139.6975
# 
SET lunch ramen:A POINT 35.6586 139.6982
SET lunch gyudon:B POINT 35.6570 139.6967
SET lunch ramen:C POINT 35.6576 139.6948
SCAN lunch MATCH ramen* IDS
NEARBY lunch POINT 35.6581 139.6975 100
EOF
}

commands | tile38-cli
