#!/usr/bin/env bash
# shellcheck disable=SC2015,SC1091,SC1090,SC2086

# Script Metadata
__secure_logic_version="1.0.0"
__secure_logic_date="$(date +%Y-%m-%d)"
__secure_logic_author="Rafael Mori"
__secure_logic_use_type="lib"
__secure_logic_init_timestamp="$(date +%s)"
__secure_logic_elapsed_time=0

# Check if verbose mode is enabled
if [[ ${MYNAME_VERBOSE:-false} == "true" ]]; then
  set -x # Enable debugging
fi

IFS=$'\n\t'

set -euo pipefail
set -o errtrace
set -o functrace
set -o posix

declare -a _main_args=("$@")

__secure_logic_sourced_name() {
  local _self="${BASH_SOURCE-}"
  _self="${_self//${_kbx_root:-$()}/}"
  _self="${_self//\.sh/}"
  _self="${_self//\-/_}"
  _self="${_self//\//_}"
  _self="${_self//\./_}"
  _self="${_self//__/_}"
  _self="${_self//[^a-zA-Z0-9_]/}"
  # Return a unique variable name based on the script name
  echo "_was_sourced_${_self//__/_}"
  return 0
}

__first() {
  if [ "$EUID" -eq 0 ] || [ "$UID" -eq 0 ]; then
    echo "Please do not run as root." >__get_output_tty
    exit 1
  elif [ -n "${SUDO_USER:-}" ]; then
    echo "Please do not run as root, but with sudo privileges." >__get_output_tty
    exit 1
  else
    # shellcheck disable=SC2155
    local _ws_name="$(__secure_logic_sourced_name)"

    if test "${BASH_SOURCE-}" != "${0}"; then
      if test ${__secure_logic_use_type:-} != "lib"; then
        echo "This script is not intended to be sourced." >__get_output_tty
        echo "Please run it directly." >__get_output_tty
        exit 1
      fi
      # If the script is sourced, we set the variable to true
      # and export it to the environment without changing
      # the shell options.
      export "${_ws_name:-}"="true"
    else
      if test ${__secure_logic_use_type:-} != "exec"; then
        echo "This script is not intended to be executed directly." >__get_output_tty
        echo "Please source it instead." >__get_output_tty
        exit 1
      fi
      # If the script is executed directly, we set the variable to false
      # and export it to the environment. We also set the shell options
      # to ensure a safe execution.
      export "${_ws_name:-}"="false"
      set -o errexit           # Exit immediately if a command exits with a non-zero status
      set -o pipefail          # Return the exit status of the last command in the pipeline that failed
      set -o errtrace          # If a command fails, the shell will exit immediately
      set -o functrace         # If a function fails, the shell will exit immediately
      shopt -s inherit_errexit # Inherit the errexit option in functions

      if [[ ${_DEBUG:-} == "true" ]]; then
        set -x
      fi
    fi
  fi
}

__first "${_main_args[@]}" >&2 || {
  echo "Error: This script must be run directly, not sourced." >&2
  exit 1
}

__source_script_if_needed() {
  local _check_declare="${1:-}"
  local _script_path="${2:-}"
  # shellcheck disable=SC2065
  if test -z "$(declare -f "${_check_declare:-}")" >/dev/null; then
    # shellcheck disable=SC2086,SC1090,SC1091,SC2065
    source "${_script_path:-}" || {
      echo "Error: Could not source ${_script_path:-}. Please ensure it exists." >&2
      return 1
    }
  fi
  return 0
}

# Load library files
_SCRIPT_DIR="$(cd "$(dirname "${0}")" && pwd)"
__source_script_if_needed "show_summary" "${_SCRIPT_DIR:-}/config.sh" || exit 1
__source_script_if_needed "apply_manifest" "${_SCRIPT_DIR:-}/apply_manifest.sh" || exit 1
__source_script_if_needed "get_current_shell" "${_SCRIPT_DIR:-}/utils.sh" || exit 1

__load_environment() {
  _PID=
  _RUN_ENV_FILE=""
  _DEBUG=${_DEBUG:-${DEBUG:-false}}
  _ROOT_DIR="${_ROOT_DIR:-${ROOT_DIR:-"$(git rev-parse --show-toplevel)"}}"
  _SCRIPT_DIR="$(dirname "${0}")"

  _MAIN_PKG="$(dirname "$(grep -risn --exclude-dir={docs,examples,version,swagger} '^package main' "$(realpath "${_ROOT_DIR}")" | head -n1 | awk -F ':' '{print $1}')")"
  _MAIN_PKG="${_MAIN_PKG:-${_ROOT_DIR}/cmd}"
  _BINARY="${XTUI_BINARY:-${_ROOT_DIR}/dist/xtui_linux_amd64}"
  _WEB_DIR="${XTUI_WEB_DIR:-${_ROOT_DIR}/frontend/dist}"
  _LOG_PATH="/tmp/$(basename "$(git rev-parse --show-toplevel)").txt"
  _FAIL_FLAG=0

  log fatal "In development process"

}

clear_screen

