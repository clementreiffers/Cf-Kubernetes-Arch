# cf-worker-architecture

Set of Kubernetes Operators which manage Cloudflare Workers Deployment into Kubernetes Cluster.

## Getting Started

You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/workerbundle:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/workerbundle:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## Architecture

```mermaid
stateDiagram-v2
    state First {
        JobBuilder
        Registry
    }
    state Second {
        WorkerBundle
        Deployment
    }
    FakeCfApi --> WorkerVersion : create
    WorkerVersion --> WorkerRelease : create or update
    WorkerRelease --> JobBuilder : create
    JobBuilder --> Registry : push
    WorkerAccount --> WorkerBundle : create
    WorkerBundle --> Deployment : create
    Registry --> Deployment : pull
    JobBuilder --> WorkerBundle : update while finished
```

### WorkerAccount

Once the architecture deployed into Kubernetes Cluster, you can create your first account into the architecture :


> **Note** : you need to install wrangler by running `npm i -g wrangler`, run `wrangler login` and get the account ID by
> running `wrangler whoami` and keep it for this step.

you can now apply this resource into kubernetes :

```yaml
apiVersion: api.cf-worker/v1
kind: WorkerAccount
metadata:
  name: YOUR-WRANGLER-ACCOUNT-ID # accounts
  labels:
    accounts: YOUR-WRANGLER-ACCOUNT-ID
spec:
  workerBundleName: "workerBundleName"
  workerReleaseSelector: 
    matchLabels: 
      accounts: YOUR-WRANGLER-ACCOUNT-ID
  podTemplate:
    imagePullSecret: "insert-secret-here"
```

## License

Copyright 2023 clementreiffers.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

