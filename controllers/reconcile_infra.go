package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	certmanager "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	sprig "github.com/go-task/slim-sprig"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	awsec2 "pepperkick.com/microenv-operator/api/crossplane/aws-ec2"
	awsroute53 "pepperkick.com/microenv-operator/api/crossplane/aws-route53"
	"pepperkick.com/microenv-operator/api/v1alpha1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
	"text/template"
	"time"
)

// ReconcileInfrastructure will ensure that the infrastructure for the cluster is in desired state
func (r *ClusterReconcilerProcess) ReconcileInfrastructure(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) error {
	r.log.Info("Reconciling infrastructure...")

	_ = r.UpdateStatusCondition(cluster, "Reconciled", metav1.ConditionFalse, "Reconciling", "Reconciling Instances...")
	err := r.ReconcileInfrastructureInstances(cluster, provider)
	if err != nil {
		r.log.Error(err, "Failed to reconcile infrastructure instances")
		return err
	}

	// Reconcile domain record if config values are present
	if !strings.EqualFold(provider.Spec.Aws.Route53HostedZone, "") && !strings.EqualFold(provider.Spec.Aws.IngressDomain, "") {
		_ = r.UpdateStatusCondition(cluster, "Reconciled", metav1.ConditionFalse, "Reconciling", "Reconciling DNS Entries...")
		err = r.ReconcileInfrastructureDomain(cluster, provider)
		if err != nil {
			r.log.Error(err, "Failed to reconcile infrastructure domain")
			return err
		}
	}

	return nil
}

// ReconcileInfrastructureInstances will ensure that the instances for the cluster is in desired state
func (r *ClusterReconcilerProcess) ReconcileInfrastructureInstances(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) error {
	r.log.Info("Reconciling infrastructure instances...")

	for i, config := range cluster.Spec.Infrastructure.Instances {
		if i == 0 {
			config.DockerSwarmRole = "manager"
		} else {
			config.DockerSwarmRole = "worker"
		}

		name := r.getEc2InstanceName(cluster, config)
		newInstance, err := r.getEc2Instance(cluster, provider, config)
		if err != nil {
			r.log.Error(err, "Failed to create EC2 instance resource", "instance", name)
			return err
		}

		existingInstance := &awsec2.Instance{}
		err = r.Client.Get(context.TODO(), types.NamespacedName{Name: name}, existingInstance)
		if err != nil {
			if k8serrors.IsNotFound(err) {
				err = controllerutil.SetControllerReference(cluster, newInstance, r.Scheme)
				if err != nil {
					r.log.Error(err, "Failed to set controller reference on EC2 instance", "instance", name)
					return err
				}
				err = r.Client.Create(context.TODO(), newInstance)
				if err != nil {
					r.log.Error(err, "Failed to create EC2 instance", "instance", name)
					return err
				}
			} else {
				return err
			}
			r.log.Info("EC2 Instance created!")
		} else if !reflect.DeepEqual(newInstance.Spec, existingInstance.Spec) {
			newInstance.ObjectMeta = existingInstance.ObjectMeta

			if !cluster.Spec.Infrastructure.AllowInstanceUpdateForUserData {
				newInstance.Spec.ForProvider.UserDataBase64 = existingInstance.Spec.ForProvider.UserDataBase64
				newInstance.Spec.ForProvider.UserData = existingInstance.Spec.ForProvider.UserData
			}

			err = controllerutil.SetControllerReference(cluster, newInstance, r.Scheme)
			if err != nil {
				r.log.Error(err, "Failed to set controller reference on EC2 instance", "instance", name)
				return err
			}
			err = r.Client.Update(context.TODO(), newInstance)
			if err != nil {
				r.log.Error(err, "Failed to update EC2 instance", "instance", name)
				return err
			}
			r.log.Info("EC2 Instance updated!", "instance", name)
		}
	}

	r.log.Info("Waiting for crossplane to pickup any changes...")
	time.Sleep(5 * time.Second)

	// Wait for EC2 instances to be ready and synced
	for index, config := range cluster.Spec.Infrastructure.Instances {
		name := r.getEc2InstanceName(cluster, config)

		r.log.Info("Waiting for EC2 Instance to be ready...", "instance", name)
		err := r.retry(10, 6*time.Second, func(r *ClusterReconcilerProcess) error {
			existingInstance := &awsec2.Instance{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{Name: name}, existingInstance)
			if err != nil {
				r.log.Error(err, "Failed to fetch EC2 instance", "instance", name)
				return err
			}

			readyFlag := false
			syncedFlag := false
			for _, condition := range existingInstance.Status.Conditions {
				if condition.Type == v1.TypeReady && condition.Status == corev1.ConditionTrue {
					readyFlag = true
				}
				if condition.Type == v1.TypeSynced && condition.Status == corev1.ConditionTrue {
					syncedFlag = true
				}
			}

			if readyFlag && syncedFlag {
				// If this is the first instance, then it is the manager instance
				if index == 0 {
					if existingInstance.Status.AtProvider.PrivateIP != nil {
						cluster.Status.ManagerInstanceIp = *existingInstance.Status.AtProvider.PrivateIP
					} else if existingInstance.Status.AtProvider.PrivateDNS != nil {
						cluster.Status.ManagerInstanceIp = r.extractPrivateIpFromDns(existingInstance)
					}
					if existingInstance.Status.AtProvider.PublicIP != nil {
						cluster.Status.ManagerInstancePublicIp = *existingInstance.Status.AtProvider.PublicIP
					}
					err := r.Status().Update(context.TODO(), cluster)
					if err != nil {
						return err
					}
				}

				return nil
			}

			r.log.Info("Waiting for EC2 Instance to be ready...", "instance", name, "ready", readyFlag, "synced", syncedFlag)
			return errors.New("ec2 instance ready and sync condition not true")
		})
		if err != nil {
			r.log.Error(err, "Timeout waiting for EC2 instance to be ready", "instance", name)
			return err
		}

		r.log.Info("EC2 Instance is ready and synced!", "instance", name)
	}

	return nil
}

