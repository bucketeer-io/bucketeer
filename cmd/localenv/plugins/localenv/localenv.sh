export HELM_EXPERIMENTAL_OCI=1
helm dependencies update
helm install localenv .
minikube addons enable ingress &
minikube tunnel