# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# Use bash explicitly in this Makefile to avoid unexpected platform
# incompatibilities among Linux distros.
#
SHELL := /bin/bash

MAKE := make --no-print-directory

#
# Allows for resources to be loaded from outside the root location of
# the kustomize config file. Ensures that resource don't need to be
# copied around the file system.
#
# See https://kubectl.docs.kubernetes.io/faq/kustomize
#
KOPTIONS := --load-restrictor LoadRestrictionsNone

#
# Include the main camel-k Makefile containing
# basic common recipes like kustomize and vars
# like VERSION
#
include script/Makefile

#
# Vars that can be overridden by external env vars
#
DRY_RUN ?= false
NAMESPACE ?= camel-k

# Global: [true|false]
# - On setup: will promote roles and bindings to cluster-level
# - On operator: set namespace to all using WATCH_NAMESPACE env var
GLOBAL ?= false
# Always Pull Images: [true|false]
ALWAYS_PULL_IMAGES ?= false
# Monitoring: [true|false]
# - On operator: will add the prometheus resources to install
MONITORING ?= false
# Monitoring Port: integer
MONITORING_PORT ?= 8080
# Health Port: integer
HEALTH_PORT ?= 8081
# Operator Logging Level: string [info, debug, 0, 1]
LOGGING_LEVEL ?= info
# Install default kamelets: [true|false]
INSTALL_DEFAULT_KAMELETS ?= true
# Install default pvc: [true|false]
INSTALL_DEFAULT_PVC ?= true
# Uninstall all Camel K resources: [true|false]
UNINSTALL_ALL ?=false

CONFIG := config
MANAGER := $(CONFIG)/manager
SAMPLES := $(CONFIG)/samples
RBAC := $(CONFIG)/rbac
RBAC_OS := $(RBAC)/openshift
RBAC_GLOBAL := global
OPERATOR := operator
PLACEHOLDER := placeholder
YAML := yaml

# Fetch the latest image name - may override the original constant
IMAGE_NAME := $(shell grep image: $(MANAGER)/operator-deployment.yaml | sed 's/.*image: \(.*\):.*/\1/')

# Setup patches
ROLE_TO_CROLE_PATCH := $(RBAC)/patch-role-to-clusterrole
ROLEBIN_TO_CROLEBIN_PATCH := $(RBAC)/patch-rolebinding-to-clusterrolebinding
# Operator patches
PORTS_PATCH := patch-ports
LOG_LEVEL_PATCH := patch-log-level
INSTALL_DEFAULT_KAMELETS_PATCH := patch-install-default-kamelets
INSTALL_DEFAULT_PVC_PATCH := patch-install-default-pvc
IMAGE_PULL_POLICY_PATCH := patch-image-pull-policy-always
WATCH_NAMESPACE_PATCH := patch-watch-namespace-global
# Platform patches
INT_PLATFORM_PATCH := patch-integration-platform

#
# Macro for editing kustomization to define
# the image reference
#
# Parameter: directory of the kustomization.yaml
#
define set-kustomize-image
	$(if $(filter $(IMAGE_NAME),$(CUSTOM_IMAGE):$(CUSTOM_VERSION)),,\
		@cd $(1) || exit 1 && \
			$(KUSTOMIZE) edit set image $(IMAGE_NAME)=$(CUSTOM_IMAGE):$(CUSTOM_VERSION))
endef

#
# Macro for editing kustomization to define
# the namespace
#
# Parameter: directory of the kustomization.yaml
#
define set-kustomize-namespace
	@cd $(1) || exit 1 && \
		$(KUSTOMIZE) edit set namespace $(NAMESPACE)
endef

#
# Add or remove a patch on a kustomization.yaml
# targetting a kind of resource
#
# Parameters:
# * directory of the kustomization.yaml
# * [add, remove]
# * path of patch
# * kind of resources, eg. Deployment, Role
#
define add-remove-kind-patch
	@cd $(1) || exit 1 && \
		$(KUSTOMIZE) edit $(2) patch --path $(3) --kind $(4) &> /dev/null
endef

