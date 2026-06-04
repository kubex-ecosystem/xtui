#!/usr/bin/env bash
# shellcheck disable=SC2015,SC1091,SC1090,SC2086

trap - EXIT HUP INT TERM

# Script Metadata
__secure_logic_version="1.0.0"
__secure_logic_date="$(date +%Y-%m-%d)"
__secure_logic_author="Rafael Mori"
__secure_logic_use_type="exec"
__secure_logic_init_timestamp="$(date +%s)"
__secure_logic_elapsed_time=0

# Check if verbose mode is enabled
if [[ ${MYNAME_VERBOSE:-false} == "true" ]]; then
  set -x # Enable debugging
fi

IFS=$'\n\t'
declare -a _main_pb_args=("$@")

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
        # echo "This script is not intended to be sourced." >__get_output_tty
        # echo "Please run it directly." >__get_output_tty
        # exit 1
        return 0
      fi
      # If the script is sourced, we set the variable to true
      # and export it to the environment without changing
      # the shell options.
      export "${_ws_name:-}"="true"
    else
      if test ${__secure_logic_use_type:-} != "exec"; then
        # echo "This script is not intended to be executed directly." >__get_output_tty
        # echo "Please source it instead." >__get_output_tty
        # exit 1
        return 0
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

__first "${_main_pb_args[@]}" >&2 || {
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
_ROOT_DIR="$(git rev-parse --show-toplevel)"
__EXP_SCRIPT_DIR="$(realpath "${_ROOT_DIR:-$(git rev-parse --show-toplevel)}/support")"
__source_script_if_needed "show_summary" "${__EXP_SCRIPT_DIR:-}/config.sh" || exit 1
__source_script_if_needed "apply_manifest" "${__EXP_SCRIPT_DIR:-}/apply_manifest.sh" || exit 1
__source_script_if_needed "get_current_shell" "${__EXP_SCRIPT_DIR:-}/utils.sh" || exit 1
__source_script_if_needed "what_platform" "${_SCRIPT_DIR:-}/platform.sh" || exit 1

function _cleanup_tmp_file {
  if test -f "${_tmp_file_path:-}"; then
    rm -f "${_tmp_file_path:-}" 2>/dev/null || {
      rm -f "/tmp/XTUIPB.*" 2>/dev/null
      printf "%s\n" "Error: Could not remove temporary file." >&2
      return 1
    }
  fi
}

__create_tmp_file() {
  local _tmp_file_path
  _tmp_file_path="$(mktemp -t XTUIPB.XXXXXX)" || {
    log fatal "Could not create temporary file."
  }

  export -f _cleanup_tmp_file
  trap '_cleanup_tmp_file' EXIT TERM INT HUP

  printf '%s\n' "0" >"${_tmp_file_path}" || {
    log fatal "Could not write to temporary file."
  }
  echo "${_tmp_file_path}"

  return 0
}

__secure_logic_main() {
  local _ws_name
  _ws_name="$(__secure_logic_sourced_name)"
  local _ws_name_val
  _ws_name_val=$(eval "echo \${${_ws_name:-}}")
  if test "${_ws_name_val:-}" != "true"; then
    __run_xtui_progress_bar_example "${_main_pb_args[@]}"
    return $?
  else
    # If the script is sourced, we export the functions
    # log error "This script is not intended to be sourced." true
    # log error "Please run it directly." true
    return 0
  fi
}

__fallback_cmd() {
  local _tmp_file_path="${1:-"/tmp/xtui_progress_bar.tmp"}"

  for i in {1..100}; do
    printf '%s\n' "$i" >"${_tmp_file_path}" || {
      log error "Could not write to temporary file." true
      return 1
    }
    sleep 0.2 || {
      log error "Could not sleep." true
      return 1
    }
  done
  return 0
}

# Simula um trabalho pesado em background que atualiza o arquivo temporário
__run_xtui_progress_bar_example() {
  # Executa o xtui no foreground para exibir e monitorar
  local _tmp_file_path
  local _cmd_to_run="${1:-}"
  local _cmd_timeout="${2:-5}"

  apply_manifest || {
    log fatal "Could not apply manifest."
  }
  local _platform_name
  if ! what_platform; then
    log error "Platform could not be determined." true
    return 1
  fi

  local _arrArgs=("${_main_pb_args[@]}")
  # local _arrArgs=( "${_args[@]::$#}" )

  local _command="${_arrArgs[0]:-help}"
  local _platform_arg="${_arrArgs[1]:-}"
  local _arch_arg="${_arrArgs[2]:-}"

  # If no platform specified, use cross-compilation mode
  if [[ -z ${_platform_arg} ]]; then
    _platform_arg="__CROSS_COMPILE__" # Special flag for cross-compilation
  fi

  # Set defaults only for specific platform requests
  if [[ ${_platform_arg} != "__CROSS_COMPILE__" ]]; then
    _arch_arg="${_arch_arg:-$(uname -m | tr '[:upper:]' '[:lower:]')}"
  fi

  local _root_dir="${_ROOT_DIR:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _cmd_path="${_CMD_PATH:-${CMD_PATH:-${_root_dir}/cmd}}"
  local _binary_name="${_BINARY_NAME:-${BINARY_NAME:-$(basename "${_cmd_path}" .go)}}"

  local _platform="${_PLATFORM:-${_CURRENT_PLATFORM:-}}"
  local _arch="${_ARCH:-${_CURRENT_ARCH:-}}"
  local _build_target="${_BUILD_TARGET:-${_platform}-${_arch}}"

  local _bin_path
  _bin_path="${_root_dir}/dist/${_BINARY_NAME:?BINARY_NAME is not set}_${_build_target//-/_}"

  if [ ! -x "${_bin_path:-}" ]; then
    log fatal "Could not find binary ${_bin_path:-}"
  fi

  _tmp_file_path=$(__create_tmp_file)

  if test -z "${_cmd_to_run:-}"; then
    export -f __fallback_cmd
    __fallback_cmd "${_tmp_file_path}" &

    ${_bin_path} f pb -f "${_tmp_file_path}" -t "Processando dados do sistema:" -w --timeout 29 -T 100 -D || {
      log fatal "Failed to run xtui progress-bar."
    }
  else
    timeout "${_cmd_timeout:-}" bash -c "$_cmd_to_run $_tmp_file_path" || {
      log fatal "Failed to run xtui progress-bar."
    }
  fi
}

main_progress_bar_example() {
  __secure_logic_main "${_main_pb_args[@]}" || {
    log fatal "Script execution failed." true
  }

  __secure_logic_elapsed_time="$(($(date +%s) - __secure_logic_init_timestamp))"

  if [[ ${MYNAME_VERBOSE:-false} == "true" || ${_DEBUG:-false} == "true" ]]; then
    log info "Script executed in ${__secure_logic_elapsed_time} seconds."
  fi
}
main_progress_bar_example "${_main_pb_args[@]}"
