package istio

import (
	"context"
	"errors"
	"fmt"
	"sync"

	accessv1alpha3 "github.com/servicemeshinterface/smi-controller-sdk/apis/access/v1alpha3"
	splitv1alpha4 "github.com/servicemeshinterface/smi-controller-sdk/apis/split/v1alpha4"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	securityv1beta1 "istio.io/api/security/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	sv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	smiLabels = map[string]string{
		"app.kubernetes.io/managed-by": "istio-smi-controller-sdk",
		"app.kubernetes.io/created-by": "istio-smi-controller-sdk",
	}
)

// Client allows the creation of Istio objects from their SMI counterparts
type Client interface {
	CreateVirtualService(ctx context.Context, r client.Writer, ts *splitv1alpha4.TrafficSplit) error
	DeleteVirtualService(ctx context.Context, r client.Writer, ts *splitv1alpha4.TrafficSplit) error

	CreateAuthorizationPolicy(ctx context.Context, r client.Writer, tt *accessv1alpha3.TrafficTarget) error
	DeleteAuthorizationPolicy(ctx context.Context, r client.Writer, tt *accessv1alpha3.TrafficTarget) error
}

// IstioClient is a concrete implementation of Client
type IstioClient struct {
	once sync.Once
}

// lazyRegisterIstioTypesToScheme lazily registers the Istio types to the
// controllers client. This allows us to use the default client that is created by the
// controller instead of manually creating one.
//
// Since the client is not available until the first time the controller is called
// this function must be called from every Client interface method.
// sync.Once ensures that types are only registered once.
func (ic *IstioClient) lazyRegisterIstioTypesToScheme(c client.Client) {
	ic.once.Do(func() {
		// register types to the scheme
		err := v1beta1.AddToScheme(c.Scheme())
		if err != nil {
			panic(err)
		}
	})
}

// CreateVirtualService creates an Istio VirtualService from the given TrafficSplit
func (ic *IstioClient) CreateVirtualService(ctx context.Context, r client.Writer, ts *splitv1alpha4.TrafficSplit) error {
	ic.lazyRegisterIstioTypesToScheme(r.(client.Client))

	vs := &v1beta1.VirtualService{}

	vs.ObjectMeta.Name = ts.ObjectMeta.Name
	vs.ObjectMeta.Namespace = ts.ObjectMeta.Namespace
	vs.ObjectMeta.SetLabels(smiLabels) // TODO - copy over labels from TrafficSplit while stripping out Istio related ones

	vs.Spec.Hosts = []string{ts.Spec.Service}

	// if matches is not present then create as a HTTPRoute
	httpRoute := &networkingv1beta1.HTTPRoute{}
	httpRoute.Route = []*networkingv1beta1.HTTPRouteDestination{}

	for _, be := range ts.Spec.Backends {
		hrd := &networkingv1beta1.HTTPRouteDestination{}
		hrd.Destination = &networkingv1beta1.Destination{
			Host:   ts.Spec.Service,
			Subset: be.Service,
		}
		hrd.Weight = int32(be.Weight)

		httpRoute.Route = append(httpRoute.Route, hrd)
	}

	vs.Spec.Http = []*networkingv1beta1.HTTPRoute{httpRoute}

	return r.Create(ctx, vs, &client.CreateOptions{})
}

func (ic *IstioClient) DeleteVirtualService(ctx context.Context, r client.Writer, ts *splitv1alpha4.TrafficSplit) error {
	ic.lazyRegisterIstioTypesToScheme(r.(client.Client))

	vs := &v1beta1.VirtualService{}

	// client only needs name and namespace to perform delete operation
	vs.ObjectMeta.Name = ts.ObjectMeta.Name
	vs.ObjectMeta.Namespace = ts.ObjectMeta.Namespace

	return r.DeleteAllOf(ctx, vs, client.MatchingLabels(smiLabels))
}

// TODO - For this to be effective, we'd need to find some way to deny all requests by default (assuming that's the design of the SMI spec)
func (ic *IstioClient) CreateAuthorizationPolicy(ctx context.Context, r client.Writer, tt *accessv1alpha3.TrafficTarget) error {
	ic.lazyRegisterIstioTypesToScheme(r.(client.Client))
	// l := log.FromContext(ctx)

	ap := &sv1beta1.AuthorizationPolicy{}

	ap.ObjectMeta.Name = tt.Name
	ap.ObjectMeta.Namespace = tt.Namespace
	ap.ObjectMeta.SetLabels(smiLabels) // TODO - copy over labels from TrafficTarget while stripping out Istio related ones

	ap.Spec.Action = securityv1beta1.AuthorizationPolicy_ALLOW // THe SMI spec is additive ALLOWs only

	// Source configuration (only principals since that's what the SMI spec supports)
	var principals []string
	for _, source := range tt.Spec.Sources {
		if source.Kind != "ServiceAccount" {
			// Istio only supports service account principals; not Groups
			continue
		}
		principals = append(principals, fmt.Sprintf("cluster.local/ns/%s/sa/%s", source.Namespace, source.Name))
	}

	// TrafficTarget essentially represents a single Istio "Rule" so we only add one here
	ap.Spec.Rules = append(ap.Spec.Rules, &securityv1beta1.Rule{
		From: []*securityv1beta1.Rule_From{
			{
				Source: &securityv1beta1.Source{
					Principals: principals,
				},
			},
		},
		To: []*securityv1beta1.Rule_To{
			{
				Operation: &securityv1beta1.Operation{},
			},
		},
	})

	return errors.New("not implemented")
}

func (ic *IstioClient) DeleteAuthorizationPolicy(ctx context.Context, w client.Writer, tt *accessv1alpha3.TrafficTarget) error {
	ic.lazyRegisterIstioTypesToScheme(w.(client.Client))

	ap := &sv1beta1.AuthorizationPolicy{}

	// client only needs name and namespace to perform delete operation
	ap.ObjectMeta.Name = tt.ObjectMeta.Name
	ap.ObjectMeta.Namespace = tt.ObjectMeta.Namespace

	return w.DeleteAllOf(ctx, ap, client.MatchingLabels(smiLabels))
}