#
# Macro for adding / removing the prometheus resources for monitoring
#
define add-remove-operator-monitoring
	@cd $(1) || exit 1 && \
		$(KUSTOMIZE) edit $(2) resource ../$(CONFIG)/prometheus &> /dev/null
endef

.PHONY: have-platform check_admin setup-cluster .setup-kubernetes .setup-openshift setup

#
# Determine the platform of the cluster based on
# either the use of querying through a go-client
# or using an installed client, ie. oc or kubectl
#
find-platform:
ifndef PLATFORM
PLATFORM=$(shell script/check_platform.sh)
endif

#
# Checks if the cluster platform has been defined correctly either by the user
# or by the platform_check script.
#
have-platform: find-platform
ifeq ($(PLATFORM),openshift)
	@echo Platform identified as 'openshift'
else ifeq ($(PLATFORM),kubernetes)
	@echo Platform identified as 'kubernetes'
else
	@echo "****"
	@echo "**** ERROR: Cannot continue as cluster platform cannot be identified ****"
	@echo "****"
	@exit 1
endif

#
# Checks if the cluster user has the necessary privileges to be a cluster-admin
# In this case if the user can list the CRDs then probably a cluster-admin
#
check-admin: kubectl
	@output=$$(kubectl get crd 2>&1) || (echo "****" && echo "**** ERROR: Cannot continue as user is not a Cluster-Admin ****" && echo "****"; exit 1)

crd-api-support: kubectl
ifndef CRD_SUPPORT
CRD_SUPPORT=$(shell script/check_crd_api_support.sh)
endif

check-crd-api-support: crd-api-support
ifneq ($(CRD_SUPPORT),OK)
	$(error *** CRD API FAILURE: $(CRD_SUPPORT) ****)
endif

#---
#
#@ setup-cluster
#
#== Setup the cluster installation by installing crds and cluster roles.
#
# Will either call setup-cluster-openshift (then setup-cluster-kubernetes) or
# setup-cluster-kubernetes depending on the identity of the cluster
#
#=== Cluster-admin privileges are required.
#
#* PARAMETERS:
#** NAMESPACE: Set the namespace for the resources
#** PLATFORM:  Override the discovered platform, if required
#** DRY_RUN:   Print the resources to be applied instead of applying them [true|false]
#
#---
setup-cluster: check-admin check-crd-api-support have-platform kustomize kubectl
# Set the namespace in the setup-cluster kustomization yaml
	@$(call set-kustomize-namespace,$@)
ifeq ($(PLATFORM), openshift)
	@for res in $(RBAC_OS)/operator-cluster*; do \
		cd $@ || exit 1 && \
			$(KUSTOMIZE) edit add resource ../$$res && \
			cd - &> /dev/null; \
	done
endif
#
# Build the resources
# Post-process ClusterRoleBindings to fix the namespace in the refs (not yet handled by kustomize)
# Either apply to the cluster or output to CLI
#
ifeq ($(DRY_RUN), false)
	@$(KUSTOMIZE) build $(KOPTIONS) $@ | \
		sed 's/$(PLACEHOLDER)/$(NAMESPACE)/' | \
		kubectl apply --server-side -f -
else
	@$(KUSTOMIZE) build $(KOPTIONS) $@ | \
		sed 's/$(PLACEHOLDER)/$(NAMESPACE)/'
endif

#---
#
#@ setup
#
#== Setup the installation by installing roles and granting privileges for the installing operator.
#
#=== Calls setup-cluster
#=== Cluster-admin privileges are required.
#
#* PARAMETERS:
#** NAMESPACE: Set the namespace for the resources
#** GLOBAL:    Convert all roles & bindings to cluster-level [true|false]
#** PLATFORM:  Override the discovered platform, if required
#** DRY_RUN:   Print the resources to be applied instead of applying them [true|false]
#
#---
setup: setup-cluster
# Set the namespace in the setup kustomization yaml
	@$(call set-kustomize-namespace,$@)
