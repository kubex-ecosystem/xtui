#!/usr/bin/env bash

set -euo pipefail

########################################################
# liste VMs e suas service accounts e scopes
#


process_output_gcloud_a01() {
  # 1. lista todas as instancias
  for scopesInfo in $(
      gcloud compute instances list \
          --format="csv[no-heading](name,id,serviceAccounts[].email.list(),
                        serviceAccounts[].scopes[].map().list(separator=;))")
  do
        IFS=',' read -r -a scopesInfoArray<<< "$scopesInfo"
        NAME="${scopesInfoArray[0]}"
        ID="${scopesInfoArray[1]}"
        EMAIL="${scopesInfoArray[2]}"
        SCOPES_LIST="${scopesInfoArray[3]}"

        echo "NAME: $NAME, ID: $ID, EMAIL: $EMAIL"
        echo ""
        IFS=';' read -r -a scopeListArray<<< "$SCOPES_LIST"
        for SCOPE in  "${scopeListArray[@]}"
        do
          echo "  SCOPE: $SCOPE"
        done
  done
}


# process_output_gcloud_a02() {
#   for scopesInfo in $(gcloud compute instances list --format="csv[no-heading](name,id,serviceAccounts[].email.list(),
#                         serviceAccounts[].scopes[].map().list(separator=;))"); do
#     IFS=',' read -r -a scopesInfoArray <<<"$scopesInfo"
#     NAME="${scopesInfoArray[0]}"
#     ID="${scopesInfoArray[1]}"
#     EMAIL="${scopesInfoArray[2]}"
#     SCOPES_LIST="${scopesInfoArray[3]}"

#     echo "NAME: $NAME, ID: $ID, EMAIL: $EMAIL"
#     echo ""
#     IFS=';' read -r -a scopeListArray <<<"$SCOPES_LIST"
#     for SCOPE in "${scopeListArray[@]}"; do
#       echo "  SCOPE: $SCOPE"
#     done
#   done
# }


process_output_gcloud_a01 "$@" || {
  echo 'Falha ao processar a saída do gcloud.' >&2
  echo "$@"
}

# process_output_gcloud_a02"$@" || {
#   echo 'Falha ao processar a saída do gcloud.' >&2
#   exit 1
# }
