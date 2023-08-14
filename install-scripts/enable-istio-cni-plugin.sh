#!/bin/bash
SIDECAR_WORKING_DIR=$HOME/sidecar-deployment/install-scripts
ISTIO_DIR=$SIDECAR_WORKING_DIR/support-files/istio
ISTIO_OPERATOR_DIR=$ISTIO_DIR/istio-operator/manifests/charts/istio-operator

## istioctl 설치
cd $ISTIO_DIR
curl -L https://git.io/getLatestIstio | ISTIO_VERSION=1.12.6 sh -
cd istio-1.12.6
sudo cp bin/istioctl /usr/local/bin/istioctl

## istio operator 설치
cd $ISTIO_DIR
git clone https://github.com/istio/istio.git -b 1.12.6 istio-operator
kubectl apply -f $ISTIO_DIR/config/istio-namespace.yaml

cat << EOF >> $ISTIO_OPERATOR_DIR/values.yaml

sidecarLabel:
  labels:
    cloudfoundry.org/istio_version: 1.12.6
EOF

INPUT_LABEL_LINE=$(grep -n "name: istio-operator{{- if not (eq .Values.revision \"\") }}-{{ .Values.revision }}{{- end }}" $ISTIO_OPERATOR_DIR/templates/deployment.yaml | cut -d: -f1)
sed -i "${INPUT_LABEL_LINE} a\{{ toYaml .Values.sidecarLabel | trim | indent 2 }}" $ISTIO_OPERATOR_DIR/templates/deployment.yaml

INPUT_LABEL_LINE=$(grep -n "{{- if .Values.podAnnotations }}" $ISTIO_OPERATOR_DIR/templates/deployment.yaml | cut -d: -f1)
sed -i "${INPUT_LABEL_LINE} i\{{ toYaml .Values.sidecarLabel.labels | trim | indent 8 }}" $ISTIO_OPERATOR_DIR/templates/deployment.yaml

helm upgrade -i istio-operator $ISTIO_OPERATOR_DIR -f $ISTIO_DIR/config/istio-operator-values.yaml -n istio-system
kubectl apply -f $ISTIO_DIR/config/istio-cni.yaml

sleep 10

## istio operator 설치 확인
kubectl get all -o wide -n istio-system

cd $SIDECAR_WORKING_DIR