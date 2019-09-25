package cmd

const (
	// https://github.com/fmiguelez/wait-for
	wait_for = `
#!/bin/sh

TIMEOUT=15
QUIET=0

CONNS=
HOST=
PORT=

echoerr() {
  if [ "$QUIET" -ne 1 ]; then printf "%s\n" "$*" 1>&2; fi
}

usage() {
  exitcode="$1"
  cat << USAGE >&2
Usage:
  $cmdname host:port [host2:port2 ... [hostN:portN]] [-t timeout] [-- command args]
  -q | --quiet                        Do not output any status messages
  -t TIMEOUT | --timeout=timeout      Timeout in seconds, zero for no timeout
  -- COMMAND ARGS                     Execute command with args after the test finishes
USAGE
  exit "$exitcode"
}

wait_for() {
  for i in ` + "`" + `seq $TIMEOUT` + "`" + ` ; do
    nc -z "$HOST" "$PORT" > /dev/null 2>&1

    result=$?
    if [ $result -eq 0 ] ; then
      return 0
    fi
    sleep 1
  done
  echo "Operation timed out waiting for $HOST:$PORT" >&2
  exit 1
}


while [ $# -gt 0 ]
do
  case "$1" in
    *:* )
    CONNS="$CONNS $1"
    shift 1
    ;;
    -q | --quiet)
    QUIET=1
    shift 1
    ;;
    -t)
    TIMEOUT="$2"
    if [ "$TIMEOUT" = "" ]; then break; fi
    shift 2
    ;;
    --timeout=*)
    TIMEOUT="${1#*=}"
    shift 1
    ;;
    --)
    shift
    break
    ;;
    --help)
    usage 0
    ;;
    *)
    echoerr "Unknown argument: $1"
    usage 1
    ;;
  esac
done

if [ "$CONNS" = "" ]; then
  echoerr "Error: you need to provide a host and port to test."
  usage 2
fi

result=

for conn in $CONNS
do
    HOST="$(printf "%s\n" "$conn"| cut -d : -f 1)"
    PORT="$(printf "%s\n" "$conn"| cut -d : -f 2)"

    wait_for "$@"
    result=$?
done

if [ $result -eq 0 -a $# -gt 0 ]
then
    exec "$@"
fi	
	`

	// https://github.com/eficode/wait-for
	wait_for_one = `
#!/bin/sh

TIMEOUT=15
QUIET=0

echoerr() {
  if [ "$QUIET" -ne 1 ]; then printf "%s\n" "$*" 1>&2; fi
}

usage() {
  exitcode="$1"
  cat << USAGE >&2
Usage:
  $cmdname host:port [-t timeout] [-- command args]
  -q | --quiet                        Do not output any status messages
  -t TIMEOUT | --timeout=timeout      Timeout in seconds, zero for no timeout
  -- COMMAND ARGS                     Execute command with args after the test finishes
USAGE
  exit "$exitcode"
}

wait_for() {
  for i in ` + "`" + `seq $TIMEOUT` + "`" + ` ; do
    nc -z "$HOST" "$PORT" > /dev/null 2>&1
    
    result=$?
    if [ $result -eq 0 ] ; then
      if [ $# -gt 0 ] ; then
        exec "$@"
      fi
      exit 0
    fi
    sleep 1
  done
  echo "Operation timed out" >&2
  exit 1
}

while [ $# -gt 0 ]
do
  case "$1" in
    *:* )
    HOST=$(printf "%s\n" "$1"| cut -d : -f 1)
    PORT=$(printf "%s\n" "$1"| cut -d : -f 2)
    shift 1
    ;;
    -q | --quiet)
    QUIET=1
    shift 1
    ;;
    -t)
    TIMEOUT="$2"
    if [ "$TIMEOUT" = "" ]; then break; fi
    shift 2
    ;;
    --timeout=*)
    TIMEOUT="${1#*=}"
    shift 1
    ;;
    --)
    shift
    break
    ;;
    --help)
    usage 0
    ;;
    *)
    echoerr "Unknown argument: $1"
    usage 1
    ;;
  esac
done

if [ "$HOST" = "" -o "$PORT" = "" ]; then
  echoerr "Error: you need to provide a host and port to test."
  usage 2
fi

wait_for "$@"
	`
)
