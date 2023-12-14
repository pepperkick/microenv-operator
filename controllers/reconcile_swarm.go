package controllers

import (
	"context"
	"errors"
	types2 "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
	awsec2 "pepperkick.com/microenv-operator/api/crossplane/aws-ec2"
	"pepperkick.com/microenv-operator/api/v1alpha1"
	"strconv"
	"strings"
	"time"
)

func (r *ClusterReconcilerProcess) ReconcileDockerSwarm(cluster *v1alpha1.Cluster) error {
	r.log.Info("Reconciling docker swarm...")

	dockerSwarmToken := ""
	dockerSwarmManagerIp := ""

	// Wait for Docker ports to be open
	for i, config := range cluster.Spec.Infrastructure.Instances {
		name, instance, err := r.getCrossplaneInstance(cluster, config)
		if err != nil {
			r.log.Error(err, "Failed to fetch EC2 instance", "instance", name)
			return err
		}

		dockerHost, cli, err := r.getDockerCliForInstance(instance)
		if err != nil {
			r.log.Error(err, "Failed to get docker cli for EC2 instance", "instance", name)
			return err
		}

		r.log.Info("Waiting for docker to be up in the instance...", "instance", name, "dockerHost", cli.DaemonHost())
		err = r.retry(60, 6*time.Second, func(r *ClusterReconcilerProcess) error {
			list, err := cli.ContainerList(context.TODO(), types2.ContainerListOptions{})
			if err != nil {
				return err
			}

			if list != nil {
				return nil
			}

			r.log.Info("Waiting for docker to be up in the instance...", "instance", name, "dockerHost", cli.DaemonHost(), "list", list)
			return errors.New("failed to list docker containers")
		})
		if err != nil {
			r.log.Error(err, "Timeout waiting for docker to be ready", "instance", name, "dockerHost", dockerHost)
			return err
		}

		inspect, err := cli.SwarmInspect(context.TODO())
		if err != nil {
			if !strings.Contains(err.Error(), "This node is not a swarm manager") {
				r.log.Error(err, "Failed to fetch docker swarm status", "instance", name, "dockerHost", dockerHost)
				return err
			}
		}

		if inspect.ID != "" {
			r.log.Info("Docker swarm for instance is in desired state.", "instance", name, "dockerHost", dockerHost)
			dockerSwarmManagerIp = r.getInstancePrivateIp(instance) + ":2377"
			dockerSwarmToken = inspect.JoinTokens.Worker

			continue
		}

		if i == 0 {
			init, err := cli.SwarmInit(context.TODO(), swarm.InitRequest{
				ListenAddr:      "0.0.0.0:2377",
				ForceNewCluster: false,
			})
			if err != nil {
				r.log.Error(err, "Failed to initialize docker swarm", "instance", name, "dockerHost", dockerHost)
				return err
			}

			inspect, err := cli.SwarmInspect(context.TODO())
			if err != nil {
				r.log.Error(err, "Failed to fetch docker swarm status", "instance", name, "dockerHost", dockerHost)
				return err
			}

			dockerSwarmManagerIp = r.getInstancePrivateIp(instance) + ":2377"
			dockerSwarmToken = inspect.JoinTokens.Worker

			r.log.Info("Docker swarm initialized", "instance", name, "dockerHost", dockerHost, "response", init, "managerIp", dockerSwarmManagerIp, "joinToken", dockerSwarmToken)
		} else {
			err := cli.SwarmJoin(context.TODO(), swarm.JoinRequest{
				ListenAddr:  "0.0.0.0:2377",
				RemoteAddrs: []string{dockerSwarmManagerIp},
				JoinToken:   dockerSwarmToken,
			})
			if err != nil {
				if strings.Contains(err.Error(), "This node is already part of a swarm.") {
					r.log.Info("Docker swarm for instance is in desired state.", "instance", name, "dockerHost", dockerHost)
					continue
				}

				r.log.Error(err, "Failed to join docker swarm", "instance", name, "dockerHost", dockerHost, "managerIp", dockerSwarmManagerIp, "joinToken", dockerSwarmToken)
				return err
			}

			r.log.Info("Docker swarm joined", "instance", name, "dockerHost", dockerHost, "managerIp", dockerSwarmManagerIp, "joinToken", dockerSwarmToken)
		}
	}

	return nil
}

func (r *ClusterReconcilerProcess) ReconcileDockerNetwork(cluster *v1alpha1.Cluster) error {
	r.log.Info("Reconciling docker network...")

	overlayNetworkName := "menv"
	overlayNetworkId := ""

	// Ensure overlay docker network is present
	managerName, managerInstance, err := r.getCrossplaneInstance(cluster, cluster.Spec.Infrastructure.Instances[0])
	if err != nil {
		r.log.Error(err, "Failed to fetch EC2 instance", "instance", managerName)
		return err
	}

	_, cli, err := r.getDockerCliForInstance(managerInstance)
	if err != nil {
		r.log.Error(err, "Failed to get docker cli for EC2 instance", "instance", managerName)
		return err
	}

	r.log.Info("Created docker cli for EC2 instance", "version", cli.ClientVersion())

	list, err := cli.NetworkList(context.TODO(), types2.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: overlayNetworkName,
		}, filters.KeyValuePair{
			Key:   "driver",
			Value: "overlay",
		}),
	})
	if err != nil {
		r.log.Error(err, "Failed to fetch docker networks")
		return err
	}

	if len(list) == 0 {
		r.log.Info("Creating docker overlay network...")
		// Docker network does not exist for cluster, create it
		create, err := cli.NetworkCreate(context.TODO(), overlayNetworkName, types2.NetworkCreate{
			CheckDuplicate: true,
			Driver:         "overlay",
			Attachable:     true,
		})
		if err != nil {
			r.log.Error(err, "Failed to create docker overlay network")
			return err
		}

		overlayNetworkId = create.ID
	} else {
		overlayNetworkId = list[0].ID
	}

	if strings.EqualFold(overlayNetworkId, "") {
		r.log.Error(err, "Failed to find docker overlay network")
		return errors.New("could not find docker overlay network")
	}

	r.log.Info("Found docker overlay network", "id", overlayNetworkId)

	return nil
}

func (r *ClusterReconcilerProcess) getDockerCliForInstance(instance *awsec2.Instance) (string, *docker.Client, error) {
	return r.getDockerCli(r.getInstancePrivateIp(instance), 2375)
}

func (r *ClusterReconcilerProcess) getDockerCli(ip string, port int) (string, *docker.Client, error) {
	dockerHost := "tcp://" + ip + ":" + strconv.Itoa(port)

	cli, err := docker.NewClientWithOpts(
		docker.WithHost(dockerHost),
		docker.WithTimeout(15*time.Second),
	)
	if err != nil {
		return "", nil, err
	}

	return dockerHost, cli, nil
}
