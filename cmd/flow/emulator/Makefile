# The short Git commit hash
SHORT_COMMIT := $(shell git rev-parse --short HEAD)
# The Git commit hash
COMMIT := $(shell git rev-parse HEAD)
# The tag of the current commit, otherwise empty
VERSION := $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
IMAGE_TAG := ${VERSION}
ifeq (${IMAGE_TAG},)
IMAGE_TAG := ${SHORT_COMMIT}
endif
# Service information
#----------------------------------------------------------------------
SERVICE_SHORT_NAME=emulator
SERVICE_LONG_NAME=dl-flow/${SERVICE_SHORT_NAME}
DEPLOYMENT_NAME=flow-emulator-v1

# Build information
#----------------------------------------------------------------------
REPO:=gcr.io
IMAGE_URL:=${REPO}/${SERVICE_LONG_NAME}

MONOREPO_ROOT=`pwd`/../../../
# The location of the k8s YAML files
K8S_YAMLS_LOCATION=./k8s

KUBECONFIG := $(shell uuidgen)

#----------------------------------------------------------------------
# CD COMMANDS
#----------------------------------------------------------------------

.PHONY: deploy-staging
deploy-staging: update-deployment-image apply-staging-files monitor-rollout

.PHONY: deploy-production
deploy-production: update-deployment-image apply-prod-files monitor-rollout

# Staging YAMLs must have 'staging' in their name.
.PHONY: apply-staging-files
apply-staging-files:
	echo "$$KUBECONFIG_STAGING_2" > ${KUBECONFIG}; \
	files=$$(find ${K8S_YAMLS_LOCATION} -type f \( -name "*.yml" -or -name "*.yaml" \) | grep staging); \
	echo "$$files" | xargs -I {} kubectl --kubeconfig=${KUBECONFIG} apply -f {}

# Production YAMLs must have 'production' in their name.
.PHONY: apply-prod-files
apply-prod-files:
	echo "$$KUBECONFIG_PRODUCTION_2" > ${KUBECONFIG}; \
	files=$$(find ${K8S_YAMLS_LOCATION} -type f \( -name "*.yml" -or -name "*.yaml" \) | grep production); \
	echo "$$files" | xargs -I {} kubectl --kubeconfig=${KUBECONFIG} apply -f {}

# Deployment YAMLs must have 'deployment' in their name.
.PHONY: update-deployment-image
update-deployment-image: CONTAINER=flow-emulator
update-deployment-image:
	@files=$$(find ${K8S_YAMLS_LOCATION} -type f \( -name "*.yml" -or -name "*.yaml" \) | grep deployment); \
	for i in $$files; do \
		patched=`openssl rand -hex 8`; \
		kubectl patch -f $$i -p '{"spec":{"template":{"spec":{"containers":[{"name":"${CONTAINER}","image":"${IMAGE_URL}:${IMAGE_TAG}"}]}}}}' --local -o yaml > $$patched; \
		mv -f $$patched $$i; \
	done

# Has to come after one of the `apply` targets
.PHONY: monitor-rollout
monitor-rollout:
	kubectl --kubeconfig=${KUBECONFIG} rollout status statefulsets.apps ${DEPLOYMENT_NAME}