# If GLOBAL then add the conversion patches for all roles and rolebindings
ifeq ($(GLOBAL),true)
	@$(call add-remove-kind-patch,setup,add,../$(ROLE_TO_CROLE_PATCH).$(YAML),Role)
	@$(call add-remove-kind-patch,setup,add,../$(ROLEBIN_TO_CROLEBIN_PATCH).$(YAML),RoleBinding)
else
	@$(call add-remove-kind-patch,setup,remove,../$(ROLE_TO_CROLE_PATCH).$(YAML),Role)
	@$(call add-remove-kind-patch,setup,remove,../$(ROLEBIN_TO_CROLEBIN_PATCH).$(YAML),RoleBinding)
endif
ifeq ($(PLATFORM), openshift)
	@for res in $(RBAC_OS)/operator-role*; do \
		cd $@ || exit 1 && \
			$(KUSTOMIZE) edit add resource ../$$res && \
			cd - &> /dev/null; \
	done
endif
#
# Build the resources
# Post-process RoleBindings to fix the namespace in the refs (not yet handled by kustomize)
# Either apply to the cluster or output to CLI
#
ifeq ($(DRY_RUN), false)
	@$(KUSTOMIZE) build $(KOPTIONS) $@ | \
		sed 's/$(PLACEHOLDER)/$(NAMESPACE)/' | \
		kubectl apply --server-side -f -
else
	@$(KUSTOMIZE) build $(KOPTIONS) $@ | \
		sed 's/$(PLACEHOLDER)/$(NAMESPACE)/'
endif

.PHONY: operator .operator-port-patch .operator-can-monitor .operator-log-level-patch

#
# Customizes the port patch
#
.operator-port-patch:
	@sed -i 's/--monitoring-port=.*/--monitoring-port=$(MONITORING_PORT)/' $(MANAGER)/$(PORTS_PATCH).$(YAML)
	@sed -i '/path:.*\/containerPort/,/- op/{s/value: .*/value: $(MONITORING_PORT)/}' $(MANAGER)/$(PORTS_PATCH).$(YAML)
	@sed -i 's/--health-port=.*/--health-port=$(HEALTH_PORT)/' $(MANAGER)/$(PORTS_PATCH).$(YAML)
	@sed -i '/path:.*\/httpGet\/port/,/- op/{s/value: .*/value: $(HEALTH_PORT)/}' $(MANAGER)/$(PORTS_PATCH).$(YAML)

#
# Customizes the log level patch
#
.operator-log-level-patch:
	@sed -i 's/    value:.*/    value: "$(LOGGING_LEVEL)"/' $(MANAGER)/$(LOG_LEVEL_PATCH).$(YAML)

.operator-can-monitor: kubectl
	@output=$$(kubectl get crd prometheusrules.monitoring.coreos.com 2>&1) || (echo "****" && echo "**** ERROR: Montoring not available as Prometheus CRDs not installed in cluster ****" && echo "****"; exit 1)

#---
#
#@ operator
#
#== Install the operator deployment and related resources.
#
#=== Cluster-admin privileges are required.
#
#* PARAMETERS:
#** NAMESPACE:                Set the namespace to install the operator into
#** PLATFORM:                 Override the discovered platform, if required
#** GLOBAL:                   Set the operator to watch all namespaces for custom resources [true|false]
#** CUSTOM_IMAGE:             Set a custom operator image name
#** CUSTOM_VERSION:           Set a custom operator image version/tag
#** ALWAYS_PULL_IMAGES:       Set whether to always pull the operator image [true|false]
#** MONITORING:               Add the prometheus monitoring resources
#** MONITORING_PORT:          Set a custom monitoring port
#** HEALTH_PORT:              Set a custom health port
#** LOGGING_LEVEL:            Set the level of logging [info|debug]
#** INSTALL_DEFAULT_KAMELETS: Install the default Kamelets from catalog [true|false]
#** DRY_RUN:                  Print the resources to be applied instead of applying them [true|false]
#
#---
operator: check-admin have-platform check-crd-api-support kustomize kubectl .operator-port-patch .operator-log-level-patch
ifeq ($(MONITORING), true)
	@$(MAKE) -s .operator-can-monitor
	@$(call add-remove-operator-monitoring,$@,add)
else
	@$(call add-remove-operator-monitoring,$@,remove)
