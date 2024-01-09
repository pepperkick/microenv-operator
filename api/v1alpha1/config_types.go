package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// AwsConfigSpec defines the AWS configuration
type AwsConfigSpec struct {
	// Region is to use for creating resources
	Region string `json:"region"`
	// BaseAmiId points to the base AMI to use for instances
	BaseAmiId string `json:"baseAmiId"`
	// SubnetId points to subnet to use for instances
	SubnetId string `json:"subnetId"`
	// SecurityGroupId points to security group to use for instances
	SecurityGroupId string `json:"securityGroupId"`
	// IamInstanceProfileName points to IAM role to use for instances
	IamInstanceProfileName string `json:"iamInstanceProfileName"`
	// InstanceType points to default instance type to use
	InstanceType string `json:"instanceType,omitempty"`

	// PcaCertIssuer points to the name of AWS cert issuer
	PcaCertIssuer string `json:"pcaCertIssuer,omitempty"`

	// Route53HostedZone points to Route53 hosted zone for creating certs and dns entries
	Route53HostedZone string `json:"route53HostedZone,omitempty"`
	// IngressDomain points to the domain to use, env domains will subdomains of it
	IngressDomain string `json:"ingressDomain,omitempty"`

	// UsePrivateIp points to using private ip or public ip
	UsePrivateIp bool `json:"usePrivateIp,omitempty"`
}

// ProxyConfigSpec defines the configuration for http and https proxy for all the instance
type ProxyConfigSpec struct {
	Enabled       bool     `json:"enabled"`
	HttpEndpoint  string   `json:"httpEndpoint"`
	HttpsEndpoint string   `json:"httpsEndpoint"`
	Exclusions    []string `json:"exclusions"`
}

type UtilImageSpec struct {
	Image            string `json:"image"`
	RegistryUsername string `json:"registryUsername"`
	RegistryPassword string `json:"registryPassword"`
}

// ConfigSpec defines the provider configuration
type ConfigSpec struct {
	Proxy        ProxyConfigSpec `json:"proxy,omitempty"`
	ProviderName string          `json:"providerName,omitempty"`

	// SystemNamespace points to namespace to store secrets and resources
	SystemNamespace string `json:"systemNamespace"`

	// InstanceSetupScript defines the script to run during instance setup
	// This will run within the normal setup script
	InstanceSetupScript string `json:"instanceSetupScript,omitempty"`

	UtilImage UtilImageSpec `json:"utilImage,omitempty"`

	// KubernetesVersion denotes the kubernetes version to deploy
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	Aws AwsConfigSpec `json:"aws"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster,categories={microenv}

// Config is the Schema for the config provider API
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ConfigSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// ConfigList contains a list of Config
type ConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Config `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Config{}, &ConfigList{})
}
