#/bin/bash
set -e

cleanup() {
  echo "Killing cosmos, PID: ${COSMOS}"
  kill -9 $COSMOS
  echo "Killing diagnostics, PID: ${DIAG}"
  kill -9 $DIAG
  echo "Killing mesos, PID: ${MESOS}"
  kill -9 $MESOS
}
trap cleanup SIGINT SIGHUP 

help() {
  echo "./test [--private-ip][-p] 10.0.0.1 [--public-ip][-pub][default:none] 52.43.22.235"
  exit 1
}

while [[ $# > 1 ]]; do
  key="$1"
  case $key in
      -priv|--private-ip)
      PRIV="$2"
      shift 
      ;;
      -pub|--public-ip)
      PUB="$2"
      shift 
      ;;
      *)
      help 
      ;;
  esac
  shift 
done

[[ $PRIV && ${PRIV-x} ]] || help 
[[ $PUB && ${PUB-x} ]] || help

echo "Forwarding $PRIV:1050 -> ${PUB}:1050"
ssh -i ~/.ssh/mesos_dev.pem -N -L 1050:$PRIV:1050 core@$PUB &
export DIAG=$!
echo "Diagnostics session PID ${DIAG}"

echo "Forwarding $PRIV:7070 -> ${PUB}:7070"
ssh -i ~/.ssh/mesos_dev.pem  -N -L 7070:$PRIV:7070 core@$PUB &
export COSMOS=$!
echo "Comsos session PID ${COSMOS}"

echo "Forwarding $PRIV:5050 -> ${PUB}:5050"
ssh -i ~/.ssh/mesos_dev.pem  -N -L 5050:$PRIV:5050 core@$PUB &
export MESOS=$!
echo "Mesos session PID ${MESOS}"

while :; do
  sleep 5
done
