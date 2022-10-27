#!/bin/bash

APP=$1

case $APP in
    "autoops") echo "auto-ops" ;;
    "eventcounter") echo "event-counter" ;;
    "eventpersister") echo "event-persister" ;;
    "goalbatch") echo "goal-batch" ;;
    "metricsevent") echo "metrics-event" ;;
    "opsevent") echo "ops-event" ;;
    *) echo $APP
esac