func (r *ClusterReconcilerProcess) ReconcileInfrastructureDomain(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) error {
	r.log.Info("Reconciling infrastructure domain...")

	name := r.getRoute53RecordName(cluster)
	newInstance, err := r.getRoute53Record(cluster, provider)
	if err != nil {
		r.log.Error(err, "Failed to create Route53 Record resource", "instance", name)
		return err
	}

	existingInstance := &awsroute53.Record{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: name}, existingInstance)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = controllerutil.SetControllerReference(cluster, newInstance, r.Scheme)
			if err != nil {
				r.log.Error(err, "Failed to set controller reference on Route53 Record", "instance", name)
				return err
			}
			err = r.Client.Create(context.TODO(), newInstance)
			if err != nil {
				r.log.Error(err, "Failed to create Route53 Record", "instance", name)
				return err
			}
		} else {
			return err
		}
		r.log.Info("Route53 Record created!")
	} else if !reflect.DeepEqual(newInstance.Spec, existingInstance.Spec) {
		newInstance.ObjectMeta = existingInstance.ObjectMeta

		err = controllerutil.SetControllerReference(cluster, newInstance, r.Scheme)
		if err != nil {
			r.log.Error(err, "Failed to set controller reference on Route53 Record", "instance", name)
			return err
		}
		err = r.Client.Update(context.TODO(), newInstance)
		if err != nil {
			r.log.Error(err, "Failed to update Route53 Record", "instance", name)
			return err
		}
		r.log.Info("Route53 Record updated!", "instance", name)
	}

	r.log.Info("Waiting for crossplane to pickup any changes...")
	time.Sleep(5 * time.Second)

	r.log.Info("Waiting for Route53 Record to be ready...", "instance", name)
	err = r.retry(30, 6*time.Second, func(r *ClusterReconcilerProcess) error {
		existingInstance := &awsroute53.Record{}
		err := r.Client.Get(context.TODO(), types.NamespacedName{Name: name}, existingInstance)
		if err != nil {
			r.log.Error(err, "Failed to fetch Route53 Record", "instance", name)
			return err
		}

		readyFlag := false
		syncedFlag := false
		for _, condition := range existingInstance.Status.Conditions {
			if condition.Type == v1.TypeReady && condition.Status == corev1.ConditionTrue {
				readyFlag = true
			}
			if condition.Type == v1.TypeSynced && condition.Status == corev1.ConditionTrue {
				syncedFlag = true
			}
		}

		if readyFlag && syncedFlag {
			return nil
		}

		r.log.Info("Waiting for Route53 Record to be ready...", "instance", name, "ready", readyFlag, "synced", syncedFlag)
		return errors.New("route53 record sync condition not true")
	})
	if err != nil {
		r.log.Error(err, "Timeout waiting for Route53 Record to be ready", "instance", name)
		return err
	}

	cluster.Status.ClusterIngressDomain = r.getRoute53RecordDomain(cluster, provider)
	err = r.Status().Update(context.TODO(), cluster)
	if err != nil {
		return err
	}

	r.log.Info("Route53 Record is ready and synced!", "instance", name)

	return nil
}

