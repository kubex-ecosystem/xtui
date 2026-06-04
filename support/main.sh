#!/usr/bin/env bash
# shellcheck disable=SC2015,SC1091

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

__get_output_tty() {
  if [[ -t 1 ]]; then
    echo '/dev/tty'
  else
    echo '&2'
  fi
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
      set -o nounset           # Treat unset variables as an error when substituting
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

_QUIET=${_QUIET:-${QUIET:-false}}
_DEBUG=${_DEBUG:-${DEBUG:-false}}
_HIDE_ABOUT=${_HIDE_ABOUT:-${HIDE_ABOUT:-false}}
_SCRIPT_DIR="$(dirname "${0}")"

__first "${_main_args[@]}" >&2 || {
  echo "Error: This script must be run directly, not sourced." >&2
  exit 1
}

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

# Load library files
_SCRIPT_DIR="$(cd "$(dirname "${0}")" && pwd)"
__source_script_if_needed "show_summary" "${_SCRIPT_DIR:-}/config.sh" || exit 1
__source_script_if_needed "apply_manifest" "${_SCRIPT_DIR:-}/apply_manifest.sh" || exit 1
__source_script_if_needed "get_current_shell" "${_SCRIPT_DIR:-}/utils.sh" || exit 1
__source_script_if_needed "what_platform" "${_SCRIPT_DIR:-}/platform.sh" || exit 1
__source_script_if_needed "check_dependencies" "${_SCRIPT_DIR:-}/validate.sh" || exit 1
__source_script_if_needed "detect_shell_rc" "${_SCRIPT_DIR:-}/install_funcs.sh" || exit 1
__source_script_if_needed "build_binary" "${_SCRIPT_DIR:-}/build.sh" || exit 1
__source_script_if_needed "optimize_media" "${_SCRIPT_DIR:-}/optimize_media.sh" || exit 1
__source_script_if_needed "run" "${_SCRIPT_DIR:-}/project_cmds.sh" || exit 1

# Initialize traps
set_trap "${_main_args[@]}"

clear_screen

__run_custom_scripts() {
  local _STAGE="${1:-post}"
  if test -d "${_SCRIPT_DIR:-}/${_STAGE:-}.d/"; then
    if ls -1A "${_SCRIPT_DIR:-}/${_STAGE:-}.d/"*.sh >/dev/null 2>&1; then
      local _CUSTOM_SCRIPTS=()
      local _print_stage_header=false

      # shellcheck disable=SC2011
      _CUSTOM_SCRIPTS=("$(ls -1A "${_SCRIPT_DIR:-}/${_STAGE:-}.d/"*.sh | xargs -I{} basename {} || true)")
      local _CUSTOM_SCRIPTS_LEN="${#_CUSTOM_SCRIPTS[@]}"

      if [[ $_CUSTOM_SCRIPTS_LEN -gt 0 ]]; then
        log info "${_CUSTOM_SCRIPTS_LEN} ${_STAGE} custom scripts found..."

        if [[ $_CUSTOM_SCRIPTS_LEN -gt 1 ]]; then
          _print_stage_header=true
        fi

        if [ ${_print_stage_header:-false} = true ]; then
          log hr "[BEGIN CUSTOM STAGE: ${_STAGE}] "
        fi

        for _CUSTOM_SCRIPT in "${_CUSTOM_SCRIPTS[@]}"; do
          if [[ -f "${_SCRIPT_DIR:-}/${_STAGE:-}.d/${_CUSTOM_SCRIPT:-}" ]]; then
            log hr "[STAGE: ${_STAGE} - START SCRIPT: $(basename "${_CUSTOM_SCRIPT:-}")] "
            log notice "Executing script: ${_CUSTOM_SCRIPT}"
            # Ensure the script is executable
            if [[ ! -x "${_SCRIPT_DIR:-}/${_STAGE:-}.d/${_CUSTOM_SCRIPT:-}" ]]; then
              log info "Making script executable: ${_CUSTOM_SCRIPT:-}"
              chmod +x "${_SCRIPT_DIR:-}/${_STAGE:-}.d/${_CUSTOM_SCRIPT:-}" || {
                log error "Failed to make script executable: ${_CUSTOM_SCRIPT:-}"
                log hr "[STAGE: ${_STAGE} - END SCRIPT: $(basename "${_CUSTOM_SCRIPT:-}")] "
                test ${_print_stage_header:-false} = true && log hr "[END CUSTOM STAGE: ${_STAGE}] "
                return 1
              }
              log notice "Made script executable: ${_CUSTOM_SCRIPT:-}"
            fi

            # Execute the script without passing build arguments
            "${_SCRIPT_DIR:-}/${_STAGE:-}.d/${_CUSTOM_SCRIPT:-}" || {
              log error "Script execution failed: ${_CUSTOM_SCRIPT:-}"
              log hr "[STAGE: ${_STAGE} - END SCRIPT: $(basename "${_CUSTOM_SCRIPT:-}")] "
              test ${_print_stage_header:-false} = true && log hr "[END CUSTOM STAGE: ${_STAGE}] "
              return 1
            }
            log success "Script executed successfully: ${_CUSTOM_SCRIPT:-}"
          else
            log warn "Script not found: ${_CUSTOM_SCRIPT:-}"
            log hr "[STAGE: ${_STAGE} - END SCRIPT: $(basename "${_CUSTOM_SCRIPT:-}")] "
            test ${_print_stage_header:-false} = true && log hr "[END CUSTOM STAGE: ${_STAGE}] "
            return 1
          fi

          log hr "[STAGE: ${_STAGE} - END SCRIPT: $(basename "${_CUSTOM_SCRIPT:-}")] "
        done

        test ${_print_stage_header:-false} = true && log hr "[END CUSTOM STAGE: ${_STAGE}] "

        return 0
      fi
    fi
  fi

  return 0
}

__main() {
  if ! what_platform; then
    log error "Platform could not be determined."
    return 1
  fi

  local _arrArgs=("${_main_args[@]}")
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

  local _force="${_FORCE:-${FORCE:-n}}"
  local _will_upx_pack_binary="${_WILL_UPX_PACK_BINARY:-${WILL_UPX_PACK_BINARY:-}}"
  local _root_dir="${_ROOT_DIR:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _cmd_path="${_CMD_PATH:-${CMD_PATH:-${_root_dir}/cmd}}"
  local _binary_name="${_BINARY_NAME:-${BINARY_NAME:-$(basename "${_cmd_path}" .go)}}"
  local _app_name="${_APP_NAME:-${APP_NAME:-$(basename "${_root_dir}")}}"
  local _version="${_VERSION:-${VERSION:-$(git describe --tags)}}"

  local _platform="${_PLATFORM:-${_CURRENT_PLATFORM:-}}"
  local _arch="${_ARCH:-${_CURRENT_ARCH:-}}"
  local _build_target="${_BUILD_TARGET:-${_platform}-${_arch}}"

  if [[ ${_platform_arg} == "__CROSS_COMPILE__" ]]; then
    _platform_arg=""
  fi

  case "${_command:-}" in
    # Help
    # Main command dispatcher
    help | HELP | -h | -H)
      log info "Help:"
      echo "Usage: make {build|build-dev|build-docs|clean|test|help}"
      echo "Commands:"
      echo "  build    - Compiles the binary for the specified platform and architecture."
      echo "  build-dev - Builds the binary in development mode (without compression)."

      echo "  build-docs - Builds the documentation for the project."
      echo "  serve-docs - Serves the generated documentation locally."
      echo "  pub-docs - Publishes the documentation to GitHub Pages."

      echo "  test     - Runs the tests for the project."
      echo "  clean    - Cleans up build artifacts."
      echo "  help     - Displays this help message."

      return 0
      ;;
    build-dev | BUILD-DEV | -bd | -BD | dev)
      log notice "Preparing to build the binary..."
      if ! validate_versions; then
        log error "Required dependencies are missing. Please install them and try again."
        return 1
      fi
      log notice "Running build command in development mode..."
      build_binary "${_platform_arg:-}" "${_arch_arg:-}" "${_force:-}" "${_will_upx_pack_binary:-false}"
      return 0
      ;;
    build | BUILD | -b | -B)
      # validate_versions
      log notice "Preparing to build the binary..."
      if ! validate_versions; then
        log error "Required dependencies are missing. Please install them and try again."
        return 1
      fi
      log notice "Running build command..."
      build_binary "${_platform_arg:-__CROSS_COMPILE__}" "${_arch_arg:-}" "${_force:-}" "${_will_upx_pack_binary:-true}"
      return 0
      ;;

    # CLEAN
    # Clean up build artifacts
    clear | clean | CLEAN | -c | -C)
      log info "Running clean command..."
      clean_artifacts || return 1
      log success "Clean completed successfully."
      ;;

    # TEST
    # Run tests for the project
    test | TEST | -t | -T)
      log info "Running test command..."
      if ! check_dependencies; then
        log error "Required dependencies are missing. Please install them and try again."
        return 1
      fi
      if ! go test -v "${_ROOT_DIR:-}/..."; then
        log error "Tests failed. Please check the output for details."
        return 1
      fi
      log success "All tests passed successfully."
      ;;

    # BUILD-DOCS
    # Build documentation for the project
    build-docs | BUILD-DOCS | -bdc | -BDC)
      log info "Generating Documentation..."

      cd "${_ROOT_DIR:-}/docs" || {
        log error "Failed to change directory to ${_ROOT_DIR:-}"
        return 1
      }

      # Validate uv
      if [[ -t 1 && ! ${NON_INTERACTIVE:-} && ! ${CI:-} && -e /dev/tty ]]; then
        if ! command -v uv >/dev/null 2>&1; then
          apt-get update && apt-get install -y uv
          if ! command -v uv >/dev/null 2>&1; then
            log error "The 'uv' tool is required to build documentation. Please install it and try again."
            return 1
          fi
        fi

        # Validate if .venv exists

        if [[ ! -d ".venv" ]]; then
          uv --no-progress --quiet venv
          . .venv/bin/activate
          uv --no-progress --quiet pip install -r "${_ROOT_DIR:-}/support/docs/requirements.txt"
        else
          . .venv/bin/activate
        fi
      fi

      # Generate the documentation
      mkdocs build -f "${_ROOT_DIR:-}/support/docs/mkdocs.yml" -d "${_ROOT_DIR:-}/dist/docs" -q || {
        log error "Failed to generate documentation."
        return 1
      }

      log success "Documentation generated successfully."
      ;;

    # SERVE-DOCS
    # Serve the generated documentation
    serve-docs | SERVE-DOCS | -sdc | -SDC)
      log info "Serving Documentation..."
      cd "${_ROOT_DIR:-}/docs" || {
        log error "Failed to change directory to ${_ROOT_DIR:-}/docs"
        return 1
      }

      # Validate uv
      if [[ -t 1 && ! ${NON_INTERACTIVE:-} && ! ${CI:-} && -e /dev/tty ]]; then
        if ! command -v uv >/dev/null 2>&1; then
          apt-get update && apt-get install -y uv
          if ! command -v uv >/dev/null 2>&1; then
            log error "The 'uv' tool is required to build documentation. Please install it and try again."
            return 1
          fi
        fi

        # Validate if .venv exists
        if [[ ! -d ".venv" ]]; then
          uv --no-progress --quiet venv
          . .venv/bin/activate
          uv --no-progress --quiet pip install -r "${_ROOT_DIR:-}/support/docs/requirements.txt"
        else
          . .venv/bin/activate
        fi
      fi

      mkdocs serve -a "0.0.0.0:8081" -f "${_ROOT_DIR:-}/support/docs/mkdocs.yml" --dirtyreload -q || {
        log error "Failed to serve documentation."
        return 1
      }

      log success "Documentation server successfully ran at http://localhost:8081/docs"
      ;;

    # PUB-DOCS
    # Publish the documentation to GitHub Pages
    pub-docs | PUB-DOCS | -pd | -PD)
      log info "Publishing Documentation..."
      cd "${_ROOT_DIR:-}/docs" || {
        log error "Failed to change directory to ${_ROOT_DIR:-}/docs"
        return 1
      }

      # Validate uv
      if [[ -t 1 && ! ${NON_INTERACTIVE:-} && ! ${CI:-} && -e /dev/tty ]]; then
        if ! command -v uv >/dev/null 2>&1; then
          apt-get update && apt-get install -y uv
          if ! command -v uv >/dev/null 2>&1; then
            log error "The 'uv' tool is required to build documentation. Please install it and try again."
            return 1
          fi
        fi

        # Validate if .venv exists
        if [[ ! -d ".venv" ]]; then
          uv --no-progress --quiet venv
          . .venv/bin/activate
          uv --no-progress --quiet pip install -r "${_ROOT_DIR:-}/support/docs/requirements.txt"
        else
          . .venv/bin/activate
        fi
      fi

      mkdocs gh-deploy -f "${_ROOT_DIR:-}/support/docs/mkdocs.yml" -d "${_ROOT_DIR:-}/dist/docs" --force --no-history -q || {
        log error "Failed to publish documentation."
        return 1
      }

      log success "Documentation published successfully."
      ;;

    # DEFAULT
    # Default command handler
    run | RUN | -run | -RUN)
      # run the server
      log info "Starting execution..."
      run "${_arrArgs[1]:-}" || {
        log error "Failed to run command."
        return 1
      }
      ;;
    status | STATUS | -s | -S)
      log info "Getting status..."
      status || {
        log error "Failed to get status."
        return 1
      }
      ;;
    logs | LOGS | -l | -L)
      log info "Getting logs..."
      logs || {
        log error "Failed to get logs."
        return 1
      }
      ;;
    stop | STOP | -k | -K)
      log info "Stopping server..."
      stop || {
        log error "Failed to stop server."
        return 1
      }
      ;;
    restart | RESTART | -restart | -RESTART)
      log info "Restarting server..."
      restart || {
        log error "Failed to restart server."
        return 1
      }
      ;;
    export-logs | EXPORT-LOGS | -el | -EL)
      log info "Exporting logs..."
      export_logs || {
        log error "Failed to export logs."
        return 1
      }
      ;;
    *) # Fallback
      log error "Invalid command: ${_arrArgs[0]:-}"
      echo "Usage: make {build|build-dev|install|build-docs|clean|test|help}"
      ;;
  esac
}

