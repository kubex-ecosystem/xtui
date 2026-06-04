#!/usr/bin/env bash
# shellcheck disable=SC2155,SC2207

declare -a _ARR_FN=()
declare -a _tmp_final_fn_list=()

_iter_pid=""
_pb_tmp_file=""
_tot_exec_perc=0
_curr_exec_perc=0

__kbx_DmsShUGFnSrc_iterator() {
  for ((_fn_idx = 0; _fn_idx < _tot_exec_perc; _fn_idx++)); do
    _curr_exec_perc=$((_curr_exec_perc + 1))

    local _error=""

    local _fn_tmp_val=$(readlink -f "$(printf '/%s\n' "$(command -V "${_ARR_FN[$_fn_idx]}" 2>/dev/null | cut --delimiter='/' -f2-)" 2>/dev/null | grep -i '/' | grep -vE 'not|found|is an autoload shell function')") || {
      _error="$?"
      kbx_log error "Error: Could not get fn source for: ${_ARR_FN[$_fn_idx]:-}"
    }
    _fn_tmp_val="${_fn_tmp_val:-}"

    if test -n "${_fn_tmp_val:-}"; then
      _tmp_final_fn_list+=("${_fn_tmp_val}") || {
        _error="$?"
        kbx_log error "Error: Could not get fn source for: ${_ARR_FN[$_fn_idx]:-}"
      }
    fi

    # Append new line for the function's source code file
    printf '%s\n' "${_fn_tmp_val:-}" >>"${_pb_tmp_file}.tmp.fn" || {
      _error="$?"
      kbx_log error "Error: Could not write to temporary file for progress bar." true
      return 1
    }

    # Update the progress bar value (current iteration percentage)
    printf '%s\n' "${_curr_exec_perc}" >"${_pb_tmp_file}"
  done

  _tmp_final_fn_list=(
    "${_tmp_final_fn_list[@]}"
  )

  printf '# LIST_OF_FILES' >"${_pb_tmp_file}.fn" || {
    kbx_log error "Failed to create temporary file for progress bar." true
    return 1
  }

  if cat "${_pb_tmp_file}.tmp.fn" | sort -u >>"${_pb_tmp_file}.fn"; then
    rm -rf "${_pb_tmp_file}.tmp.fn" 2>/dev/null || true
  else
    kbx_log error "Error: Could not sort and uniq temporary file for progress bar." true
    return 1
  fi

  printf '\n# LIST_OF_FUNCTIONS\n' >>"${_pb_tmp_file}.fn" || {
    kbx_log error "Failed to create temporary file for progress bar." true
    return 1
  }

  printf '%s\n' "${_ARR_FN[@]}" | sort -u >>"${_pb_tmp_file}.fn" || {
    kbx_log error "Failed to create temporary file for progress bar." true
    return 1
  }

  return 0
}

__kbx_DmsShUGFnSrc_iter() {
  local _monitor_state
  _monitor_state="$(set +o | grep monitor 2>/dev/null || true)"
  set +m 2>/dev/null || true

  __kbx_DmsShUGFnSrc_iterator &

  _iter_pid=$! || {
    kbx_log error "Failed to start fn source iteration."
    return 1
  }

  disown "${_iter_pid}" 2>/dev/null || true

  eval "${_monitor_state}" 2>/dev/null || true

  if ! kill -CONT "${_iter_pid}" 2>/dev/null; then
    kbx_log error "Failed to start fn source iteration." true
    return 1
  fi

  return 0
}

__kbx_DmsShUGFnSrc() {
  local _fn_arr=()
  local _grep_str="${1:-}"
  local _tmp_cmd=(
    "declare -f"
    "grep ' () {'"
  )

  if test -n "${_grep_str:-}"; then
    _tmp_cmd+=("grep -iE '${_grep_str}'")
  fi

  _tmp_cmd+=(
    "awk '{ print \$1 }'"
    "grep -vE 'zsh|zle|nvm_|_x_|__kbx_|___error___|source_script_if_needed'"
  )

  local _declare_cmd="$(printf '%s | ' "${_tmp_cmd[@]}" | sed 's/ | $//')"
  _ARR_FN=($(eval "${_declare_cmd}" || true))

  _tot_exec_perc=${#_ARR_FN[@]}
  _curr_exec_perc=0

  _pb_tmp_file="$(mktemp -t XTUIPB_FN.XXXXXX)" || {
    kbx_log error "Failed to create temporary file for progress bar." true
    return 1
  }

  printf '0\n' >"${_pb_tmp_file}" || {
    kbx_log error "Failed to write to temporary file for progress bar." true
    return 1
  }
  touch "${_pb_tmp_file}.tmp.fn" || {
    kbx_log error "Failed to create temporary file for progress bar." true
    return 1
  }

  # Start the iteration in the background for the progress bar to stay in the foreground
  # to ensure perfect synchronization of the rendered bar.
  __kbx_DmsShUGFnSrc_iter || {
    kbx_log error "Failed to start fn source iteration." true
    return 1
  }

  # Start the XTUI in the background, directing the output to the user's TTY
  # to avoid mixing with this function's stdout (return).
  # We assume that the 'xtui' binary is in the PATH or uses the $_XTUI_BIN variable
  ${_XTUI_BIN:-xtui} f pb -f "${_pb_tmp_file}" -t "Processando shell functions..." -T "$((_tot_exec_perc))" -w --timeout 120 -D || {
    kbx_log error "Failed to run xtui progress-bar." true
    return 1
  }

  if ! [[ ${_iter_pid:-} =~ ^[0-9]+$ ]]; then
    kbx_log error "Failed to start fn source iteration." true
    return 1
  fi

  if ps -p ${_iter_pid:-} >/dev/null 2>&1; then
    wait ${_iter_pid} || {
      kbx_log error "Failed to wait for fn source iteration." true
      return 1
    }
  fi

  return 0
}

function get_fn_source() {
  __kbx_DmsShUGFnSrc "${@}"

  local _final_fn_list=("${_tmp_final_fn_list[@]}")
  local _tot_fn="${#_tmp_final_fn_list[@]}"
  kbx_log info "Total functions: ${_tot_exec_perc}" true
  kbx_log info "List file: ${_pb_tmp_file}.fn" true
  local _fn_content="$(cat "${_pb_tmp_file}.fn" || true)"
  kbx_log info "Fn content: " true
  printf '%s\n' "${_fn_content}" >&2

  # rm -rf "${_pb_tmp_file}" "${_pb_tmp_file}.fn" 2>/dev/null || true

  return 0
}
