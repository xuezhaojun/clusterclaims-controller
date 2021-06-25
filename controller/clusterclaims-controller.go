// Copyright Contributors to the Open Cluster Management project.

package controller

import (
	"context"
	"crypto/sha256"

	"github.com/go-logr/logr"
	mcv1 "github.com/open-cluster-management/api/cluster/v1"
	kacv1 "github.com/open-cluster-management/klusterlet-addon-controller/pkg/apis/agent/v1"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const CredentialHash = "credential-hash"
const ProviderTypeLabel = "cluster.open-cluster-management.io/type"
const copiedFromNamespaceLabel = "cluster.open-cluster-management.io/copiedFromNamespace"
const copiedFromNameLabel = "cluster.open-cluster-management.io/copiedFromSecretName"
const CredentialLabel = "cluster.open-cluster-management.io/credentials"

var hash = sha256.New()

// ProviderCredentialSecretReconciler reconciles a Provider secret
type ClusterClaimsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *ClusterClaimsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := r.Log.WithValues("ClusterClaimsReconciler", req.NamespacedName)

	var cc hivev1.ClusterClaim
	if err := r.Get(ctx, req.NamespacedName, &cc); err != nil {
		log.V(0).Info("Resource deleted")

		//TODO: cleaneUpManagedCluster( req.NamespacedName)

		return ctrl.Result{}, nil
	}

	target := cc.ClusterName
	log.V(0).Info("Initalize cluster: " + target + " if needed")

	// ManagedCluster
	var mc mcv1.ManagedCluster
	err := r.Get(ctx, types.NamespacedName{Namespace: cc.Namespace, Name: target}, &mc)
	if k8serrors.IsNotFound(err) {

		log.V(0).Info("Create a new ManagedCluster resource and KlusterletAddonConfig")
		mc.Name = target
		mc.Spec.HubAcceptsClient = true
		mc.ObjectMeta.Labels = map[string]string{"vendor": "OpenShift"}

		if err = r.Create(ctx, &mc, &client.CreateOptions{}); err != nil {

			log.V(2).Info("Could not create ManagedCluster resource: " + target)
			return ctrl.Result{}, err
		}

	} else if err != nil {

		log.V(1).Info("Error when attempting to retreive the ManagedCluster resource: " + target)
		return ctrl.Result{}, err
	}

	// KlusterletAddonConfig
	var kac kacv1.KlusterletAddonConfig
	err = r.Get(ctx, types.NamespacedName{Namespace: target, Name: target}, &kac)
	if k8serrors.IsNotFound(err) {

		log.V(0).Info("Create a new KlusterletAddonConfig resource")
		kac.Name = target
		kac.Namespace = target
		kac.Spec.ClusterName = target
		kac.Spec.ClusterNamespace = target
		kac.Spec.ApplicationManagerConfig.Enabled = true
		kac.Spec.CertPolicyControllerConfig.Enabled = true
		kac.Spec.IAMPolicyControllerConfig.Enabled = true
		kac.Spec.PolicyController.Enabled = true
		kac.Spec.SearchCollectorConfig.Enabled = true

		if err = r.Create(ctx, &kac, &client.CreateOptions{}); err != nil {

			log.V(2).Info("Could not create KlusterletAddonConfig resource: " + target)
			return ctrl.Result{}, err
		}

	} else if err != nil {

		log.V(1).Info("Error when attempting to retreive the KlusterletAddonConfig resource: " + target)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ClusterClaimsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Secret{}).WithEventFilter(predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
	}).WithOptions(controller.Options{
		MaxConcurrentReconciles: 1, // This is the default
	}).Complete(r)
}
