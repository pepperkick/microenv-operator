package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	awsec2 "pepperkick.com/microenv-operator/api/crossplane/aws-ec2"
)

// InstanceSpec defines the configuration of instances of the Cluster
type InstanceSpec struct {
	// Name denotes the name to suffix the EC2 instance name with
	Name string `json:"name"`
	// Type denotes the machine type of the EC2 instance
	Type string `json:"type"`
	// RootVolume defines the configuration of the root volume for the EC2 instance
	RootVolume awsec2.RootBlockDeviceParameters `json:"rootVolume,omitempty"`

	// Nodes denotes configuration for Kubernetes nodes created for this instance
	Nodes []InstanceNodesSpec `json:"nodes,omitempty"`

	// DockerSwarmRole denotes the role that EC2 instance will use for docker swarm.
	// This key is ignored. First instance will be manager and rest will be workers.
	DockerSwarmRole string `json:"dockerSwarmRole,omitempty"`
}

type InstanceNodesSpec struct {
	// Labels denotes the labels to apply to Kubernetes node
	Labels map[string]string `json:"labels,omitempty"`

	// Taints denotes the labels to apply to Kubernetes node
	Taints []corev1.Taint `json:"taints,omitempty"`
}

// InfrastructureSpec defines the desired infrastructure state of Cluster
type InfrastructureSpec struct {
	Instances []InstanceSpec `json:"instances"`

	// CertIssuer points to which certificate issuer to use
	// Following values are supported
	// - cert-manager
	// - letsencrypt
	CertIssuer string `json:"certIssuer"`

	// AllowInstanceUpdateForUserData will update the instance CR when there is a user data change
	// This means the instances might restart if needed
	AllowInstanceUpdateForUserData bool `json:"allowInstanceUpdateForUserData,omitempty"`

	Proxy ProxyConfigSpec `json:"proxy,omitempty"`
}

// FeaturesSpec defines the desired usage state of Cluster
type FeaturesSpec struct {
	// === INGRESS NGINX OPTIONS === \\

	// InstallIngressNginx defines if ingress nginx controller should be installed or not
	InstallIngressNginx bool `json:"InstallIngressNginx,omitempty"`

	// === ARGO WORKFLOW OPTIONS === \\

	// InstallArgoWorkflow defines if argo workflow should be installed or not
	InstallArgoWorkflow bool `json:"installArgoWorkflow,omitempty"`

	// ArgoWorkflow defines the workflow to run, ignored if InstallArgoWorkflow is set to false
	// The contents has to be gzipped and base64 encoded to ensure it is read properly
	ArgoWorkflow string `json:"argoWorkflow,omitempty"`
}

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	// Provider points to name of existing Provider resource
	Provider       string             `json:"provider,omitempty"`
	Email          string             `json:"email"`
	Infrastructure InfrastructureSpec `json:"infrastructure"`
	Features       FeaturesSpec       `json:"features,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// Conditions contains the current status of the cluster
	Conditions []metav1.Condition `json:"conditions"`

	// ManagerInstanceIp points to private IP of cluster's manager instance
	ManagerInstanceIp string `json:"managerInstanceIp,omitempty"`
	// ManagerInstancePublicIp points to public IP of cluster's manager instance
	ManagerInstancePublicIp string `json:"managerInstancePublicIp,omitempty"`
	// ClusterIngressDomain points to ingress domain for the cluster
	ClusterIngressDomain string `json:"clusterIngressDomain,omitempty"`
	// CertificateSecret points to secret which contains the domain certificate
	CertificateSecret string `json:"certificateSecret,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster,categories={microenv}

// Cluster is the Schema for the clusters API
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type=='Reconciled')].message"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Domain",type="string",JSONPath=".status.clusterIngressDomain"
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
