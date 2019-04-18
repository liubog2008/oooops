#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE}")/..
# defines absolute root path
ROOT_PATH=${ROOT_PATH:-$(cd ${SCRIPT_ROOT}; pwd -P)}
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo k8s.io/code-generator)}


OUTPUT_DIR=${ROOT_PATH}/_output
mkdir ${OUTPUT_DIR}

# Register function to be called on EXIT to remove generated binary.
function cleanup {
  rm -rf "${OUTPUT_DIR:-}"
}
trap cleanup EXIT

EXT_APIS_PKG=${ROOT_PATH#"$GOPATH/src/"}/pkg/apis
OUTPUT_PKG=${ROOT_PATH#"$GOPATH/src/"}/pkg/client
# apps:v1,v2 othergroup:v1alpha1,v1alpha2
GROUP_VERSIONS="flow:v1alpha1"

function join() { local IFS="$1"; shift; echo "$*"; }

EXT_APIS=()

for GVs in ${GROUP_VERSIONS}; do
  IFS=: read G Vs <<<"${GVs}"
  # enumerate versions
  for V in ${Vs//,/ }; do
    EXT_APIS+=("${EXT_APIS_PKG}/${G}/${V}")
  done
done

echo "Building register-gen"
REGISTER_GEN="${OUTPUT_DIR}/register-gen"
go build -o "${REGISTER_GEN}" ${CODEGEN_PKG}/cmd/register-gen

echo "Generating register func for ${GROUP_VERSIONS}"
${REGISTER_GEN} --input-dirs $(join , "${EXT_APIS[@]}")



echo "Building deepcopy-gen"
DEEPCOPY_GEN="${OUTPUT_DIR}/deepcopy-gen"
go build -o "${DEEPCOPY_GEN}" ${CODEGEN_PKG}/cmd/deepcopy-gen

echo "Generating deepcopy funcs for ${GROUP_VERSIONS}"
${DEEPCOPY_GEN} --input-dirs $(join , "${EXT_APIS[@]}") -O zz_generated.deepcopy --bounding-dirs ${EXT_APIS_PKG}



echo "Building defaulter-gen"
DEFAULTER_GEN="${OUTPUT_DIR}/defaulter-gen"
go build -o "${DEFAULTER_GEN}" ${CODEGEN_PKG}/cmd/defaulter-gen

echo "Generating defaulters for ${GROUP_VERSIONS}"
${DEFAULTER_GEN}  --input-dirs $(join , "${EXT_APIS[@]}") -O zz_generated.defaults



echo "Building client-gen"
CLIENT_GEN="${OUTPUT_DIR}/client-gen"
go build -o "${CLIENT_GEN}" ${CODEGEN_PKG}/cmd/client-gen

echo "Generating clientset for ${GROUP_VERSIONS} at ${OUTPUT_PKG}/clientset"
${CLIENT_GEN} --clientset-name clientset --input-base "" --input $(join , "${EXT_APIS[@]}") --output-package ${OUTPUT_PKG}



echo "Building lister-gen"
LISTER_GEN="${OUTPUT_DIR}/lister-gen"
go build -o "${LISTER_GEN}" ${CODEGEN_PKG}/cmd/lister-gen

echo "Generating listers for ${GROUP_VERSIONS} at ${OUTPUT_PKG}/listers"
${LISTER_GEN} --input-dirs $(join , "${EXT_APIS[@]}") --output-package ${OUTPUT_PKG}/listers



echo "Building informer-gen"
INFORMER_GEN="${OUTPUT_DIR}/informer-gen"
go build -o "${INFORMER_GEN}" ${CODEGEN_PKG}/cmd/informer-gen

echo "Generating informers for ${GROUP_VERSIONS} at ${OUTPUT_PKG}/informers"
${INFORMER_GEN} \
    --input-dirs $(join , "${EXT_APIS[@]}") \
    --versioned-clientset-package ${OUTPUT_PKG}/clientset \
    --single-directory \
    --listers-package ${OUTPUT_PKG}/listers \
    --output-package ${OUTPUT_PKG}/informers
