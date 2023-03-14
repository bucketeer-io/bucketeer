####################################
# Generate Certificates
####################################

TLS_TARGETS := generate-tls-certificate-local
OAUTH_TARGETS := generate-oauth-local
SERVICE_TOKEN_TARGETS := generate-service-token-local

.PHONY: $(TLS_TARGETS)
$(TLS_TARGETS): generate-tls-certificate-%:
	mkdir -p ${CURDIR}/$*
	openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
		-keyout ${CURDIR}/$*/tls.key \
		-out ${CURDIR}/$*/tls.crt \
		-subj "/CN=*.default.svc.cluster.local" \
		-config ${CURDIR}/service.config

.PHONY: $(OAUTH_TARGETS)
$(OAUTH_TARGETS): generate-oauth-%:
	mkdir -p ${CURDIR}/$*
	openssl genrsa -out ${CURDIR}/$*/oauth-private.pem 4096
	openssl rsa -in ${CURDIR}/$*/oauth-private.pem -pubout -out ${CURDIR}/$*/oauth-public.pem

.PHONY: $(SERVICE_TOKEN_TARGETS)
$(SERVICE_TOKEN_TARGETS): generate-service-token-%:
	./generate-service-token.sh "$*"

###############################
# Create or Update k8s secrets
###############################

SERVICE_CERT_SECRET_TARGETS := service-cert-secret-local
OAUTH_KEY_SECRET_TARGETS := oauth-key-secret-local
SERVICE_TOKEN_SECRET_TARGETS := service-token-secret-local

.PHONY: $(SERVICE_CERT_SECRET_TARGETS)
$(SERVICE_CERT_SECRET_TARGETS): service-cert-secret-%:
	kubectl delete secret bucketeer-service-cert --ignore-not-found
	kubectl create secret tls bucketeer-service-cert --key ${CURDIR}/$*/tls.key --cert ${CURDIR}/$*/tls.crt

.PHONY: $(OAUTH_KEY_SECRET_TARGETS)
$(OAUTH_KEY_SECRET_TARGETS): oauth-key-secret-%:
	kubectl delete secret bucketeer-oauth-key --ignore-not-found
	kubectl create secret generic bucketeer-oauth-key --from-file public.pem=${CURDIR}/$*/oauth-public.pem --from-file private.pem=${CURDIR}/$*/oauth-private.pem

.PHONY: $(SERVICE_TOKEN_SECRET_TARGETS)
$(SERVICE_TOKEN_SECRET_TARGETS): service-token-secret-%:
	kubectl delete secret bucketeer-service-token --ignore-not-found
	kubectl create secret generic bucketeer-service-token --from-file token=${CURDIR}/$*/service.token

# For local environment, we probably won't need this, but we might need to set the localhost ssl certificate
# to this secret so it won't fail while deploying
# This secret bucketeer-jp-cert-${DATE} must match the secretName set in the api-gateway.yaml and web-gateway.yaml
BUCKETEER_JP_KEY := bucketeer-jp.key
BUCKETEER_JP_CERT := bucketeer-jp.crt

.PHONY: bucketeer-jp-certs
bucketeer-jp-certs:
	$(eval DATE=${shell date +"%Y%m%d"})
	kubectl create secret tls bucketeer-jp-cert-${DATE} --key ${CURDIR}/$(BUCKETEER_JP_KEY) --cert ${CURDIR}/$(BUCKETEER_JP_CERT)
