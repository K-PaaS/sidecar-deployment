#!/bin/bash
source variables.yml

# delete sidecar
helm uninstall sidecar --namespace="$sidecar_namespace" 
