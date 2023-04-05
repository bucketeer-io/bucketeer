# localenv

## How to run 

To run the component as a part of the main Bucketeer chart, it could be invoked as (from the general Bucketeer chart folder):

        $ helm install bucketeer . -f values/account-apikey-cacher.yaml -f values/event-persister-evaluation-events-dwh.yaml	-f values/migration-mysql.yaml -f values/account.yaml	-f values/event-persister-evaluation-events-evaluation-count.yaml	-f values/notification-sender.yaml -f values/api-gateway.yaml	-f values/event-persister-evaluation-events-ops.yaml	-f values/notification.yaml -f values/auditlog-persister.yaml	-f values/event-persister-goal-events-dwh.yaml -f values/ops-event-batch.yaml -f values/auditlog.yaml -f values/event-persister-goal-events-ops.yaml -f values/push-sender.yaml -f values/auth.yaml	-f values/experiment.yaml	-f values/push.yaml -f values/auto-ops.yaml	-f values/feature-recorder.yaml -f values/user-persister.yaml -f values/calculator.yaml	-f values/feature-segment-persister.yaml -f values/user.yaml -f values/feature-tag-cacher.yaml -f values/web-gateway.yaml -f values/dex.yaml	-f values/feature.yaml -f values/web.yaml -f values/environment.yaml -f values/global.yaml -f values/event-counter.yaml -f values/metrics-event-persister.yaml -f charts/localenv/values.yaml

## How to update dependencies

To update all dependencies (optional) the following command could be run:

        $ helm dependency update .

## Bucketeer adjustments

1. In all values.yaml files the following parts of code should be commented (otherwise no pods will be deployed in Minikube):

        - matchExpressions:
            - key: cloud.google.com/gke-nodepool

1. Setup a secret with the following name (or choose another, but then it should be mentioned in values.yaml):

        bucketeer-jp-cert-20220411