#!/bin/bash
source variables.yml
kubectl patch sc $storageclass_name -p '{"metadata":{"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