# Função para limpar artefatos de build
clean_artifacts() {
  log info "Cleaning up build artifacts..."
  local _platforms=("windows" "darwin" "linux")
  local _archs=("amd64" "386" "arm64")
  for _platform in "${_platforms[@]}"; do
    for _arch in "${_archs[@]}"; do
      local _output_name
      _output_name=$(printf '%s_%s_%s' "${_BINARY:-}" "${_platform:-}" "${_arch:-}")
      if [[ ${_platform:-} != "windows" ]]; then
        local _compress_name="${_output_name:-}.tar.gz"
      else
        _output_name="${_output_name:-}.exe"
        local _compress_name="${_BINARY:-}_${_platform:-}_${_arch:-}.zip"
      fi
      rm -f "${_output_name:-}" || true
      rm -f "${_compress_name:-}" || true
    done
  done
  log success "Build artifacts removed."
}

__secure_logic_main() {
  local _ws_name
  _ws_name="$(__secure_logic_sourced_name)"
  local _ws_name_val
  _ws_name_val=$(eval 'echo ${_ws_name}')
  if test "${_ws_name_val:-}" != "true"; then
    __main "${_main_args[@]}"
    return $?
  else
    # If the script is sourced, we export the functions
    log error "This script is not intended to be sourced."
    log error "Please run it directly."
    return 1
  fi
}

