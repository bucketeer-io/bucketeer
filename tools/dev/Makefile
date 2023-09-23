NAMESPACE := default

generate-tls-certificate:
	mkdir -p $(CURDIR)/cert
	openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
		-keyout $(CURDIR)/cert/tls.key \
		-out $(CURDIR)/cert/tls.crt \
		-subj "/CN=*.bucketeer.org" \
		-config service.config

generate-oauth:
	mkdir -p ${CURDIR}/$*
	openssl genrsa -out ${CURDIR}/cert/oauth-private.pem 4096
	openssl rsa -in ${CURDIR}/cert/oauth-private.pem -pubout -out ${CURDIR}/cert/oauth-public.pem

service-cert-secret:
	kubectl delete secret bucketeer-service-cert --namespace $(NAMESPACE) --ignore-not-found
	kubectl create secret tls bucketeer-service-cert --key ${CURDIR}/cert/tls.key --cert ${CURDIR}/cert/tls.crt --namespace $(NAMESPACE)

oauth-key-secret:
	kubectl delete secret bucketeer-oauth-key --namespace $(NAMESPACE) --ignore-not-found
	kubectl create secret generic bucketeer-oauth-key --from-file public.pem=${CURDIR}/cert/oauth-public.pem --from-file private.pem=${CURDIR}/cert/oauth-private.pem --namespace $(NAMESPACE)

generate-github-token:
	kubectl delete secret github-token --namespace ${NAMESPACE} --ignore-not-found
	kubectl create secret generic github-token --from-literal bucketeer-bot-access-token=${GITHUB_TOKEN}

setup-minikube:
	minikube start --memory 6g
	minikube addons enable ingress
	kubectl patch deployment -n ingress-nginx ingress-nginx-controller --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value":"--enable-ssl-passthrough"}]'
	helm install localenv ../../manifests/localenv

start-minikube: 
	minikube start
	# used to enable ssl-passthrough for nginx ingress container
	kubectl patch deployment -n ingress-nginx ingress-nginx-controller --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value":"--enable-ssl-passthrough"}]'