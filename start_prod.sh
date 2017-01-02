#!/usr/bin/env bash

ps aux|grep -v grep|grep ana_prod|awk '{print $2}'|xargs kill -9
nohup ./bin/ana_prod --conf=/home/fe/ana_prod_0.2/conf/conf.prod.toml >/dev/null 2>&1 &

