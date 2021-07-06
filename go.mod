module github.com/nicholasjackson/istio-smi-controller

go 1.15

require (
	github.com/cucumber/godog v0.11.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/go-logr/logr v0.4.0
	github.com/hashicorp/go-hclog v0.16.1
	github.com/servicemeshinterface/smi-controller-sdk v0.0.0-20210706112634-f50376b590e8
	github.com/shipyard-run/shipyard v0.3.12
	github.com/stretchr/testify v1.7.0
	k8s.io/client-go v0.19.3
	k8s.io/kubectl v0.19.3
	sigs.k8s.io/controller-runtime v0.7.2
)