# ====== Utility functions ======
have() {
  command -v "${1}" >/dev/null 2>&1 && return 0 || return 1
}

cleanup() {
  if [[ -n $_PID ]]; then
    kill "$_PID" 2>/dev/null || true
    wait "$_PID" 2>/dev/null || true
  fi

  if ! pgrep -f "${_BINARY}" 2>/dev/null; then
    # Removendo o arquivo de log
    rm -f "$_LOG_PATH" || true
  else
    log info "Processo ainda em execução. Aguarde até o encerramento para remover o arquivo de log."
  fi
}

# Trap para limpeza de resíduos desse script
trap cleanup EXIT HUP INT QUIT ABRT ALRM TERM

__launch_server() {
  if [ ! -f "${_BINARY}" ]; then
    log info "Compilando o Gateway XTUI..."
    cd "${_ROOT_DIR}" || {
      log error "não foi possível navegar até o diretório raiz"
      exit 1
    }
    if command -v make 2>/dev/null; then
      make build-dev || {
        log error "falha ao compilar"
        exit 1
      }
    else
      go build -ldflags "-X main.envFile=${_RUN_ENV_FILE}" -trimpath -o "$_BINARY" "${_MAIN_PKG}/main.go" || {
        log error "falha ao compilar"
        exit 1
      }
    fi
  fi

  local wait=${1:-false}

  if [ ! -f "${_LOG_PATH}" ]; then
    echo "" >"${_LOG_PATH}" || {
      log error "não foi possível criar o arquivo de log"
      return 1
    }
    chmod 666 "${_LOG_PATH}" 2>/dev/null || true
  fi

  # ---------- lançar servidor ----------
  log info "iniciando xtui gateway na porta ${_PORT}..."
  local cmd=(
    "${_BINARY}"
    "gateway"
    up
    "-e"
    "${_RUN_ENV_FILE}"
    "--web-dir"
    "${_WEB_DIR}"
    "--port"
    "${_PORT}"
    "-D"
  )
  "${cmd[@]}" >"${_LOG_PATH}" 2>&1 &
  _PID=$!

  if [ -z "$_PID" ]; then
    log error "Failed to start ${APP_NAME:-$(basename "$(git rev-parse --show-toplevel)")} server"
    return $_FAIL_FLAG
  fi

  if [[ -f ${_LOG_PATH} ]]; then
    tail -n 10 "${_LOG_PATH}"
  else
    log error "${_LOG_PATH} not found"
    return 1
  fi

  kill -CONT "$_PID" 2>/dev/null || true

  log info "${APP_NAME:-$(basename "$(git rev-parse --show-toplevel)")} PID: $_PID"

  sleep 1

  if ! __wait_for_health; then
    log error "Failed to start ${APP_NAME:-$(basename "$(git rev-parse --show-toplevel)")} server"
    return $_FAIL_FLAG
  fi

  if [ "$wait" == "true" ]; then
    log info "Iniciando acompanhamento do log do processo..."
    log info "Para encerrar o acompanhamento do processo, pressione Ctrl+C."
    tail -f "${_LOG_PATH}" || {
      log error "Falha ao iniciar acompanhamento do log"
      return 1
    }
    __kill_server || true
    wait "$_PID" 2>/dev/null || true
    cleanup
    log success "${APP_NAME:-} PID: $_PID encerrado com sucesso"
  else
    log success "Inicializacao concluida com sucesso"
  fi
  return 0
}

check() {
  local label="$1" url="$2" expect="$3"
  local body
  # -s sem -f: queremos inspecionar o corpo mesmo em respostas 4xx/5xx
  body=$(curl -s "$url" 2>/dev/null || echo "CURL_FAILED")
  if echo "$body" | grep -q "$expect"; then
    log info "OK   $label"
  else
    log error "FAIL $label — esperado '$expect', obtido: $body"
    _FAIL_FLAG=1
  fi
}

__wait_for_health() {
  log info "Waiting for ${APP_NAME:-$(basename "$(git rev-parse --show-toplevel)")} server to be healthy..."
  # aguardar health
  local ATTEMPTS=0

  until curl -X GET "${BASE}/api/v1/health" >/dev/null 2>&1; do
    sleep 1
    ATTEMPTS=$((ATTEMPTS + 1))
    if [[ $ATTEMPTS -gt 5 ]]; then
      log error "servidor não respondeu em 5s. (${ATTEMPTS} tentativas)"
      return 1
    fi
  done
  log info "servidor up após ${ATTEMPTS}s"
  return 0
}

__kill_server() {
  if [[ -n $_PID ]]; then
    kill "$_PID" 2>/dev/null || true
    wait "$_PID" 2>/dev/null || true
  else
    log info "PID não informado. Procurando processos para encerrar"
    local pid=

    # Try to find the process using pgrep, if not found, try to find it using ps aux
    pid=$(pgrep -f "${_BINARY}")
    # shellcheck disable=SC2009
    pid=${pid:-$(ps aux | grep "${_BINARY}" | grep -v grep | awk '{print $2}')}

    if [[ -n $pid ]]; then
      log info "Encerrando processo com PID: $pid"
      kill "$pid" 2>/dev/null || true
      timeout 5 wait "$pid" 2>/dev/null || true
      if [[ -n "$(ps -p "$pid" -o pid= 2>/dev/null)" ]]; then
        log error "Falha ao encerrar processo com PID: $pid"
        _FAIL_FLAG=1
      else
        log success "Processo com PID: $pid encerrado com sucesso"
      fi
    else
      log warn "Processo com PID: $pid não encontrado"
    fi
  fi
}

