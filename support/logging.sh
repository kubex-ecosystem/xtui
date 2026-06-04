#!/usr/bin/env bash

## <editor-fold defaultstate="collapsed" desc="Convert case">
function _camel_to_snake() {
  printf '%s' "${1:-}" | sed -E 's/([A-Z])/_\L\1/g' | sed 's/^_//' || return 1
}
function _camel_to_kebab() {
  printf '%s' "${1:-}" | sed -E 's/([A-Z])/-\L\1/g' | sed 's/^-//' || return 1
}
function _convert_case() {
  local text="$1"
  local target_case="$2"

  case "$target_case" in
    snake) printf '%s' "$text" | sed -E 's/([A-Z])/_\L\1/g' | sed 's/^_//' ;;
    kebab) printf '%s' "$text" | sed -E 's/([A-Z])/-\L\1/g' | sed 's/^-//' ;;
    camel) printf '%s' "$text" | awk -F'[_-]' '{for (i=1; i<=NF; i++) $i=toupper(substr($i,1,1)) substr($i,2)}1' OFS='' ;;
    *) printf '%s\n' "Invalid case. Use 'snake', 'kebab', or 'camel'." && return 1 ;;
  esac
}
## </editor-fold>

## <editor-fold defaultstate="collapsed" desc="Colors">
function _get_colors() {
  if [ -n "${NO_COLOR:-}" ] || [ -n "${ANSI_COLORS_DISABLED:-}" ]; then
    export nocolor="" bold="" nobold="" underline="" nounderline="" red="" green="" yellow="" blue="" magenta="" cyan=""
    return 0
  fi

  if ! tput colors &>/dev/null || [ "$(tput colors)" -lt 8 ]; then
    export nocolor="" bold="" nobold="" underline="" nounderline="" red="" green="" yellow="" blue="" magenta="" cyan=""
    return 0
  fi

  export nocolor="\033[0m" bold="\033[1m" nobold="\033[22m"
  export underline="\033[4m" nounderline="\033[24m"
  export red="\033[31m" green="\033[32m" yellow="\033[33m" blue="\033[34m" magenta="\033[35m" cyan="\033[36m"
}
## </editor-fold>

log_shell() {
  local _type=${1:-info}
  local _message=${2:-}
  local _debug=${3:-}
  local _quiet=${4:-}
  local _verbose=${5:-}
  _debug="${_debug:-${DEBUG:-${_DEBUG:-false}}}"
  _quiet="${_quiet:-${QUIET:-${_QUIET:-false}}}"
  _verbose="${_verbose:-${VERBOSE:-${_VERBOSE:-false}}}"

  case $_type in
    question | _QUESTION | -q | -Q)
      if [[ ${_debug:-false} == "true" ]]; then
        printf '%b[QUESTION]%b ❓  %s: ' "${_NOTICE:-\033[0;35m}" "${_NC:-\033[0m}" "$_message" >&2
      fi
      ;;
    notice | _NOTICE | -n | -N)
      if [[ ${_debug:-false} == "true" ]] || [[ ${_quiet:-false} == "false" && ${_verbose:-false} == "true" ]]; then
        printf '%b[NOTICE]%b  %s\n' "${_NOTICE:-\033[0;35m}" "${_NC:-\033[0m}" "$_message" >&2
      fi
      ;;
    info | _INFO | -i | -I)
      if [[ ${_debug:-false} == "true" ]]; then
        printf '%b[INFO]%b  %s\n' "${_INFO:-\033[0;36m}" "${_NC:-\033[0m}" "$_message" >&2
      fi
      ;;
    warn | _WARN | -w | -W)
      if [[ ${_debug:-false} == "true" ]]; then
        printf '%b[WARN]%b  %s\n' "${_WARN:-\033[0;33m}" "${_NC:-\033[0m}" "$_message" >&2
      fi
      ;;
    error | _ERROR | -e | -E)
      printf '%b[ERROR]%b  %s\n' "${_ERROR:-\033[0;31m}" "${_NC:-\033[0m}" "$_message" >&2
      ;;
    success | _SUCCESS | -s | -S)
      printf '%b[SUCCESS]%b  %s\n' "${_SUCCESS:-\033[0;32m}" "${_NC:-\033[0m}" "$_message" >&2
      ;;
    fatal | _FATAL | -f | -F)
      printf '%b[FATAL]%b 💀  %s\n' "${_FATAL:-\033[0;41m}" "${_NC:-\033[0m}" "Exiting due to fatal error: $_message" >&2
      clear_build_artifacts || true
      clear_script_cache || true

      # shellcheck disable=SC2317
      killall --process-group --signal=KILL make 2>/dev/null || true

      exit 1
      ;;
    separator | _SEPARATOR | hr | -hr | -HR | line)
      # if [[ "${_debug:-false}" != "true" ]]; then
      local _columns=${COLUMNS:-$(tput cols || echo 80)}
      local _margin=$((_columns - (_columns / 2)))
      _message="${_message// /¬}"
      _message="$(printf '%b%s%b %*s' "${_TRACE:-\033[0;34m}" "${_message:-}" "${_NC:-\033[0m}" "$((_columns - ("${#_message}" + _margin)))" '')"
      _message="${_message// /\#}"
      _message="${_message//¬/ }"
      printf '%s\n' "${_message:-}" >&2
      # fi
      ;;
    debug | _DEBUG | -d | -D)
      if [[ ${_debug:-false} == "true" ]]; then
        printf '%b[DEBUG]%b  %s\n' "${_TRACE:-\033[0;34m}" "${_NC:-\033[0m}" "$_message" >&2
      fi
      ;;
    *)
      log "info" "$_message" "${_debug:-false}" || true
      ;;
  esac

  return 0
}

