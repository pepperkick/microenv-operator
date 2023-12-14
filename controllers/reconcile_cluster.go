package controllers

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	types2 "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	sprig "github.com/go-task/slim-sprig"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"gopkg.in/yaml.v2"
	"io"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"pepperkick.com/microenv-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
	"text/template"
	"time"
)

type KinstInstance struct {
	Name       string      `yaml:"name"`
	DockerHost string      `yaml:"docker"`
	Nodes      []KinstNode `yaml:"nodes"`
}

type KinstNode struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels,omitempty"`
	Taints []corev1.Taint    `yaml:"taints,omitempty"`
}

func (r *ClusterReconcilerProcess) ReconcileKinstCluster(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) error {
	_ = r.UpdateStatusCondition(cluster, "Reconciled", metav1.ConditionFalse, "Reconciling", "Reconciling KINST Cluster...")
	r.log.Info("Reconciling kinst cluster...")

	managerContainerName := "menv-reconciler"
	managerImage := provider.Spec.UtilImage.Image

	_, cli, err := r.getDockerCli(cluster.Status.ManagerInstanceIp, 2375)
	if err != nil {
		r.log.Error(err, "Failed to get docker cli for manager EC2 instance")
		return err
	}

	r.log.Info("Created docker cli for manager EC2 instance", "version", cli.ClientVersion(), "host", cli.DaemonHost())

	imageList, err := cli.ImageList(context.TODO(), types2.ImageListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("reference", managerImage)),
	})
	if err != nil {
		r.log.Error(err, "Failed to check for manager image", "image", managerImage, "host", cli.DaemonHost())
		return err
	}

	if len(imageList) == 0 {
		r.log.Info("Pulling menv manager image...", "image", managerImage, "host", cli.DaemonHost())

		authConfig := types2.AuthConfig{
			Username: provider.Spec.UtilImage.RegistryUsername,
			Password: provider.Spec.UtilImage.RegistryPassword,
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			r.log.Error(err, "Failed to encode auth for manager image", "image", managerImage, "host", cli.DaemonHost())
			return err
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)

		_, err = cli.ImagePull(context.TODO(), managerImage, types2.ImagePullOptions{
			RegistryAuth: authStr,
		})
		if err != nil {
			r.log.Error(err, "Failed to pull menv manager image", "image", managerImage, "auth", authStr, "host", cli.DaemonHost())
			return err
		}

		r.log.Info("Waiting for manager image to download...", "image", managerImage, "host", cli.DaemonHost())
		err = r.retry(50, 6*time.Second, func(r *ClusterReconcilerProcess) error {
			imageList, err := cli.ImageList(context.TODO(), types2.ImageListOptions{
				All:     true,
				Filters: filters.NewArgs(filters.Arg("reference", managerImage)),
			})
			if err != nil {
				r.log.Error(err, "Failed to check for manager image", "image", managerImage, "host", cli.DaemonHost())
				return err
			}

			if len(imageList) > 0 {
				return nil
			}

			r.log.Info("Waiting for manager image to download..", "image", managerImage, "host", cli.DaemonHost())
			return errors.New("manager image not downloaded")
		})
		if err != nil {
			r.log.Error(err, "Timeout waiting for manager image to be downloaded", "image", managerImage, "host", cli.DaemonHost())
			return err
		}
	}

	// Create a manager container to manage the instance via docker to avoid using SSH or AWS SSM
	// Check if container already exists
	containerList, err := cli.ContainerList(context.TODO(), types2.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: managerContainerName,
		}),
	})
	if err != nil {
		return err
	}

	if len(containerList) != 0 {
		r.log.Info("Removing old menv manager container...", "container", containerList[0].ID)
		err := cli.ContainerRemove(context.TODO(), containerList[0].ID, types2.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			r.log.Error(err, "Failed to delete menv manager container", "container", containerList[0].ID)
			return err
		}
	}

	r.log.Info("Creating menv manager container...")
	newContainer, err := cli.ContainerCreate(context.TODO(),
		&container.Config{
			Cmd:        []string{"-c", "sleep infinity"},
			Image:      managerImage,
			Entrypoint: []string{"/bin/sh"},
		},
		&container.HostConfig{
			Privileged: true,
			Binds: []string{
				"/bin:/host/bin",
				"/var/run/docker.sock:/var/run/docker.sock",
				"/home/ec2-user/menv-cluster:/home/ec2-user/menv-cluster",
			},
		},
		&network.NetworkingConfig{},
		&ocispec.Platform{}, managerContainerName)
	if err != nil {
		r.log.Error(err, "Failed to create menv manager container", "response", newContainer)
		return err
	}

	r.log.Info("Starting menv manager container...")
	err = cli.ContainerStart(context.TODO(), newContainer.ID, types2.ContainerStartOptions{})
	if err != nil {
		r.log.Error(err, "Failed to start menv manager container", "containerId", newContainer.ID)
		return err
	}

	script, err := r.generateKinstScript(cluster, provider)
	if err != nil {
		r.log.Error(err, "Failed to generate setup script for Kinst")
		return err
	}

	r.log.Info("Generated setup script for Kinst", "script", script)

	var scriptEnvs []string
	if strings.EqualFold(cluster.Spec.Infrastructure.CertIssuer, "cert-manager") {
		r.log.Info("Fetching certificate...")
		certSecret := corev1.Secret{}
		err = r.Get(context.TODO(), r.getCertificateNamespaceName(cluster, provider), &certSecret)
		if err != nil {
			r.log.Error(err, "Failed to get certificate secret")
			return err
		}

		scriptEnvs = append(scriptEnvs, "SCRIPT_CERT_ISSUER="+cluster.Spec.Infrastructure.CertIssuer)
		scriptEnvs = append(scriptEnvs, "SCRIPT_CERT_CRT="+string(certSecret.Data["tls.crt"]))
		scriptEnvs = append(scriptEnvs, "SCRIPT_CERT_KEY="+string(certSecret.Data["tls.key"]))
	}

	r.log.Info("Creating exec script for container...")
	exec, err := cli.ContainerExecCreate(context.TODO(), newContainer.ID, types2.ExecConfig{
		Env: scriptEnvs,
		Cmd: []string{"/bin/sh", "-c", "echo " + script + " | base64 -d > /script.sh; chmod +x /script.sh; /script.sh > /script.log 2>&1"},
	})
	if err != nil {
		r.log.Error(err, "Failed to create setup script for menv manager container", "containerId", newContainer.ID, "response", exec)
		return err
	}

	r.log.Info("Executing script in container...")
	err = cli.ContainerExecStart(context.TODO(), exec.ID, types2.ExecStartCheck{})
	if err != nil {
		r.log.Error(err, "Failed to execute setup script for menv manager container", "containerId", newContainer.ID, "execId", exec.ID)
		return err
	}

	r.log.Info("Waiting for exec to complete...")
	inspect := types2.ContainerExecInspect{}
	err = r.retry(60, 6*time.Second, func(r *ClusterReconcilerProcess) error {
		inspect, err = cli.ContainerExecInspect(context.TODO(), exec.ID)
		if err != nil {
			return err
		}

		if !inspect.Running {
			return nil
		}

		r.log.Info("Waiting for exec to complete...", "inspect", inspect)
		return errors.New("exec command not complete")
	})
	if err != nil {
		r.log.Error(err, "Timeout waiting for exec to complete", "containerId", newContainer.ID, "execId", exec.ID, "execInfo", inspect)
		return err
	}

	if inspect.ExitCode != 0 {
		r.log.Error(err, "Setup script failed with non-zero exit code", "containerId", newContainer.ID, "execId", exec.ID, "execInfo", inspect)
		return errors.New("cluster setup script failed with non-zero exit code")
	}

	r.log.Info("Executed setup script for Kinst", "containerId", newContainer.ID, "execId", exec.ID, "execInfo", inspect)

	r.log.Info("Fetching kubeconfig file...", "container", newContainer.ID)
	tarFile, _, err := cli.CopyFromContainer(context.TODO(), newContainer.ID, "/home/ec2-user/menv-cluster/kubeconfig.external")
	if err != nil {
		r.log.Error(err, "Failed to fetch kubeconfig file from control plane container")
		return err
	}

	tr := tar.NewReader(tarFile)
	h, err := tr.Next()
	if err != nil {
		r.log.Error(err, "Failed to read tar file")
		return err
	}
	bs, err := io.ReadAll(tr)
	if err != nil {
		r.log.Error(err, "Failed to read file", "file", h.Name)
		return err
	}
	r.log.Info("Fetched kubeconfig file", "container", newContainer.ID)

	secretName := "kubeconfig-" + cluster.Name
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: provider.Spec.SystemNamespace,
		},
		Type: "Opaque",
		Data: map[string][]byte{
			"kubeconfig": bs,
		},
	}

	existingSecret := &corev1.Secret{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: provider.Spec.SystemNamespace}, existingSecret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = controllerutil.SetControllerReference(cluster, secret, r.Scheme)
			if err != nil {
				r.log.Error(err, "Failed to set controller reference on kubeconfig secret")
				return err
			}
			err = r.Create(context.TODO(), secret)
			if err != nil {
				r.log.Error(err, "Failed to create kubeconfig secret")
				return err
			}

			r.log.Info("Created kubeconfig secret", "secret", secret.Name)
		} else {
			r.log.Error(err, "Failed to get secret")
			return err
		}
	} else {
		secret.ObjectMeta = existingSecret.ObjectMeta

		err = controllerutil.SetControllerReference(cluster, secret, r.Scheme)
		if err != nil {
			r.log.Error(err, "Failed to set controller reference on kubeconfig secret")
			return err
		}
		err = r.Update(context.TODO(), secret)
		if err != nil {
			r.log.Error(err, "Failed to update kubeconfig secret")
			return err
		}
		r.log.Info("Updated kubeconfig secret", "secret", secret.Name)
	}

	return nil
}