endif
# Set the namespace in the operator kustomization yaml
	@$(call set-kustomize-namespace,$@)
# Set the image reference of the kustomization
	@$(call set-kustomize-image,$@)
# Set the WATCH NAMESPACE env var depending on GLOBAL var
ifeq ($(GLOBAL), true)
	@$(call add-remove-kind-patch,$(MANAGER),add,$(WATCH_NAMESPACE_PATCH).$(YAML),Deployment)
else
	@$(call add-remove-kind-patch,$(MANAGER),remove,$(WATCH_NAMESPACE_PATCH).$(YAML),Deployment)
endif
# Set the ALWAYS_PULL_IMAGES config depending on var
ifeq ($(ALWAYS_PULL_IMAGES),true)
	@$(call add-remove-kind-patch,$(MANAGER),add,$(IMAGE_PULL_POLICY_PATCH).$(YAML),Deployment)
else
	@$(call add-remove-kind-patch,$(MANAGER),remove,$(IMAGE_PULL_POLICY_PATCH).$(YAML),Deployment)
endif
# Set the PORTS depending on vars
ifneq ($(MONITORING_PORT), 8080)
	@$(call add-remove-kind-patch,$(MANAGER),add,$(PORTS_PATCH).$(YAML),Deployment)
else ifneq ($(HEALTH_PORT), 8081)
	@$(call add-remove-kind-patch,$(MANAGER),add,$(PORTS_PATCH).$(YAML),Deployment)
endif
# Set the Log level of the operator
ifneq ($(LOGGING_LEVEL), info)
	@$(call add-remove-kind-patch,$(MANAGER),add,$(LOG_LEVEL_PATCH).$(YAML),Deployment)
else ifneq ($(LOGGING_LEVEL), 0)
		@$(call add-remove-kind-patch,$(MANAGER),add,$(LOG_LEVEL_PATCH).$(YAML),Deployment)
else
	@$(call add-remove-kind-patch,$(MANAGER),remove,$(LOG_LEVEL_PATCH).$(YAML),Deployment)
endif
# Set the INSTALL_DEFAULT_KAMELETS config depending on var
ifeq ($(INSTALL_DEFAULT_KAMELETS),false)
	@$(call add-remove-kind-patch,$(MANAGER),add,$(INSTALL_DEFAULT_KAMELETS_PATCH).$(YAML),Deployment)
else
	@$(call add-remove-kind-patch,$(MANAGER),remove,$(INSTALL_DEFAULT_KAMELETS_PATCH).$(YAML),Deployment)
endif
# Set the INSTALL_DEFAULT_PVC config depending on var
ifeq ($(INSTALL_DEFAULT_PVC),true)
	@$(call add-remove-kind-patch,$(MANAGER),add,$(INSTALL_DEFAULT_PVC_PATCH).$(YAML),Deployment)
else
	@$(call add-remove-kind-patch,$(MANAGER),remove,$(INSTALL_DEFAULT_PVC_PATCH).$(YAML),Deployment)
endif
ifeq ($(DRY_RUN), false)
	@$(KUSTOMIZE) build $(KOPTIONS) $@ | kubectl apply -f -
else
	@$(KUSTOMIZE) build $(KOPTIONS) $@
endif

.PHONY: platform .platform-openshift-patch .platform-kubernetes-patch

#
# Customizes the samples patches for kubernetes
#
.platform-kubernetes-patch:
	@sed -i 's/.*profile:.*/  profile: Kubernetes/' $(SAMPLES)/$(INT_PLATFORM_PATCH).$(YAML)

#
# Customizes the samples patches for openshift
#
.platform-openshift-patch:
	@sed -i 's/.*profile:.*/  profile: OpenShift/' $(SAMPLES)/$(INT_PLATFORM_PATCH).$(YAML)

#---
#
#@ platform
#
#== Install the integration platform.
#
#=== Cluster-admin privileges are required.
#
#* PARAMETERS:
#** NAMESPACE: Set the namespace to install the operator into
#** PLATFORM:  Override the discovered platform, if required
#** DRY_RUN:   Print the resources to be applied instead of applying them [true|false]
#
#---
platform: have-platform kustomize kubectl
# Cannot be a dependency as PLATFORM could contain 'ERROR: '
	@$(MAKE) .platform-$(PLATFORM)-patch
