# Clientset

A simple clientset interface to the k8s api-server.

## Instructions

1. Install dependencies

        $ mkdir -p ${GOPATH}/src/k8s.io
        $ cd ${GOPATH}/src/k8s.io
        $ git clone https://github.com/kubernetes/code-generator.git
        $ git clone https://github.com/kubernetes/gengo.git

2. Generate the clientset

        Note: - the code-generator:
                - is not aware of Go modules
                - will read definitions from the pkg/apis/itemresource/v1 package in this (github.com/keith-cullen/clientset) module
                - place generated code in ${GOPATH}/src/github.com/keith-cullen/clientset
              - the generated code must be copied from $GOPATH/src/github.com/keith-cullen/clientset into this module

        $ cd clientset
        $ ${GOPATH}/src/k8s.io/code-generator/generate-groups.sh all "github.com/keith-cullen/clientset/pkg/client" "github.com/keith-cullen/clientset/pkg/apis" itemresource:v1
        $ cp -r ${GOPATH}/src/github.com/keith-cullen/clientset/pkg .

3. Create the CRD

        $ kubectl apply -f item_crd.yaml
        $ kubectl get ns
        $ kubectl get crds

4. Build and run

        edit run.sh

            KUBECONFIG=
            KUBERNETES_MASTER=
            KUBERNETES_SERVICE_HOST=
            KUBERNETES_SERVICE_PORT=

        $ go build
        $ ./run.sh

5. Show the CRD instance

        $ kubectl describe item item1 -n itemns

6. Clean the project

        $ ./clean.sh
