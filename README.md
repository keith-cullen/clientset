# Clientset

A simple clientset interface to the k8s api-server.

## Instructions

1. Install dependencies

        $ mkdir -p "${GOPATH}/src/k8s.io"
        $ cd "${GOPATH}/src/k8s.io"
        $ git clone https://github.com/kubernetes/code-generator.git
        $ git clone https://github.com/kubernetes/gengo.git

2. Generate the clientset

        $ cd clientset
        $ ./gen.sh

3. Create the CRD

        $ kubectl apply -f item_crd.yaml
        $ kubectl get ns
        $ kubectl get crds

4. Build and run

        $ go build
        edit run.sh
            KUBECONFIG=
            KUBERNETES_MASTER=
            KUBERNETES_SERVICE_HOST=
            KUBERNETES_SERVICE_PORT=
        $ ./run.sh

5. Show the CRD instance

        $ kubectl describe item item1 -n itemns

6. Clean the project

        $ ./clean.sh
