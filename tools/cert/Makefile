
BUCKETEER_MANIFESTS := $(wildcard ../../manifests/bucketeer/charts/*)

TLS_TARGETS := $(addprefix generate-tls-certificate-, $(notdir $(BUCKETEER_MANIFESTS)))
OAUTH_TARGETS := $(addprefix generate-oauth-, $(notdir $(BUCKETEER_MANIFESTS)))

NAMESPACE := default

SERVICE_CERT_SECRET_TARGETS := $(addprefix service-cert-secret-, $(notdir $(BUCKETEER_MANIFESTS)))
OAUTH_KEY_SECRET_TARGETS := $(addprefix oauth-key-secret-, $(notdir $(BUCKETEER_MANIFESTS)))


.PHONY: $(TLS_TARGETS)
$(TLS_TARGETS): generate-tls-certificate-%:
	mkdir -p $(CURDIR)/$*
	openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
		-keyout $(CURDIR)/$*/tls.key \
		-out $(CURDIR)/$*/tls.crt \
		-subj "/CN=*.default.svc.cluster.local"

.PHONY: $(OAUTH_TARGETS)
$(OAUTH_TARGETS): generate-oauth-%:
	mkdir -p ${CURDIR}/$*
	openssl genrsa -out ${CURDIR}/$*/oauth-private.pem 4096
	openssl rsa -in ${CURDIR}/$*/oauth-private.pem -pubout -out ${CURDIR}/$*/oauth-public.pem


.PHONY: $(SERVICE_CERT_SECRET_TARGETS)
$(SERVICE_CERT_SECRET_TARGETS): service-cert-secret-%:
	kubectl delete secret bucketeer-service-cert --namespace $(NAMESPACE) --ignore-not-found
	kubectl create secret tls bucketeer-service-cert --key ${CURDIR}/$*/tls.key --cert ${CURDIR}/$*/tls.crt --namespace $(NAMESPACE)

.PHONY: $(OAUTH_KEY_SECRET_TARGETS)
$(OAUTH_KEY_SECRET_TARGETS): oauth-key-secret-%:
	kubectl delete secret bucketeer-oauth-key --namespace $(NAMESPACE) --ignore-not-found
	kubectl create secret generic bucketeer-oauth-key --from-file public.pem=${CURDIR}/$*/oauth-public.pem --from-file private.pem=${CURDIR}/$*/oauth-private.pem --namespace $(NAMESPACE)