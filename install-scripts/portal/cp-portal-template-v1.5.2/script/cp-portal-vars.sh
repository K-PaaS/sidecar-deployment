# COMMON VARIABLE (Please change the value of the variables below.)
K8S_MASTER_NODE_IP="${kubectl config view | yq e '.clusters[0].cluster.server' | cut -d '/' -f 3 | cut -d ':' -f 1
}" 									  # Kubernetes Master Node Public IP
K8S_CLUSTER_API_SERVER="https://${K8S_MASTER_NODE_IP}:6443"               # kubernetes API Server (e.g. https://${K8S_MASTER_NODE_IP}:6443)
K8S_STORAGECLASS="cp-storageclass"                                        # Kubernetes StorageClass Name (e.g. cp-storageclass)
HOST_CLUSTER_IAAS_TYPE="1"                                                # Kubernetes Cluster IaaS Type ([1] AWS, [2] OPENSTACK, [3] NAVER, [4] NHN, [5] KT)
HOST_DOMAIN=INGRESS_HOST_DOMAIN                                               # Host Domain (e.g. xx.xxx.xxx.xx.nip.io)

# The belows are the default values.
# If you change the values below, there will be a problem with the deploy. Please keep the values.
NAMESPACE=(
"vault"
"harbor"
"mariadb"
"keycloak"
"cp-portal"
)

IAAS_TYPE=(
"AWS"
"OPENSTACK"
"NAVER"
"NHN"
"KT"
)

CHART_NAME=(
"cp-vault"
"cp-harbor"
"cp-mariadb"
"cp-keycloak"
"cp-app"
"cp-cert-setup"
)

IMAGE_NAME=(
"cp-portal-common-api"
"cp-portal-metric-api"
)

IMAGE_TAGS="latest"                                                                          # image tag
IMAGE_PULL_POLICY="Always"                                                                   # image pull policy
IMAGE_PULL_SECRET="cp-secret"                                                                # image pull secret
SERVICE_TYPE="ClusterIP"                                                                     # service type in kubernetes
SERVICE_PROTOCOL="TCP"                                                                       # service protocol in kubernetes

# K8S
K8S_CLUSTER_ADMIN="cp-cluster-admin"                                                         # kubernetes cluster-admin role name
K8S_CLUSTER_ADMIN_NAMESPACE="kube-system"                                                    # kubernetes cluster-admin role namespace

# INGRESS CONTROLLER
CP_DEFAULT_INGRESS_NAMESPACE="ingress-nginx"                                                 # container platform default ingress namespace
CP_DEFAULT_INGRESS_CLASS_NAME="nginx"                                                        # container platform default ingress name
CP_DEFAULT_INGRESS_CONTROLLER_SELECTOR="app.kubernetes.io/component=controller"              # container platform default ingress controller selector

# VAULT
VAULT_URL="http://vault.${HOST_DOMAIN}"                                                      # vault url
VAULT_ROLE_NAME="cp_role"                                                                    # vault role name

# HARBOR
REPOSITORY_URL="https://harbor.${HOST_DOMAIN}"                                               # harbor url
REPOSITORY_USERNAME="admin"                                                                  # harbor admin username (e.g. admin)
REPOSITORY_PASSWORD="Harbor12345"                                                            # harbor admin password (e.g. Harbor12345)
REPOSITORY_PROJECT_NAME="cp-portal-repository"                                               # harbor project name

# MARIADB
DATABASE_URL="cp-mariadb.mariadb.svc.cluster.local:3306"                                     # database url
DATABASE_USER_ID="cp-admin"                                                                  # database user name (e.g. cp-admin)
DATABASE_USER_PASSWORD="cpAdmin!12345"                                                       # database user password (e.g. cpAdmin!12345)
DATABASE_TERRAMAN_ID="terraman"                                                              # database user name (e.g. terraman)
DATABASE_TERRAMAN_PASSWORD="cpAdmin!12345"                                                   # database user name (e.g. cpAdmin!12345)

# KEYCLOAK
KEYCLOAK_URL="http://keycloak.${HOST_DOMAIN}"                                                # keycloak url (if apply TLS, https:// )
KEYCLOAK_DB_VENDOR="mariadb"                                                                 # keycloak database vendor
KEYCLOAK_DB_SCHEMA="keycloak"                                                                # keycloak database schema
KEYCLOAK_ADMIN_USERNAME="admin"                                                              # keycloak admin username (e.g. admin)
KEYCLOAK_ADMIN_PASSWORD="admin"                                                              # keycloak admin password (e.g. admin)
KEYCLOAK_SESSIONS_COUNT="2"                                                                  # keycloak sessions count
KEYCLOAK_LOG_LEVEL="INFO"                                                                    # keycloak log level
KEYCLOAK_CP_REALM="cp-realm"                                                                 # keycloak realm for container platform portal
KEYCLOAK_CP_CLIENT_ID="cp-client"                                                            # keycloak client id for container platform portal
KEYCLOAK_CP_CLIENT_SECRET="bfbc9f86-81f4-4307-b57e-c84e943a5bc5"                             # keycloak client secret for container platform portal
KEYCLOAK_INGRESS_TLS_ENABLED="false"                                                         # keycloak ingress tls enabled (if apply TLS, true)
KEYCLOAK_TLS_CERT_PATH="path/to/cert/file"                                                   # keycloak tls cert file path (if apply TLS, cert file path)
KEYCLOAK_TLS_KEY_PATH="path/to/key/file"                                                     # keycloak tls key file path (if apply TLS, key file path)
KEYCLOAK_TLS_SECRET="cp-keycloak-tls-secret"                                                 # keycloak tls secret name

# CP_PORTAL
RELEASE_NAME="cp-portal"                                                                     # container platform portal release name
CP_PORTAL_URL="http://portal.${HOST_DOMAIN}"                                                 # container platform portal url
HOST_CLUSTER_NAME="host-cluster"                                                             # host cluster name

# CP_SERVICE
CP_SERVICE_PIPELINE_NAMESPACE="cp-pipeline"                                                  # container platform service pipeline namespace
CP_SERVICE_SOURCE_CONTROL_NAMESPACE="cp-source-control"                                      # container platform service source control namespace
CP_SERVICE_PIPELINE_URL="http://pipeline.${HOST_DOMAIN}"                                     # container platform service pipeline url
CP_SERVICE_SOURCE_CONTROL_URL="http://scm.${HOST_DOMAIN}"                                    # container platform service source control url

# CP CERT SETUP
CP_CERT_SETUP_NAMESPACE="kube-system"                                                        # container platform cert setup namespace
CP_CERT_SETUP_SELECTOR="app=${CHART_NAME[5]}"                                                # container platform cert setup selector
