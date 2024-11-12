#!/bin/bash

source portal-deploy-variables.yml

cf delete-org $ORG_NAME -f