func (r *ClusterReconcilerProcess) ReconcileClusterFeatures(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) error {
	managerContainerName := "menv-reconciler"

	_, cli, err := r.getDockerCli(cluster.Status.ManagerInstanceIp, 2375)
	if err != nil {
		r.log.Error(err, "Failed to get docker cli for manager EC2 instance")
		return err
	}

	// Create a manager container to manage the instance via docker to avoid using SSH or AWS SSM
	// Check if container already exists
	containerList, err := cli.ContainerList(context.TODO(), types2.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: managerContainerName,
		}),
	})
	if err != nil {
		return err
	}

	var managerContainer types2.Container
	if len(containerList) != 0 {
		managerContainer = containerList[0]
	}

	script, err := r.generateFeaturesScript(cluster, provider)
	if err != nil {
		r.log.Error(err, "Failed to generate features script for cluster")
		return err
	}

	var scriptEnvs []string
	if cluster.Spec.Features.InstallArgoWorkflow {
		wf := cluster.Spec.Features.ArgoWorkflow
		scriptEnvs = append(scriptEnvs, "SCRIPT_USAGE_INSTALL_ARGO_WORKFLOW=true")
		scriptEnvs = append(scriptEnvs, "SCRIPT_USAGE_ARGO_WORKFLOW="+wf)
	}

	r.log.Info("Creating features script for container...", "containerId", managerContainer.ID, "script", script, "envs", scriptEnvs)
	exec, err := cli.ContainerExecCreate(context.TODO(), managerContainer.ID, types2.ExecConfig{
		Env: scriptEnvs,
		Cmd: []string{"/bin/sh", "-c", "echo " + script + " | base64 -d > /features.sh; chmod +x /features.sh; /features.sh > /features.log 2>&1"},
	})
	if err != nil {
		r.log.Error(err, "Failed to create features script for menv manager container", "containerId", managerContainer.ID, "response", exec)
		return err
	}

	r.log.Info("Executing features script in container...", "containerId", managerContainer.ID)
	err = cli.ContainerExecStart(context.TODO(), exec.ID, types2.ExecStartCheck{})
	if err != nil {
		r.log.Error(err, "Failed to execute features script for menv manager container", "containerId", managerContainer.ID, "execId", exec.ID)
		return err
	}

	r.log.Info("Waiting for features exec to complete...")
	inspect := types2.ContainerExecInspect{}
	err = r.retry(60, 6*time.Second, func(r *ClusterReconcilerProcess) error {
		inspect, err = cli.ContainerExecInspect(context.TODO(), exec.ID)
		if err != nil {
			return err
		}

		if !inspect.Running {
			return nil
		}

		r.log.Info("Waiting for features exec to complete...", "inspect", inspect)
		return errors.New("features exec command not complete")
	})
	if err != nil {
		r.log.Error(err, "Timeout waiting for features exec to complete", "containerId", managerContainer.ID, "execId", exec.ID, "execInfo", inspect)
		return err
	}

	if inspect.ExitCode != 0 {
		r.log.Error(err, "Features script failed with non-zero exit code", "containerId", managerContainer.ID, "execId", exec.ID, "execInfo", inspect)
		return errors.New("cluster features script failed with non-zero exit code")
	}

	r.log.Info("Executed features script for cluster", "containerId", managerContainer.ID, "execId", exec.ID, "execInfo", inspect)

	return nil
}