func (r *ClusterReconcilerProcess) ReconcileCertificate(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) error {
	r.log.Info("Reconciling certificate...")

	nsname := r.getCertificateNamespaceName(cluster, provider)
	newInstance, err := r.getCertificate(cluster, provider)
	if err != nil {
		r.log.Error(err, "Failed to create Certificate resource", "instance", nsname.String())
		return err
	}

	existingInstance := &certmanager.Certificate{}
	err = r.Client.Get(context.TODO(), nsname, existingInstance)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = controllerutil.SetControllerReference(cluster, newInstance, r.Scheme)
			if err != nil {
				r.log.Error(err, "Failed to set controller reference on Certificate", "instance", nsname.String())
				return err
			}
			err = r.Client.Create(context.TODO(), newInstance)
			if err != nil {
				r.log.Error(err, "Failed to create Certificate", "instance", nsname.String())
				return err
			}
		} else {
			return err
		}
		r.log.Info("Certificate created!")
	} else if !reflect.DeepEqual(newInstance.Spec, existingInstance.Spec) {
		newInstance.ObjectMeta = existingInstance.ObjectMeta

		err = controllerutil.SetControllerReference(cluster, newInstance, r.Scheme)
		if err != nil {
			r.log.Error(err, "Failed to set controller reference on Certificate", "instance", nsname.String())
			return err
		}
		err = r.Client.Update(context.TODO(), newInstance)
		if err != nil {
			r.log.Error(err, "Failed to update Certificate", "instance", nsname.String())
			return err
		}
		r.log.Info("Certificate updated!", "instance", nsname.String())
	}

	r.log.Info("Waiting for cert-manager to pickup any changes...")
	time.Sleep(5 * time.Second)

	r.log.Info("Waiting for Certificate to be ready...", "instance", nsname.String())
	err = r.retry(10, 6*time.Second, func(r *ClusterReconcilerProcess) error {
		existingInstance := &certmanager.Certificate{}
		err := r.Client.Get(context.TODO(), nsname, existingInstance)
		if err != nil {
			r.log.Error(err, "Failed to fetch Certificate", "instance", nsname.String())
			return err
		}

		readyFlag := false
		for _, condition := range existingInstance.Status.Conditions {
			if condition.Type == certmanager.CertificateConditionReady && condition.Status == cmmeta.ConditionTrue {
				readyFlag = true
			}
		}

		if readyFlag {
			return nil
		}

		r.log.Info("Waiting for Certificate to be ready...", "instance", nsname.String(), "ready", readyFlag)
		return errors.New("certificate sync condition not true")
	})
	if err != nil {
		r.log.Error(err, "Timeout waiting for Certificate to be ready", "instance", nsname.String())
		return err
	}

	cluster.Status.CertificateSecret = r.getCertificateNamespaceName(cluster, provider).String()
	err = r.Status().Update(context.TODO(), cluster)
	if err != nil {
		return err
	}

	r.log.Info("Certificate is ready!", "instance", nsname.String())

	return nil
}

func (r *ClusterReconcilerProcess) getCrossplaneInstance(cluster *v1alpha1.Cluster, config v1alpha1.InstanceSpec) (string, *awsec2.Instance, error) {
	name := r.getEc2InstanceName(cluster, config)

	instance := &awsec2.Instance{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: name}, instance)
	if err != nil {
		return "", nil, err
	}

	return name, instance, nil
}

func (r *ClusterReconcilerProcess) retry(attempts int, sleep time.Duration, f func(r *ClusterReconcilerProcess) error) (err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			time.Sleep(sleep)
		}
		err = f(r)
		if err == nil {
			return nil
		}
	}
	r.log.Error(err, "Function failed with multiple retries", "attempts", attempts)
	return err
}