run() {
  __load_environment "${_main_args[@]}" || {
    # o log fatal já encerra o processo que foi iniciado pelo makefile
    log fatal "Falha ao carregar ambiente"
  }
  __launch_server "${_main_args[@]}" || {
    # o log fatal já encerra o processo que foi iniciado pelo makefile
    log fatal "Falha ao iniciar servidor"
  }
}

stop() {
  __load_environment "${_main_args[@]}" || {
    # o log fatal já encerra o processo que foi iniciado pelo makefile
    log fatal "Falha ao carregar ambiente"
  }
  __kill_server || {
    log fatal "Falha ao encerrar servidor"
  }
}

status() {
  __load_environment "${_main_args[@]}" || {
    # o log fatal já encerra o processo que foi iniciado pelo makefile
    log fatal "Falha ao carregar ambiente"
  }
  local status_code=
  status_code=$(curl -s -o /dev/null -w "%{http_code}" "${BASE}/api/v1/health" | cut -d '%' -f1 || echo "000")
  local status_body=
  status_body=$(curl -sf "${BASE}/api/v1/health" 2>/dev/null || echo "000")
  if [[ $status_code != "200" ]]; then
    log warn "Server Status Code:"
    log warn "$status_code"
    log warn "Server Status Body:"
    log warn "$(echo "${status_body}" | jq .)"
    _FAIL_FLAG=1
  else
    log info "Server Status Code:"
    log info "$status_code"
    log info "Server Status Body:"
    log info "$(echo "${status_body}" | jq .)"
  fi
}

restart() {
  __load_environment "${_main_args[@]}" || {
    # o log fatal já encerra o processo que foi iniciado pelo makefile
    log fatal "Falha ao carregar ambiente"
  }
  __kill_server || {
    log fatal "Falha ao encerrar servidor"
  }
  __launch_server || {
    log fatal "Falha ao iniciar servidor"
  }
}

logs() {
  __load_environment "${_main_args[@]}" || {
    # o log fatal já encerra o processo que foi iniciado pelo makefile
    log fatal "Falha ao carregar ambiente"
  }
  local follow="${1:-false}"
  local size="${2:-100}"

  if [[ $follow == "true" ]]; then
    tail -n "$size" -f "$_LOG_PATH" || {
      # o fatal independe de quiet, verbose, etc... Ele sempre será exibido.
      log fatal "Falha ao ler logs"
      return 1
    }
  else
    tail -n "$size" "$_LOG_PATH" || {
      log fatal "Falha ao ler logs"
      return 1
    }
  fi

  # o true no terceiro argumento força a exibição no stdout mesmo se estiver quiet habilitado
  log success "Encerrada a exibição de logs" #true

  return $_FAIL_FLAG
}

export_logs() {
  __load_environment "${_main_args[@]}" || {
    # o log fatal já encerra o processo que foi iniciado pelo makefile
    log fatal "Falha ao carregar ambiente"
  }
  local target_path="${1:-$(date +"${_INSTALL_DIR}/logs/%Y/%m/%d/%H/%M/%S.log")}"
  local _dest_dir=""

  if [[ ! -f $_LOG_PATH ]]; then
    log error "Arquivo de log não encontrado: $_LOG_PATH"
    return 1
  fi

  _dest_dir="$(dirname "$target_path")"

  if [[ ! -d $_dest_dir ]]; then
    mkdir -p "$_dest_dir" || {
      log error "Falha ao criar diretório de destino: $_dest_dir"
      return 1
    }
  fi

  cp "$_LOG_PATH" "$target_path" || {
    log error "Falha ao copiar arquivo de log para: $target_path"
    return 1
  }

  log success "Arquivo de log exportado com sucesso para: $target_path"

  return $_FAIL_FLAG
}

get_cmd() {
  __load_environment "${_main_args[@]}" || {
    # o log fatal já encerra o processo que foi iniciado pelo makefile
    log fatal "Falha ao carregar ambiente"
  }
  local _function_name="${1:-}"

  if [[ -z $_function_name ]]; then
    log error "Nome da função não especificado"
    return 1
  fi

  type "$_function_name" || {
    log error "Função $_function_name não encontrada"
    return 1
  }

  return $_FAIL_FLAG
}

# Não podemos usar trap pois o server é enviado para background e o script atual não irá manter-se
# em execução até o encerramento do server. Logo, se ele acabar de subir e encerrar o fluxo do próprio script
# com sucesso, ele irá encerrar o servidor.
# trap __kill_server EXIT

export -f run
export -f stop
export -f logs
export -f status
export -f restart
export -f export_logs
export -f have
