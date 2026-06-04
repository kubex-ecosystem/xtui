#!/usr/bin/env bash

set -euo pipefail

########################################################
# lista todas as chaves de serviço de todos os projetos
# salva em arquivos no /tmp
#
# só funciona se o user logado tiver permissão

# Checa se está autenticado no GCP
if ! gcloud auth list | grep -q "ACTIVE"; then
  echo "Nenhum usuário autenticado no GCP." >&2
  exit 1
fi

list_gcp_keys() {
  # Faremos tudo em /tmp para não poluir o home do usuário
  cd /tmp || exit 1

  for project in $(gcloud projects list --format="value(projectId)")
  do
    echo "ProjectId:  $project"
    for robot in $(gcloud iam service-accounts list --project "$project" --format="value(email)")
    do
      echo "    -> Robot $robot"
      for key in $(gcloud iam service-accounts keys list --iam-account "$robot" --project "$project" --format="value(name.basename())")
          do
            echo "        $key"
      done
    done
  done

  return 0
}


list_gcp_keys "$@" || {
  echo 'Falha ao listar as chaves de serviço do GCP.' >&2
  exit 1
}
