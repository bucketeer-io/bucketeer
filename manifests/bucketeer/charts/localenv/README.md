# localenv

## How does it work?

The chart of localenv is stored under charts/localenv directory. It included the following dependencies:
- bq
- pubsub
- hashicorp vault
- redis
- mysql
- nginx ingress controller (currently implemented via minikube addon)

Because chart uses the ingress controller for local usage on minikube, it uses minikube ingress addon wrapped into helm plugin.

## How to run

From the root folder of the localenv chart, call the following commands:

                $ minikube start
                $ helm plugin install plugins/localenv
                $ helm localenv

The plugin call (helm localenv) will do the following operations automatically:

- install the chart with all dependencies
- start ingress addon
- start local tunnel so that kubernetes cluster loadbalancer will be available in local browser as localhost.