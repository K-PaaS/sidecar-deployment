#!/bin/bash
source variables.yml

if [[ ${use_external_blobstore} = "true" ]]; then
  if [[ ${use_external_db} = "true" ]]; then
    if [[ ${external_db_kind} = "postgres" ]]; then
      # External_Blobstore & External Postgres
      ytt -f ../config -f "manifest/sidecar-values.yml" -f "manifest/external-blobstore-values.yml" -f "manifest/external-db-values-postgresql.yml" > "manifest/sidecar-rendered.yml"
    elif [[ ${external_db_kind} = "mysql" ]]; then
      # External_Blobstore & External MySQL
      ytt -f ../config -f "manifest/sidecar-values.yml" -f "manifest/external-blobstore-values.yml" -f "manifest/external-db-values-mysql.yml" > "manifest/sidecar-rendered.yml"
    else
      # Error Check : external_db_kind
      echo "plz check variables.yml : external_db_kind"
      return
    fi
  elif [[ ${use_external_db} = "false" ]]; then
    # External_Blobstore
    ytt -f ../config -f "manifest/sidecar-values.yml" -f "manifest/external-blobstore-values.yml" > "manifest/sidecar-rendered.yml"
  else
    # Error Check : use_external_db
    echo "plz check variables.yml : use_external_db"
    return
  fi

elif [[ ${use_external_blobstore} = "false" ]]; then
  if [[ ${use_external_db} = "true" ]]; then
    if [[ ${external_db_kind} = "postgres" ]]; then
      # External Postgres
      ytt -f ../config -f "manifest/sidecar-values.yml" -f "manifest/external-db-values-postgresql.yml" > "manifest/sidecar-rendered.yml"
    elif [[ ${external_db_kind} = "mysql" ]]; then
      # External MySQL
      ytt -f ../config -f "manifest/sidecar-values.yml" -f "manifest/external-db-values-mysql.yml" > "manifest/sidecar-rendered.yml"
    else
      # Error Check : external_db_kind
      echo "plz check variables.yml : external_db_kind"
      return
    fi
  elif [[ ${use_external_db} = "false" ]]; then
    #normal
    ytt -f ../config -f "manifest/sidecar-values.yml" > "manifest/sidecar-rendered.yml"
  else
    # Error Check : use_external_db
    echo "plz check variables.yml : use_external_db"
    return
  fi
else
  echo "plz check variables.yml : use_external_blobstore"
  return
fi