// generateKinstScript function generates a setup script for kinst
func (r *ClusterReconcilerProcess) generateKinstScript(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) (string, error) {
	fileName := "instance-kinst.tmpl.sh"
	filePath := "./assets/" + fileName
	tpl, err := template.New(fileName).Funcs(sprig.FuncMap()).ParseFiles(filePath)
	if err != nil {
		return "", err
	}

	config, err := r.generateMicroenvConfig(cluster, provider, false)
	if err != nil {
		return "", err
	}

	var scriptBuf bytes.Buffer
	scriptData := map[string]any{
		"microenvConfig": config,
	}
	err = tpl.Execute(&scriptBuf, scriptData)
	if err != nil {
		return "", err
	}

	r.log.Info("Generated kinst setup script", "file", filePath, "content", scriptBuf.String())
	return base64.StdEncoding.EncodeToString(scriptBuf.Bytes()), nil
}

// generateFeaturesScript function generates a usage script for the k8s cluster
func (r *ClusterReconcilerProcess) generateFeaturesScript(cluster *v1alpha1.Cluster, provider *v1alpha1.Config) (string, error) {
	fileName := "instance-features.tmpl.sh"
	filePath := "./assets/" + fileName
	tpl, err := template.New(fileName).Funcs(sprig.FuncMap()).ParseFiles(filePath)
	if err != nil {
		return "", err
	}

	var scriptBuf bytes.Buffer
	scriptData := map[string]any{}
	err = tpl.Execute(&scriptBuf, scriptData)
	if err != nil {
		return "", err
	}

	r.log.Info("Generated cluster features script", "file", filePath, "content", scriptBuf.String())
	return base64.StdEncoding.EncodeToString(scriptBuf.Bytes()), nil
}

