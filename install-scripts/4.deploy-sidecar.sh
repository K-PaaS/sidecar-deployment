#!/bin/bash

source variables.yml

kapp deploy -a sidecar -f manifest/sidecar-rendered.yml
