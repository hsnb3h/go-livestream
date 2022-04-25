#!/bin/bash

export DISPLAY=:44
pulseaudio --fail -D --exit-idle-time=-1
pacmd load-module module-virtual-sink sink_name=v1
pacmd set-default-sink v1
pacmd set-default-source v1.monitor
xvfb-run --listen-tcp  --server-num 44  -s "-ac -screen 0 1504x846x24" /bin/live
