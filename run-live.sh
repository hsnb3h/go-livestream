#!/bin/bash

export DISPLAY=:1

xvfb-run --listen-tcp  --server-num 1  -s "-ac -screen 0 1504x846x24" /bin/live