// getEc2Instance function generates a EC2 Instance CR with filled values
func (r *ClusterReconcilerProcess) getEc2Instance(cluster *v1alpha1.Cluster, provider *v1alpha1.Config, config v1alpha1.InstanceSpec) (*awsec2.Instance, error) {
	providerName := "default"
	if strings.EqualFold(config.Type, "") {
		config.Type = provider.Spec.Aws.InstanceType
	}

	if strings.EqualFold(config.Type, "") {
		config.Type = "m5.large"
	}

	if !strings.EqualFold(provider.Spec.ProviderName, "") {
		providerName = provider.Spec.ProviderName
	}

	name := r.getEc2InstanceName(cluster, config)
	script, err := r.generateEc2InstanceScript(cluster, provider)
	if err != nil {
		r.log.Error(err, "Failed to generate startup script for instance", "instance", name)
		return nil, err
	}

	trueFlag := true

	httpResponseHopLimit := 3.0
	httpEndpointEnabled := "enabled"
	httpTokensOptional := "optional"
	instance := &awsec2.Instance{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: awsec2.InstanceSpec{
			ResourceSpec: v1.ResourceSpec{
				ProviderConfigReference: &v1.Reference{
					Name: providerName,
				},
			},
			ForProvider: awsec2.InstanceParameters{
				AMI: &provider.Spec.Aws.BaseAmiId,
				RootBlockDevice: []awsec2.RootBlockDeviceParameters{
					config.RootVolume,
				},
				InstanceType: &config.Type,
				Region:       &provider.Spec.Aws.Region,
				SubnetID:     &provider.Spec.Aws.SubnetId,
				Tags: map[string]*string{
					"Name": &name,
				},
				VPCSecurityGroupIds: []*string{&provider.Spec.Aws.SecurityGroupId},
				UserDataBase64:      &script,
				IAMInstanceProfile:  &provider.Spec.Aws.IamInstanceProfileName,
				MetadataOptions: []awsec2.MetadataOptionsParameters{
					{
						HTTPPutResponseHopLimit: &httpResponseHopLimit,
						HTTPEndpoint:            &httpEndpointEnabled,
						HTTPTokens:              &httpTokensOptional,
					},
				},
				AssociatePublicIPAddress: &trueFlag,
			},
		},
		Status: awsec2.InstanceStatus{},
	}

	return instance, nil
}

// getEc2InstanceName function returns a name for EC2 instance
func (r *ClusterReconcilerProcess) getEc2InstanceName(cluster *v1alpha1.Cluster, config v1alpha1.InstanceSpec) string {
	name := "menv-" + cluster.Name + "-" + config.Name
	name = strings.TrimSuffix(name, "-")
	return name
}

// generateEc2InstanceScript function generates a startup script for EC2 instance
func (r *ClusterReconcilerProcess) generateEc2InstanceScript(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) (string, error) {
	fileName := "ec2-instance-script.tmpl.sh"
	filePath := "./assets/" + fileName
	tpl, err := template.New(fileName).Funcs(sprig.FuncMap()).ParseFiles(filePath)
	if err != nil {
		r.log.Error(err, "Failed to parse file", "file", filePath)
		return "", err
	}

	config, err := r.generateMicroenvConfig(cluster, provider, true)
	if err != nil {
		r.log.Error(err, "Failed to generate menv config", "file", filePath)
		return "", err
	}

	var scriptBuf bytes.Buffer
	scriptData := map[string]any{
		"microenvConfig":            config,
		"customInstanceSetupScript": provider.Spec.InstanceSetupScript,
	}
	err = tpl.Execute(&scriptBuf, scriptData)
	if err != nil {
		r.log.Error(err, "Failed to template file", "file", filePath)
		return "", err
	}

	r.log.Info("Generated ec2 startup script", "file", filePath, "content", scriptBuf.String())
	return base64.StdEncoding.EncodeToString(scriptBuf.Bytes()), nil
}

