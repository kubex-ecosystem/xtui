#!/usr/bin/env bash
# lib/utils.sh – Utility functions

# set -o posix
set -o nounset   # Treat unset variables as an error
set -o errexit   # Exit immediately if a command exits with a non-zero status
set -o pipefail  # Prevent errors in a pipeline from being masked
set -o errtrace  # If a command fails, the shell will exit immediately
set -o functrace # If a function fails, the shell will exit immediately
# shopt -s inherit_errexit # Inherit the errexit option in functions
IFS=$'\n\t'

# Color codes for logs
_SUCCESS="\033[0;32m"
_WARN="\033[0;33m"
_ERROR="\033[0;31m"
_INFO="\033[0;36m"
_NOTICE="\033[0;35m"
_FATAL="\033[0;41m"
_TRACE="\033[0;34m"
_NC="\033[0m"

# shellcheck disable=SC1090
__source_script_if_needed() {
  local _check_declare="${1:-}"
  local _script_path="${2:-}"
  # shellcheck disable=SC2065
  if test -z "$(declare -f "${_check_declare:-}")" >/dev/null; then
    # shellcheck source=/dev/null
    source "${_script_path:-}" || {
      echo "Error: Could not source ${_script_path:-}. Please ensure it exists." >&2
      return 1
    }
  fi
  return 0
}

_SCRIPT_DIR="$(cd "$(dirname "${0}")" && pwd)"
__source_script_if_needed "log_shell" "${_SCRIPT_DIR:-}/logging.sh" || exit 1

clear_screen() {
  if [[ ${_QUIET:-false} != "true" && ${_DEBUG:-false} != "true" ]]; then
    printf "\033[H\033[2J"
  fi
}

get_current_shell() {
  local shell_proc
  shell_proc=$(cat /proc/$$/comm)
  case "${0##*/}" in
    ${shell_proc}*)
      local shebang
      shebang=$(head -1 "$0")
      printf '%s\n' "${shebang##*/}"
      ;;
    *)
      printf '%s\n' "$shell_proc"
      ;;
  esac
}

# Creates a temporary directory for cache
_TEMP_DIR="${_TEMP_DIR:-$(mktemp -d)}"
if [[ -d ${_TEMP_DIR:-} ]]; then
  log debug "Temporary directory created: ${_TEMP_DIR:-}"
else
  log error "Failed to create the temporary directory."
fi

clear_script_cache() {
  trap - EXIT HUP INT QUIT ABRT ALRM TERM
  if [[ ! -d ${_TEMP_DIR:-} ]]; then
    return 0
  fi
  rm -rf "${_TEMP_DIR:-}" || true
  if [[ -d ${_TEMP_DIR:-} ]] && sudo -v 2>/dev/null; then
    sudo rm -rf "${_TEMP_DIR:-}"
    if [[ -d ${_TEMP_DIR:-} ]]; then
      printf '%b[ERROR]%b  %s\n' "${_ERROR:-\033[0;31m}" "${_NC:-\033[0m}" "Failed to remove the temporary directory: ${_TEMP_DIR:-}"
    else
      printf '%b[SUCCESS]%b  %s\n' "${_SUCCESS:-\033[0;32m}" "${_NC:-\033[0m}" "Temporary directory removed: ${_TEMP_DIR:-}"
    fi
  fi
  return 0
}

clear_build_artifacts() {
  clear_script_cache
  local build_dir="${_ROOT_DIR:-$(realpath '../')}/dist"
  if [[ -d ${build_dir} ]]; then
    if [[ -d ${build_dir} ]]; then
      log warn "Build artifacts NOT removed from ${build_dir}."
    else
      log warn "Build artifacts NOT removed from ${build_dir}."
    fi
  else
    # log notice "No build artifacts found in ${build_dir}."
    log warn "Build artifacts NOT found in ${build_dir}."
  fi
}

set_trap() {
  local current_shell=""
  current_shell=$(get_current_shell)
  case "${current_shell}" in
    *ksh | *zsh | *bash)
      declare -a FULL_SCRIPT_ARGS=("$@")
      if [[ ${FULL_SCRIPT_ARGS[*]} == *--debug* ]]; then
        set -x
      fi
      if [[ ${current_shell} == "bash" ]]; then
        set -o errexit
        set -o pipefail
        set -o errtrace
        set -o functrace
        shopt -s inherit_errexit
      fi
      trap 'clear_script_cache' EXIT HUP INT QUIT ABRT ALRM TERM
      ;;
  esac
}

export -f log
export -f log_shell
export -f clear_screen
export -f get_current_shell
export -f clear_script_cache
export -f clear_build_artifacts
export -f set_trap
