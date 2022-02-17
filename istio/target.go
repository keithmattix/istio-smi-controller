package istio

import (
	"context"

	"github.com/go-logr/logr"
	accessv1alpha3 "github.com/servicemeshinterface/smi-controller-sdk/apis/access/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (l *API) UpsertTrafficTarget(
	ctx context.Context,
	c client.Client,
	log logr.Logger,
	tt *accessv1alpha3.TrafficTarget,
) (ctrl.Result, error) {

	log.Info("UpsertTrafficTarget called", "api", "v1alpha3", "target", tt)

	l.client.CreateAuthorizationPolicy(ctx, c, tt)

	return ctrl.Result{}, nil
}

func (l *API) DeleteTrafficTarget(
	ctx context.Context,
	c client.Client,
	log logr.Logger,
	tt *accessv1alpha3.TrafficTarget,
) (ctrl.Result, error) {

	log.Info("DeleteTrafficTarget called", "api", "v1alpha3", "target", tt)
	l.client.DeleteAuthorizationPolicy(ctx, c, tt)

	return ctrl.Result{}, nil
}
