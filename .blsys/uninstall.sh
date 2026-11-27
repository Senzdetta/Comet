function blsys::uninstall() {
    import "(
        universe
    )"

    while [[ ${#} -gt 0 ]]; do 
        case "${1}" in
            '--home') args::takevalue 'HOME' "${@}"; shift ;;
        esac
        shift
    done

    format::argsfmt --print-init
    format::argsfmt '--home' "${HOME}"
    format::argsfmt --print-endline

    command blsys --remove "${__app_path__}"
    command blsys --remove "${__bindir__}/${__appname__}"

    trap - ERR EXIT
}; readonly -f blsys::uninstall