_show_info() {
  if ! what_platform; then
    log error "Platform could not be determined."
    return 1
  fi

  local _arrArgs=("${_main_args[@]}")

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

  local _force="${_FORCE:-${FORCE:-n}}"
  local _will_upx_pack_binary="${_WILL_UPX_PACK_BINARY:-${WILL_UPX_PACK_BINARY:-true}}"
  local _root_dir="${_ROOT_DIR:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _cmd_path="${_CMD_PATH:-${CMD_PATH:-${_root_dir}/cmd}}"
  local _binary_name="${_BINARY_NAME:-${BINARY_NAME:-$(basename "${_cmd_path}" .go)}}"
  local _app_name="${_APP_NAME:-${APP_NAME:-$(basename "${_root_dir}")}}"
  local _version="${_VERSION:-${VERSION:-$(git describe --tags)}}"

  local _platform="${_PLATFORM:-${_CURRENT_PLATFORM:-}}"
  local _arch="${_ARCH:-${_CURRENT_ARCH:-}}"
  local _build_target="${_BUILD_TARGET:-${_platform}-${_arch}}"

  if [[ ${_platform_arg} == "__CROSS_COMPILE__" ]]; then
    _platform_arg=""
  fi

  log notice "Command: ${_command:-}"
  log notice "Platform: $(_get_os_from_args "${_platform_arg:-$(uname -s | tr '[:upper:]' '[:lower:]')}")"
  log notice "Architecture: $(_get_arch_from_args "${_platform_arg:-$(uname -s | tr '[:upper:]' '[:lower:]')}" "${_arch_arg:-$(uname -m | tr '[:upper:]' '[:lower:]')}")"

  show_headers || log fatal "Failed to display headers."
}