// generateMicroenvConfig function generates a config file for microenv
func (r *ClusterReconcilerProcess) generateMicroenvConfig(cluster *v1alpha1.Cluster, provider *v1alpha1.Config, ignoreNotFoundInstance bool) (string, error) {
	fileName := "microenv-config.tmpl.yaml"
	filePath := "./assets/" + fileName
	tpl, err := template.New(fileName).Funcs(sprig.FuncMap()).ParseFiles(filePath)
	if err != nil {
		r.log.Error(err, "Failed to parse file", "file", filePath)
		return "", err
	}

	var instances []KinstInstance

	if !ignoreNotFoundInstance {
		for _, instance := range cluster.Spec.Infrastructure.Instances {
			name, ec2, err := r.getCrossplaneInstance(cluster, instance)
			if err != nil {
				r.log.Error(err, "Failed to get EC2 instance", "name", name)
				return "", err
			}

			kinstInstance := KinstInstance{
				Name:       "menv-" + instance.Name,
				DockerHost: r.getInstancePrivateIp(ec2) + ":2375",
				Nodes: []KinstNode{
					{
						Name: instance.Name + "-worker",
					},
				},
			}

			if len(instance.Nodes) > 0 {
				var nodes []KinstNode
				for i, node := range instance.Nodes {
					nodes = append(nodes, KinstNode{
						Name:   fmt.Sprintf("%s-worker%d", instance.Name, i+1),
						Labels: node.Labels,
						Taints: node.Taints,
					})
				}
				kinstInstance.Nodes = nodes
			}

			instances = append(instances, kinstInstance)
		}
	}

	yamlString, err := yaml.Marshal(instances)
	if err != nil {
		r.log.Error(err, "Failed to marshal menv config", "file", filePath)
		return "", err
	}

	var configBuf bytes.Buffer
	configData := map[string]any{
		"cluster":   cluster,
		"provider":  provider,
		"instances": string(yamlString),
	}
	err = tpl.Execute(&configBuf, configData)
	if err != nil {
		r.log.Error(err, "Failed to template menv config", "file", filePath)
		return "", err
	}

	r.log.Info("Generated menv config", "file", filePath, "instances", instances, "yaml", yamlString, "content", configBuf.String())
	return configBuf.String(), nil
}
