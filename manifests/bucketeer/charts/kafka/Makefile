VERSION := $(shell helm dep list | grep strimzi-kafka-operator | awk '{ print $$2 }')

.PHONY: upgrade-crd
upgrade-crd: .unzip-charts
	kubectl replace -f .charts/strimzi-kafka-operator/crds

.PHONY: .unzip-charts
.unzip-charts:
	tar xvf charts/strimzi-kafka-operator-helm-3-chart-${VERSION}.tgz -C charts/
