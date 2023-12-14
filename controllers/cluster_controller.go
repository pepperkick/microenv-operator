package controllers

import (
	"context"
	"fmt"
	logger "github.com/go-logr/logr"
	"github.com/hibiken/asynq"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"pepperkick.com/microenv-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Queue  *asynq.Client
	log    logger.Logger
}

type ClusterReconcilerProcess struct {
	client.Client
	log    logger.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=cert-manager.io,resources=certificates,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=microenv.pepperkick.com,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=microenv.pepperkick.com,resources=clusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=microenv.pepperkick.com,resources=clusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=microenv.pepperkick.com,resources=configs,verbs=get;list;watch
//+kubebuilder:rbac:groups=ec2.aws.upbound.io,resources=instances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ec2.aws.upbound.io,resources=instances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ec2.aws.upbound.io,resources=instances/finalizers,verbs=update
//+kubebuilder:rbac:groups=route53.aws.upbound.io,resources=records,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=route53.aws.upbound.io,resources=records/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=route53.aws.upbound.io,resources=records/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.log = log.FromContext(ctx)
	r.log.Info("Reconciling cluster...", "name", req.NamespacedName.String())

	task, err := NewReconcilerTask(req.NamespacedName)
	if err != nil {
		r.log.Error(err, "Failed to create reconcile task", "name", req.NamespacedName.String())
		return ctrl.Result{}, err
	}

	info, err := r.Queue.EnqueueContext(ctx, task)
	if err != nil {
		r.log.Info("Task already queued, ignoring...")
		return ctrl.Result{}, nil
	}

	r.log.Info("Enqueued task", "ID", info.ID, "Queue", info.Queue, "name", req.NamespacedName.String())

	return ctrl.Result{}, nil
}

