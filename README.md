<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Keto-maester](#keto-maester)
  - [Prerequisites](#prerequisites)
  - [Design](#design)
  - [How to use it](#how-to-use-it)
    - [Command-line flags](#command-line-flags)
  - [Development](#development)
    - [Testing](#testing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Keto-maester


This project contains a Kubernetes controller that uses Custom Resources (CR) to manage Keto access control policies and roles. ORY Keto Maester watches for instances of `oryaccesscontrolpolicyrole.keto.ory.sh/v1alpha1` and `oryaccesscontrolpolicy.keto.ory.sh/v1alpha1` CRs and creates, updates, or deletes corresponding ORY access control policies and roles by communicating with ORY Keto's API.

View [sample ORY access control policy resources](config/samples) to learn more about the `oryaccesscontrolpolicyrole.keto.ory.sh/v1alpha1` and `oryaccesscontrolpolicy.keto.ory.sh/v1alpha1` CRs.

The project is based on [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder).

## Prerequisites

- recent version of Go language with support for modules (e.g: 1.12.6)
- make
- kubectl
- kustomize
- [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) for running tests
- [ginkgo](https://onsi.github.io/ginkgo/) for local integration testing
- access to K8s environment: minikube or a remote K8s cluster
- [mockery](https://github.com/vektra/mockery) to generate mocks for testing purposes

## Design

Take a look at [Design Readme](./docs/README.md).

## How to use it

- `make test` to run tests
- `make test-integration` to run integration tests
- `make install` to generate CRD file from go sources and install it on the cluster
- `export KETO_URL={KETO_SERVICE_URL} && make run` to run the controller

To deploy the controller, edit the value of the ```--keto-url``` argument in the [manager.yaml](config/manager/manager.yaml) file and run ```make deploy```.

### Command-line flags

| Name            | Required | Description                  | Default value | Example values                                       |
|-----------------|----------|------------------------------|---------------|------------------------------------------------------|
| **keto-url**    | yes      | ORY Keto's service address   | -             | ` ory-keto-api.ory.svc.cluster.local`                |
| **keto-port**   | no       | ORY Keto's service port      | `4444`        | `4444`                                               |

## Development

### Testing

Use mockery to generate mock types that implement existing interfaces. To generate a mock type for an interface, navigate to the directory containing that interface and run this command:
```
mockery -name={INTERFACE_NAME}
```
