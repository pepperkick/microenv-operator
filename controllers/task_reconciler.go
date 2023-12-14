package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	certmanager "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/hibiken/asynq"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	awsec2 "pepperkick.com/microenv-operator/api/crossplane/aws-ec2"
	awsroute53 "pepperkick.com/microenv-operator/api/crossplane/aws-route53"
	"pepperkick.com/microenv-operator/api/v1alpha1"
	microenvv1alpha1 "pepperkick.com/microenv-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	k8sCLient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"strings"
)

const (
	TypeReconciler = "process:reconciler"
)

type ReconcilerPayload struct {
	Name string `json:"name"`
}

type ReconcilerProcessor struct {
	manager manager.Manager
}

func NewReconcilerTask(ref types.NamespacedName) (*asynq.Task, error) {
	payload, err := json.Marshal(ReconcilerPayload{
		Name: ref.Name,
	})
	if err != nil {
		return nil, err
	}

	id := strings.ReplaceAll(ref.Name, "-", "_")

	return asynq.NewTask(TypeReconciler, payload, asynq.TaskID(id)), nil
}

func (pr *ReconcilerProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ReconcilerPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	scheme := runtime.NewScheme()
	logger := log.FromContext(ctx).WithValues("name", p.Name, "type", t.Type())

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(microenvv1alpha1.AddToScheme(scheme))
	utilruntime.Must(awsec2.AddToScheme(scheme))
	utilruntime.Must(awsroute53.AddToScheme(scheme))
	utilruntime.Must(certmanager.AddToScheme(scheme))

	client, err := k8sCLient.New(ctrl.GetConfigOrDie(), k8sCLient.Options{Scheme: scheme})
	if err != nil {
		logger.Info("Failed to create client")
		return nil
	}

	cluster := &v1alpha1.Cluster{}
	err = client.Get(ctx, types.NamespacedName{Name: p.Name}, cluster)
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			logger.Error(err, "Failed to fetch cluster resource")
			return err
		}

		logger.Info("Cluster got deleted!")
		return nil
	}

	provider := &v1alpha1.Config{}
	providerName := cluster.Spec.Provider
	if strings.EqualFold(providerName, "") {
		providerName = "default"
	}
	err = client.Get(context.TODO(), types.NamespacedName{Name: providerName}, provider)
	if err != nil {
		logger.Error(err, "Failed to fetch provider config resource", "name", providerName)
		return err
	}

	process := &ClusterReconcilerProcess{
		Client: client,
		Scheme: scheme,
		log:    logger,
	}

	return process.ReconcileCluster(cluster, provider)
}

func NewReconcilerProcessor() *ReconcilerProcessor {
	return &ReconcilerProcessor{}
}
