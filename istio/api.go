package istio

import (
	"context"

	"github.com/go-logr/logr"
	specsv1alpha4 "github.com/servicemeshinterface/smi-controller-sdk/apis/specs/v1alpha4"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type API struct {
	client Client
}

func New(c Client) *API {
	return &API{c}
}

func (l *API) UpsertHTTPRouteGroup(
	ctx context.Context,
	r client.Client,
	log logr.Logger,
	tt *specsv1alpha4.HTTPRouteGroup) (ctrl.Result, error) {

	log.Info("UpdateHTTPRouteGroup called", "api", "v1alpha4", "target", tt)

	return ctrl.Result{}, nil
}

func (l *API) DeleteHTTPRouteGroup(
	ctx context.Context,
	r client.Client,
	log logr.Logger,
	tt *specsv1alpha4.HTTPRouteGroup) (ctrl.Result, error) {

	log.Info("DeleteHTTPRouteGroup called", "api", "v1alpha4", "target", tt)

	return ctrl.Result{}, nil
}

func (l *API) UpsertTCPRoute(
	ctx context.Context,
	r client.Client,
	log logr.Logger,
	tt *specsv1alpha4.TCPRoute) (ctrl.Result, error) {

	log.Info("UpdateTCPRoute called", "api", "v1alpha4", "target", tt)

	return ctrl.Result{}, nil
}

func (l *API) DeleteTCPRoute(
	ctx context.Context,
	r client.Client,
	log logr.Logger,
	tt *specsv1alpha4.TCPRoute) (ctrl.Result, error) {

	log.Info("DeleteTCPRoute called", "api", "v1alpha4", "target", tt)

	return ctrl.Result{}, nil
}

func (l *API) UpsertUDPRoute(
	ctx context.Context,
	r client.Client,
	log logr.Logger,
	tt *specsv1alpha4.UDPRoute) (ctrl.Result, error) {

	log.Info("UpdateUDPRoute called", "api", "v1alpha4", "target", tt)

	return ctrl.Result{}, nil
}

func (l *API) DeleteUDPRoute(
	ctx context.Context,
	r client.Client,
	log logr.Logger,
	tt *specsv1alpha4.UDPRoute) (ctrl.Result, error) {

	log.Info("DeleteUDPRoute called", "api", "v1alpha4", "target", tt)

	return ctrl.Result{}, nil
}