// getRoute53Record function generates a Route53 Record CR with filled values
func (r *ClusterReconcilerProcess) getRoute53Record(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) (*awsroute53.Record, error) {
	name := r.getRoute53RecordName(cluster)

	aRecord := "A"
	ttl := 30.0
	domain := "*." + r.getRoute53RecordDomain(cluster, provider)

	ip := ""
	if provider.Spec.Aws.UsePrivateIp && !strings.EqualFold(cluster.Status.ManagerInstanceIp, "") {
		ip = cluster.Status.ManagerInstanceIp
	} else {
		ip = cluster.Status.ManagerInstancePublicIp
	}

	providerName := "default"
	if !strings.EqualFold(provider.Spec.ProviderName, "") {
		providerName = provider.Spec.ProviderName
	}

	instance := &awsroute53.Record{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: awsroute53.RecordSpec{
			ResourceSpec: v1.ResourceSpec{
				ProviderConfigReference: &v1.Reference{
					Name: providerName,
				},
			},
			ForProvider: awsroute53.RecordParameters{
				Name: &domain,
				Records: []*string{
					&ip,
				},
				Region: &provider.Spec.Aws.Region,
				ZoneID: &provider.Spec.Aws.Route53HostedZone,
				Type:   &aRecord,
				TTL:    &ttl,
			},
		},
		Status: awsroute53.RecordStatus{},
	}

	return instance, nil
}

// getRoute53RecordName function returns a name for Route53 record
func (r *ClusterReconcilerProcess) getRoute53RecordName(cluster *v1alpha1.Cluster) string {
	return "menv-" + cluster.Name + "-domain"
}

// getRoute53RecordDomain function returns a domain for Route53 record
func (r *ClusterReconcilerProcess) getRoute53RecordDomain(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) string {
	return "env-" + cluster.Name + "." + provider.Spec.Aws.IngressDomain
}

// getCertificate function generates a Certificate Request CR with filled values
func (r *ClusterReconcilerProcess) getCertificate(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) (*certmanager.Certificate, error) {
	nsname := r.getCertificateNamespaceName(cluster, provider)

	duration := metav1.Duration{}
	duration.Duration, _ = time.ParseDuration("4320h")

	instance := &certmanager.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nsname.Name,
			Namespace: nsname.Namespace,
		},
		Spec: certmanager.CertificateSpec{
			Subject: &certmanager.X509Subject{
				OrganizationalUnits: []string{"microenv"},
			},
			CommonName: cluster.Status.ClusterIngressDomain,
			Duration:   &duration,
			DNSNames:   []string{cluster.Status.ClusterIngressDomain},
			SecretName: nsname.Name,
			IssuerRef: cmmeta.ObjectReference{
				Name:  provider.Spec.Aws.PcaCertIssuer,
				Kind:  "AWSPCAClusterIssuer",
				Group: "awspca.cert-manager.io",
			},
			Usages: []certmanager.KeyUsage{
				certmanager.UsageServerAuth,
				certmanager.UsageClientAuth,
			},
			PrivateKey: &certmanager.CertificatePrivateKey{
				Algorithm: certmanager.RSAKeyAlgorithm,
				Size:      2048,
			},
		},
		Status: certmanager.CertificateStatus{},
	}

	return instance, nil
}

// getCertificateName function returns a name for Certificate
func (r *ClusterReconcilerProcess) getCertificateNamespaceName(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) types.NamespacedName {
	return types.NamespacedName{Name: cluster.Status.ClusterIngressDomain, Namespace: provider.Spec.SystemNamespace}
}

func (r *ClusterReconcilerProcess) getInstancePrivateIp(instance *awsec2.Instance) string {
	if instance.Status.AtProvider.PrivateIP != nil {
		return *instance.Status.AtProvider.PrivateIP
	} else if instance.Status.AtProvider.PrivateDNS != nil {
		r.log.Info("WARNING: Private IP not available, extracting from Private DNS...")
		return r.extractPrivateIpFromDns(instance)
	} else {
		r.log.Info("WARNING: Private IP not available, using public IP...")
		return *instance.Status.AtProvider.PublicIP
	}
}

func (r *ClusterReconcilerProcess) extractPrivateIpFromDns(instance *awsec2.Instance) string {
	// ip-10-128-105-24.us-west-2.compute.internal
	// ip-10-128-105-24
	part := strings.Split(*instance.Status.AtProvider.PrivateDNS, ".")[0]
	// ip, 10, 128, 105, 24
	parts := strings.Split(part, "-")
	return fmt.Sprintf("%s.%s.%s.%s", parts[1], parts[2], parts[3], parts[4])
}
