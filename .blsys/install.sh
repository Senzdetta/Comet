function blsys::install() {
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

    package::finsdep
    path::ensure_path
    app::automove_app

    (
        cd "${__app_path__}"
        command blsys --runcmd "
            command go mod tidy
            command go build \
                -ldflags='-O3 -s -w' \
                -trimpath \
                -o ${__app_path__}/${__appname__} \
                .
        " \
        --print "Compiling: ${color_GG}${__appname__}${color_N}"
    )

    iofs::chmod "(
        0755|${__app_path__}/${__appname__}
    )"

    iofs::makelink "(
        default|${__app_path__}/${__appname__} -> ${__bindir__}/${__appname__}
    )"

    system::cmdcheck "${__appname__}"
    trap - ERR EXIT
}; readonly -f blsys::install