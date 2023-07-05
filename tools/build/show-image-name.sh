#!/bin/bash

APP=$1

case $APP in
    "autoops") echo "auto-ops" ;;
    "eventcounter") echo "event-counter" ;;
    "eventpersister") echo "event-persister" ;;
    "eventpersisterdwh") echo "event-persister-dwh" ;;
    "eventpersisterops") echo "event-persister-ops" ;;
    "metricsevent") echo "metrics-event" ;;
    "opsevent") echo "ops-event" ;;
    "experimentcalculator") echo "experiment-calculator" ;;
    *) echo $APP
esac
