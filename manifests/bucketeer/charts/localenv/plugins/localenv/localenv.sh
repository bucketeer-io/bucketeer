helm dependencies update
helm install localenv .
minikube addons enable ingress &
minikube tunnel