## <editor-fold defaultstate="collapsed" desc="Log">
function log() {
  local log_type="${1:-info}"
  local _log_msg="${2:-No message}"
  local _force_stdout="${3:-false}"
  local log_z_msg=()
  local color=""
  local output_target="stderr"

  if [[ ${_force_stdout:-false} == "true" ]]; then
    output_target="stdout"
  fi

  if command -v logz &>/dev/null; then
    case "$log_type" in
      error | fatal | warn* | alert | info* | debug | notice | success)
        log_type=$(printf '%s%s' '-l' "${log_type}")
        log_z_msg=('-o' "${output_target}" '--message' "${_log_msg:-No message}")
        ;;
      help | -h) log_type=$(printf '--help') ;;
      version | -v) log_type=$(printf '--version') ;;
      separator | _SEPARATOR | hr | -hr | -HR | line | -line | -LINE | --hr)
        log_shell hr
        return 0
        ;;
      *)
        log_type=$(printf '%s%s' '-l' "info")
        log_z_msg=('-o' "${output_target}" '--message' "${_log_msg:-No message}")
        ;;
    esac
    logz "${log_type:-help}" "${log_z_msg[@]}" && return 0 || true
  else
    case "$log_type" in
      error | fatal) color="${bold}${red}" ;;
      warn* | alert) color="${bold}${yellow}" ;;
      info* | debug | notice) color="${bold}${blue}" ;;
      success) color="${bold}${green}" ;;
      help | -h)
        log_type='-h'
        ;;
      separator | _SEPARATOR | hr | -hr | -HR | line | -line | -LINE | --hr)
        log_shell hr
        return 0
        ;;
      *)
        color="${nocolor}"
        _log_msg="${log_type}"
        log_type="info"
        ;;
    esac
    if [[ ${_force_stdout:-false} == "true" ]]; then
      printf '%b[%s]%b %s\n' "$color" "$log_type" "$nocolor" "${_log_msg}"
    else
      printf '%b[%s]%b %s\n' "$color" "$log_type" "$nocolor" "${_log_msg}" >&2
    fi
  fi
}

## </editor-fold>

_get_colors

# Expose only needed functions
export -f _camel_to_snake
export -f _camel_to_kebab
export -f _convert_case
export -f _get_colors
export -f log
export -f log_shell