# Set the namespace in the platform kustomization yaml
	@$(call set-kustomize-namespace,$@)
ifeq ($(DRY_RUN), false)
	@$(KUSTOMIZE) build $(KOPTIONS) $@ | kubectl apply -f -
else
	@$(KUSTOMIZE) build $(KOPTIONS) $@
endif

.PHONY: example

#---
#
#@ example
#
#== Install the example integration.
#
#* PARAMETERS:
#** NAMESPACE: Set the namespace to install the example into
#** PLATFORM:  Override the discovered platform, if required
#** DRY_RUN:   Print the resources to be applied instead of applying them [true|false]
#
#---
example: kubectl
# Set the namespace in the example kustomization yaml
	@$(call set-kustomize-namespace,$@)
ifeq ($(DRY_RUN), false)
	@$(KUSTOMIZE) build $(KOPTIONS) $@ | kubectl apply -f -
else
	@$(KUSTOMIZE) build $(KOPTIONS) $@
endif

.PHONY: uninstall

#---
#
#@ uninstall
#
#== Uninstall the resources previously installed by the platform, operator, setup, and setup-cluster targets.
#
#=== Cluster-admin privileges are required.
#
#* PARAMETERS:
#** NAMESPACE:     Set the namespace to uninstall the resources from
#** UNINSTALL_ALL: Uninstall all Camel K resources including crds and cluster roles installed by setup-cluster [true|false]
#** DRY_RUN:       Print the resources to be applied instead of applying them [true|false]
#
#---
uninstall: kubectl
# Set the namespace in the all target kustomization yaml
	@$(call set-kustomize-namespace, platform)
	@$(call set-kustomize-namespace, operator)
	@$(call set-kustomize-namespace, setup)
ifeq ($(UNINSTALL_ALL), true)
	@$(call set-kustomize-namespace, setup-cluster)
endif
ifeq ($(DRY_RUN), false)
	@$(KUSTOMIZE) build $(KOPTIONS) platform | kubectl delete --ignore-not-found=true -f -
	@$(KUSTOMIZE) build $(KOPTIONS) operator | kubectl delete --ignore-not-found=true -f -
	@$(KUSTOMIZE) build $(KOPTIONS) setup    | kubectl delete --ignore-not-found=true -f -
ifeq ($(UNINSTALL_ALL), true)
	@$(KUSTOMIZE) build $(KOPTIONS) setup-cluster | kubectl delete --ignore-not-found=true -f -
endif
else
	@$(KUSTOMIZE) build $(KOPTIONS) platform | kubectl delete --dry-run=client -f -
	@$(KUSTOMIZE) build $(KOPTIONS) operator | kubectl delete --dry-run=client -f -
	@$(KUSTOMIZE) build $(KOPTIONS) setup    | kubectl delete --dry-run=client -f -
ifeq ($(UNINSTALL_ALL), true)
	@$(KUSTOMIZE) build $(KOPTIONS) setup-cluster | kubectl delete --dry-run=client -f -
endif
endif

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show this help screen.
	@#./script/help.sh $(MAKEFILE_LIST)
	@awk 'BEGIN { printf "\nUsage: make \033[31m<PARAM1=val1 PARAM2=val2>\033[0m \033[36m<target>\033[0m\n"; printf "\nAvailable targets are:\n" } /^#@/ { printf "\033[36m%-15s\033[0m", $$2; subdesc=0; next } /^#===/ { printf "%-14s \033[32m%s\033[0m\n", " ", substr($$0, 5); subdesc=1; next } /^#==/ { printf "\033[0m%s\033[0m\n\n", substr($$0, 4); next } /^#\*\*/ { printf "%-14s \033[31m%s\033[0m\n", " ", substr($$0, 4); next } /^#\*/ && (subdesc == 1) { printf "\n"; next } /^#\-\-\-/ { printf "\n"; next }' $(MAKEFILE_LIST)
