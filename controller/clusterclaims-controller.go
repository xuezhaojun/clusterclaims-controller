// Copyright Contributors to the Open Cluster Management project.

package clusterclaims

import (
	"context"
	"crypto/sha256"

	"github.com/go-logr/logr"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	APIReader client.Reader
	Log       logr.Logger
	Scheme    *runtime.Scheme
}

func (r *ClusterClaimsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := r.Log.WithValues("ClusterClaimsReconciler", req.NamespacedName)

	var cc hivev1.ClusterClaim
	if err := r.Get(ctx, req.NamespacedName, &cc); err != nil {
		log.V(1).Info("Resource deleted")

		//TODO: cleaneUpManagedCluster( req.NamespacedName)

		return ctrl.Result{}, err
	}

	log.V(1).Info("Create managedCluster kind if needed")

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
