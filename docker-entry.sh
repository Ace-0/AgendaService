usage() {
    echo "docker run agenda-service (server | client | repl) args..."
}

if [ $# -eq 0 ]; then
    usage
else
    case $1 in
        server)
            shift
            /go/bin/service "$@"
            ;;
        client)
            shift
            /go/bin/cli "$@"
            ;;
        repl)
            exec /bin/bash
            ;;
        *)
            usage
            ;;
    esac
fi