// ReconcileCluster will ensure that the entire cluster is in desired state
func (r *ClusterReconcilerProcess) ReconcileCluster(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) error {
	_ = r.UpdateStatusCondition(cluster, "Ready", v1.ConditionFalse, "Reconciling", "Reconciling Infrastructure...")
	_ = r.UpdateStatusCondition(cluster, "Reconciled", v1.ConditionFalse, "Reconciling", "Reconciling Infrastructure...")
	_ = r.UpdateStatusCondition(cluster, "Infrastructure", v1.ConditionFalse, "Reconciling", "")

	err := r.ReconcileInfrastructure(cluster, provider)
	if err != nil {
		r.log.Error(err, "Failed to reconcile infrastructure")
		_ = r.UpdateStatusCondition(cluster, "Reconciled", v1.ConditionFalse, "Failed",
			fmt.Sprintf("Failed to reconcile infrastructure: %s", err))
		_ = r.UpdateStatusCondition(cluster, "Infrastructure", v1.ConditionFalse, "Failed",
			fmt.Sprintf("Failed to reconcile infrastructure: %s", err))
		return err
	}

	if strings.EqualFold(cluster.Spec.Infrastructure.CertIssuer, "cert-manager") {
		_ = r.UpdateStatusCondition(cluster, "Ready", v1.ConditionFalse, "Reconciling", "Reconciling Certificates...")
		_ = r.UpdateStatusCondition(cluster, "Reconciled", v1.ConditionFalse, "Reconciling", "Reconciling Certificates...")
		err = r.ReconcileCertificate(cluster, provider)
		if err != nil {
			r.log.Error(err, "Failed to reconcile certificate")
			_ = r.UpdateStatusCondition(cluster, "Reconciled", v1.ConditionFalse, "Failed",
				fmt.Sprintf("Failed to reconcile certificate: %s", err))
			_ = r.UpdateStatusCondition(cluster, "Infrastructure", v1.ConditionFalse, "Failed",
				fmt.Sprintf("Failed to reconcile certificate: %s", err))
			return err
		}
	}

	err = r.UpdateStatusCondition(cluster, "Infrastructure", v1.ConditionTrue, "Available", "")
	if err != nil {
		r.log.Error(err, "Failed to update infrastructure status")
		return err
	}

	_ = r.UpdateStatusCondition(cluster, "Ready", v1.ConditionFalse, "Reconciling", "Reconciling Swarm...")
	_ = r.UpdateStatusCondition(cluster, "Reconciled", v1.ConditionFalse, "Reconciling", "Reconciling Swarm...")
	_ = r.UpdateStatusCondition(cluster, "Cluster", v1.ConditionFalse, "Reconciling", "")
	err = r.ReconcileDockerSwarm(cluster)
	if err != nil {
		r.log.Error(err, "Failed to reconcile docker swarm")
		_ = r.UpdateStatusCondition(cluster, "Reconciled", v1.ConditionFalse, "Failed",
			fmt.Sprintf("Failed to reconcile docker swarm: %s", err))
		_ = r.UpdateStatusCondition(cluster, "Cluster", v1.ConditionFalse, "Failed",
			fmt.Sprintf("Failed to reconcile docker swarm: %s", err))
		return err
	}

	_ = r.UpdateStatusCondition(cluster, "Ready", v1.ConditionFalse, "Reconciling", "Reconciling Cluster...")
	_ = r.UpdateStatusCondition(cluster, "Reconciled", v1.ConditionFalse, "Reconciling", "Reconciling Cluster...")
	err = r.ReconcileKinstCluster(cluster, provider)
	if err != nil {
		r.log.Error(err, "Failed to reconcile kinst cluster")
		_ = r.UpdateStatusCondition(cluster, "Reconciled", v1.ConditionFalse, "Failed",
			fmt.Sprintf("Failed to reconcile kinst cluster: %s", err))
		_ = r.UpdateStatusCondition(cluster, "Cluster", v1.ConditionFalse, "Failed",
			fmt.Sprintf("Failed to reconcile kinst cluster: %s", err))
		return err
	}

	err = r.UpdateStatusCondition(cluster, "Cluster", v1.ConditionTrue, "Available", "")
	if err != nil {
		r.log.Error(err, "Failed to update cluster status")
		return err
	}

	err = r.UpdateStatusCondition(cluster, "Ready", v1.ConditionTrue, "Available", "")
	if err != nil {
		r.log.Error(err, "Failed to update ready status")
		return err
	}

	_ = r.UpdateStatusCondition(cluster, "Reconciled", v1.ConditionFalse, "Reconciling", "Applying Features...")
	_ = r.UpdateStatusCondition(cluster, "Features", v1.ConditionFalse, "Reconciling", "")
	err = r.ReconcileClusterFeatures(cluster, provider)
	if err != nil {
		r.log.Error(err, "Failed to reconcile cluster features")
		_ = r.UpdateStatusCondition(cluster, "Reconciled", v1.ConditionFalse, "Failed",
			fmt.Sprintf("Failed to reconcile cluster features: %s", err))
		_ = r.UpdateStatusCondition(cluster, "Features", v1.ConditionFalse, "Failed",
			fmt.Sprintf("Failed to reconcile cluster features: %s", err))
		return err
	}

	err = r.UpdateStatusCondition(cluster, "Features", v1.ConditionTrue, "Available", "")
	if err != nil {
		r.log.Error(err, "Failed to update features status")
		return err
	}

	r.log.Info("Reconciled cluster!")

	_ = r.UpdateStatusCondition(cluster, "Reconciled", v1.ConditionTrue, "Complete", "Cluster is Ready!")

	return nil
}

// UpdateStatusCondition updates the status condition of the cluster
func (r *ClusterReconcilerProcess) UpdateStatusCondition(cluster *v1alpha1.Cluster, conditionType string, status v1.ConditionStatus, reason string, message string) error {
	meta.SetStatusCondition(&cluster.Status.Conditions, v1.Condition{
		Type:    conditionType,
		Status:  status,
		Reason:  reason,
		Message: message,
	})

	return r.Status().Update(context.TODO(), cluster)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Cluster{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(r)
}