main() {
  _show_info "${_main_args[@]}" || {
    log fatal "Failed to display process information."
  }

  # O padrão é ler os hooks, caso existam
  # irá executar-los, então definir a variável como true é opcional
  # porém para evitar a execução, só definir a variável como false
  # Ex:
  # export RUN_PRE_SCRIPTS=false && make build # ou make build-dev, etc...
  log debug "RUN_PRE_SCRIPTS: ${_RUN_PRE_SCRIPTS:-${RUN_PRE_SCRIPTS:-false}}"
  if [[ ${_RUN_PRE_SCRIPTS:-${RUN_PRE_SCRIPTS:-false}} != "false" ]]; then
    __run_custom_scripts "pre" "${_main_args[@]}" || {
      log error "pre-hook-build: $?"
      log fatal "Failed to execute pre-installation scripts."
    }
  fi

  __secure_logic_main "${_main_args[@]}" || {
    log fatal "Script execution failed."
  }

  # O padrão é ler os hooks, caso existam
  # irá executar-los, então definir a variável como true é opcional
  # porém para evitar a execução, só definir a variável como false
  # Ex:
  # export RUN_POST_SCRIPTS=false && make build # ou make build-dev, etc...
  log debug "RUN_POST_SCRIPTS: ${_RUN_POST_SCRIPTS:-${RUN_POST_SCRIPTS:-true}}"
  if [[ ${_RUN_POST_SCRIPTS:-${RUN_POST_SCRIPTS:-true}} != "false" ]]; then
    __run_custom_scripts "post" "${_main_args[@]}" || {
      log error "post-hook-build: $?"
      log fatal "Failed to execute post-installation scripts."
    }
  fi

  __secure_logic_elapsed_time="$(($(date +%s) - __secure_logic_init_timestamp))"

  if [[ ${MYNAME_VERBOSE:-true} != "true" || ${_DEBUG:-false} == "true" ]]; then
    log info "Script executed in ${__secure_logic_elapsed_time} seconds."
  fi
}

main "${_main_args[@]}"

# End of script logic
