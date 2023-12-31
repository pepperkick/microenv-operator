//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package aws_route53

import (
	"github.com/crossplane/crossplane-runtime/apis/common/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AliasInitParameters) DeepCopyInto(out *AliasInitParameters) {
	*out = *in
	if in.EvaluateTargetHealth != nil {
		in, out := &in.EvaluateTargetHealth, &out.EvaluateTargetHealth
		*out = new(bool)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.ZoneID != nil {
		in, out := &in.ZoneID, &out.ZoneID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AliasInitParameters.
func (in *AliasInitParameters) DeepCopy() *AliasInitParameters {
	if in == nil {
		return nil
	}
	out := new(AliasInitParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AliasObservation) DeepCopyInto(out *AliasObservation) {
	*out = *in
	if in.EvaluateTargetHealth != nil {
		in, out := &in.EvaluateTargetHealth, &out.EvaluateTargetHealth
		*out = new(bool)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.ZoneID != nil {
		in, out := &in.ZoneID, &out.ZoneID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AliasObservation.
func (in *AliasObservation) DeepCopy() *AliasObservation {
	if in == nil {
		return nil
	}
	out := new(AliasObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AliasParameters) DeepCopyInto(out *AliasParameters) {
	*out = *in
	if in.EvaluateTargetHealth != nil {
		in, out := &in.EvaluateTargetHealth, &out.EvaluateTargetHealth
		*out = new(bool)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.ZoneID != nil {
		in, out := &in.ZoneID, &out.ZoneID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AliasParameters.
func (in *AliasParameters) DeepCopy() *AliasParameters {
	if in == nil {
		return nil
	}
	out := new(AliasParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CidrRoutingPolicyInitParameters) DeepCopyInto(out *CidrRoutingPolicyInitParameters) {
	*out = *in
	if in.CollectionID != nil {
		in, out := &in.CollectionID, &out.CollectionID
		*out = new(string)
		**out = **in
	}
	if in.LocationName != nil {
		in, out := &in.LocationName, &out.LocationName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CidrRoutingPolicyInitParameters.
func (in *CidrRoutingPolicyInitParameters) DeepCopy() *CidrRoutingPolicyInitParameters {
	if in == nil {
		return nil
	}
	out := new(CidrRoutingPolicyInitParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CidrRoutingPolicyObservation) DeepCopyInto(out *CidrRoutingPolicyObservation) {
	*out = *in
	if in.CollectionID != nil {
		in, out := &in.CollectionID, &out.CollectionID
		*out = new(string)
		**out = **in
	}
	if in.LocationName != nil {
		in, out := &in.LocationName, &out.LocationName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CidrRoutingPolicyObservation.
func (in *CidrRoutingPolicyObservation) DeepCopy() *CidrRoutingPolicyObservation {
	if in == nil {
		return nil
	}
	out := new(CidrRoutingPolicyObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CidrRoutingPolicyParameters) DeepCopyInto(out *CidrRoutingPolicyParameters) {
	*out = *in
	if in.CollectionID != nil {
		in, out := &in.CollectionID, &out.CollectionID
		*out = new(string)
		**out = **in
	}
	if in.LocationName != nil {
		in, out := &in.LocationName, &out.LocationName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CidrRoutingPolicyParameters.
func (in *CidrRoutingPolicyParameters) DeepCopy() *CidrRoutingPolicyParameters {
	if in == nil {
		return nil
	}
	out := new(CidrRoutingPolicyParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FailoverRoutingPolicyInitParameters) DeepCopyInto(out *FailoverRoutingPolicyInitParameters) {
	*out = *in
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FailoverRoutingPolicyInitParameters.
func (in *FailoverRoutingPolicyInitParameters) DeepCopy() *FailoverRoutingPolicyInitParameters {
	if in == nil {
		return nil
	}
	out := new(FailoverRoutingPolicyInitParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FailoverRoutingPolicyObservation) DeepCopyInto(out *FailoverRoutingPolicyObservation) {
	*out = *in
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FailoverRoutingPolicyObservation.
func (in *FailoverRoutingPolicyObservation) DeepCopy() *FailoverRoutingPolicyObservation {
	if in == nil {
		return nil
	}
	out := new(FailoverRoutingPolicyObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FailoverRoutingPolicyParameters) DeepCopyInto(out *FailoverRoutingPolicyParameters) {
	*out = *in
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FailoverRoutingPolicyParameters.
func (in *FailoverRoutingPolicyParameters) DeepCopy() *FailoverRoutingPolicyParameters {
	if in == nil {
		return nil
	}
	out := new(FailoverRoutingPolicyParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeolocationRoutingPolicyInitParameters) DeepCopyInto(out *GeolocationRoutingPolicyInitParameters) {
	*out = *in
	if in.Continent != nil {
		in, out := &in.Continent, &out.Continent
		*out = new(string)
		**out = **in
	}
	if in.Country != nil {
		in, out := &in.Country, &out.Country
		*out = new(string)
		**out = **in
	}
	if in.Subdivision != nil {
		in, out := &in.Subdivision, &out.Subdivision
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeolocationRoutingPolicyInitParameters.
func (in *GeolocationRoutingPolicyInitParameters) DeepCopy() *GeolocationRoutingPolicyInitParameters {
	if in == nil {
		return nil
	}
	out := new(GeolocationRoutingPolicyInitParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeolocationRoutingPolicyObservation) DeepCopyInto(out *GeolocationRoutingPolicyObservation) {
	*out = *in
	if in.Continent != nil {
		in, out := &in.Continent, &out.Continent
		*out = new(string)
		**out = **in
	}
	if in.Country != nil {
		in, out := &in.Country, &out.Country
		*out = new(string)
		**out = **in
	}
	if in.Subdivision != nil {
		in, out := &in.Subdivision, &out.Subdivision
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeolocationRoutingPolicyObservation.
func (in *GeolocationRoutingPolicyObservation) DeepCopy() *GeolocationRoutingPolicyObservation {
	if in == nil {
		return nil
	}
	out := new(GeolocationRoutingPolicyObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeolocationRoutingPolicyParameters) DeepCopyInto(out *GeolocationRoutingPolicyParameters) {
	*out = *in
	if in.Continent != nil {
		in, out := &in.Continent, &out.Continent
		*out = new(string)
		**out = **in
	}
	if in.Country != nil {
		in, out := &in.Country, &out.Country
		*out = new(string)
		**out = **in
	}
	if in.Subdivision != nil {
		in, out := &in.Subdivision, &out.Subdivision
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeolocationRoutingPolicyParameters.
func (in *GeolocationRoutingPolicyParameters) DeepCopy() *GeolocationRoutingPolicyParameters {
	if in == nil {
		return nil
	}
	out := new(GeolocationRoutingPolicyParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LatencyRoutingPolicyInitParameters) DeepCopyInto(out *LatencyRoutingPolicyInitParameters) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LatencyRoutingPolicyInitParameters.
func (in *LatencyRoutingPolicyInitParameters) DeepCopy() *LatencyRoutingPolicyInitParameters {
	if in == nil {
		return nil
	}
	out := new(LatencyRoutingPolicyInitParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LatencyRoutingPolicyObservation) DeepCopyInto(out *LatencyRoutingPolicyObservation) {
	*out = *in
	if in.Region != nil {
		in, out := &in.Region, &out.Region
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LatencyRoutingPolicyObservation.
func (in *LatencyRoutingPolicyObservation) DeepCopy() *LatencyRoutingPolicyObservation {
	if in == nil {
		return nil
	}
	out := new(LatencyRoutingPolicyObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LatencyRoutingPolicyParameters) DeepCopyInto(out *LatencyRoutingPolicyParameters) {
	*out = *in
	if in.Region != nil {
		in, out := &in.Region, &out.Region
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LatencyRoutingPolicyParameters.
func (in *LatencyRoutingPolicyParameters) DeepCopy() *LatencyRoutingPolicyParameters {
	if in == nil {
		return nil
	}
	out := new(LatencyRoutingPolicyParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Record) DeepCopyInto(out *Record) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Record.
func (in *Record) DeepCopy() *Record {
	if in == nil {
		return nil
	}
	out := new(Record)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Record) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RecordInitParameters) DeepCopyInto(out *RecordInitParameters) {
	*out = *in
	if in.Alias != nil {
		in, out := &in.Alias, &out.Alias
		*out = make([]AliasInitParameters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AllowOverwrite != nil {
		in, out := &in.AllowOverwrite, &out.AllowOverwrite
		*out = new(bool)
		**out = **in
	}
	if in.CidrRoutingPolicy != nil {
		in, out := &in.CidrRoutingPolicy, &out.CidrRoutingPolicy
		*out = make([]CidrRoutingPolicyInitParameters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.FailoverRoutingPolicy != nil {
		in, out := &in.FailoverRoutingPolicy, &out.FailoverRoutingPolicy
		*out = make([]FailoverRoutingPolicyInitParameters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.GeolocationRoutingPolicy != nil {
		in, out := &in.GeolocationRoutingPolicy, &out.GeolocationRoutingPolicy
		*out = make([]GeolocationRoutingPolicyInitParameters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.LatencyRoutingPolicy != nil {
		in, out := &in.LatencyRoutingPolicy, &out.LatencyRoutingPolicy
		*out = make([]LatencyRoutingPolicyInitParameters, len(*in))
		copy(*out, *in)
	}
	if in.MultivalueAnswerRoutingPolicy != nil {
		in, out := &in.MultivalueAnswerRoutingPolicy, &out.MultivalueAnswerRoutingPolicy
		*out = new(bool)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Records != nil {
		in, out := &in.Records, &out.Records
		*out = make([]*string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(string)
				**out = **in
			}
		}
	}
	if in.SetIdentifier != nil {
		in, out := &in.SetIdentifier, &out.SetIdentifier
		*out = new(string)
		**out = **in
	}
	if in.TTL != nil {
		in, out := &in.TTL, &out.TTL
		*out = new(float64)
		**out = **in
	}
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
	if in.WeightedRoutingPolicy != nil {
		in, out := &in.WeightedRoutingPolicy, &out.WeightedRoutingPolicy
		*out = make([]WeightedRoutingPolicyInitParameters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RecordInitParameters.
func (in *RecordInitParameters) DeepCopy() *RecordInitParameters {
	if in == nil {
		return nil
	}
	out := new(RecordInitParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RecordList) DeepCopyInto(out *RecordList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Record, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RecordList.
func (in *RecordList) DeepCopy() *RecordList {
	if in == nil {
		return nil
	}
	out := new(RecordList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RecordList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RecordObservation) DeepCopyInto(out *RecordObservation) {
	*out = *in
	if in.Alias != nil {
		in, out := &in.Alias, &out.Alias
		*out = make([]AliasObservation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AllowOverwrite != nil {
		in, out := &in.AllowOverwrite, &out.AllowOverwrite
		*out = new(bool)
		**out = **in
	}
	if in.CidrRoutingPolicy != nil {
		in, out := &in.CidrRoutingPolicy, &out.CidrRoutingPolicy
		*out = make([]CidrRoutingPolicyObservation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.FailoverRoutingPolicy != nil {
		in, out := &in.FailoverRoutingPolicy, &out.FailoverRoutingPolicy
		*out = make([]FailoverRoutingPolicyObservation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Fqdn != nil {
		in, out := &in.Fqdn, &out.Fqdn
		*out = new(string)
		**out = **in
	}
	if in.GeolocationRoutingPolicy != nil {
		in, out := &in.GeolocationRoutingPolicy, &out.GeolocationRoutingPolicy
		*out = make([]GeolocationRoutingPolicyObservation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.HealthCheckID != nil {
		in, out := &in.HealthCheckID, &out.HealthCheckID
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.LatencyRoutingPolicy != nil {
		in, out := &in.LatencyRoutingPolicy, &out.LatencyRoutingPolicy
		*out = make([]LatencyRoutingPolicyObservation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.MultivalueAnswerRoutingPolicy != nil {
		in, out := &in.MultivalueAnswerRoutingPolicy, &out.MultivalueAnswerRoutingPolicy
		*out = new(bool)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Records != nil {
		in, out := &in.Records, &out.Records
		*out = make([]*string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(string)
				**out = **in
			}
		}
	}
	if in.SetIdentifier != nil {
		in, out := &in.SetIdentifier, &out.SetIdentifier
		*out = new(string)
		**out = **in
	}
	if in.TTL != nil {
		in, out := &in.TTL, &out.TTL
		*out = new(float64)
		**out = **in
	}
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
	if in.WeightedRoutingPolicy != nil {
		in, out := &in.WeightedRoutingPolicy, &out.WeightedRoutingPolicy
		*out = make([]WeightedRoutingPolicyObservation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ZoneID != nil {
		in, out := &in.ZoneID, &out.ZoneID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RecordObservation.
func (in *RecordObservation) DeepCopy() *RecordObservation {
	if in == nil {
		return nil
	}
	out := new(RecordObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RecordParameters) DeepCopyInto(out *RecordParameters) {
	*out = *in
	if in.Alias != nil {
		in, out := &in.Alias, &out.Alias
		*out = make([]AliasParameters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AllowOverwrite != nil {
		in, out := &in.AllowOverwrite, &out.AllowOverwrite
		*out = new(bool)
		**out = **in
	}
	if in.CidrRoutingPolicy != nil {
		in, out := &in.CidrRoutingPolicy, &out.CidrRoutingPolicy
		*out = make([]CidrRoutingPolicyParameters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.FailoverRoutingPolicy != nil {
		in, out := &in.FailoverRoutingPolicy, &out.FailoverRoutingPolicy
		*out = make([]FailoverRoutingPolicyParameters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.GeolocationRoutingPolicy != nil {
		in, out := &in.GeolocationRoutingPolicy, &out.GeolocationRoutingPolicy
		*out = make([]GeolocationRoutingPolicyParameters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.HealthCheckID != nil {
		in, out := &in.HealthCheckID, &out.HealthCheckID
		*out = new(string)
		**out = **in
	}
	if in.HealthCheckIDRef != nil {
		in, out := &in.HealthCheckIDRef, &out.HealthCheckIDRef
		*out = new(v1.Reference)
		(*in).DeepCopyInto(*out)
	}
	if in.HealthCheckIDSelector != nil {
		in, out := &in.HealthCheckIDSelector, &out.HealthCheckIDSelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
	if in.LatencyRoutingPolicy != nil {
		in, out := &in.LatencyRoutingPolicy, &out.LatencyRoutingPolicy
		*out = make([]LatencyRoutingPolicyParameters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.MultivalueAnswerRoutingPolicy != nil {
		in, out := &in.MultivalueAnswerRoutingPolicy, &out.MultivalueAnswerRoutingPolicy
		*out = new(bool)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Records != nil {
		in, out := &in.Records, &out.Records
		*out = make([]*string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(string)
				**out = **in
			}
		}
	}
	if in.Region != nil {
		in, out := &in.Region, &out.Region
		*out = new(string)
		**out = **in
	}
	if in.SetIdentifier != nil {
		in, out := &in.SetIdentifier, &out.SetIdentifier
		*out = new(string)
		**out = **in
	}
	if in.TTL != nil {
		in, out := &in.TTL, &out.TTL
		*out = new(float64)
		**out = **in
	}
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
	if in.WeightedRoutingPolicy != nil {
		in, out := &in.WeightedRoutingPolicy, &out.WeightedRoutingPolicy
		*out = make([]WeightedRoutingPolicyParameters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ZoneID != nil {
		in, out := &in.ZoneID, &out.ZoneID
		*out = new(string)
		**out = **in
	}
	if in.ZoneIDRef != nil {
		in, out := &in.ZoneIDRef, &out.ZoneIDRef
		*out = new(v1.Reference)
		(*in).DeepCopyInto(*out)
	}
	if in.ZoneIDSelector != nil {
		in, out := &in.ZoneIDSelector, &out.ZoneIDSelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RecordParameters.
func (in *RecordParameters) DeepCopy() *RecordParameters {
	if in == nil {
		return nil
	}
	out := new(RecordParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RecordSpec) DeepCopyInto(out *RecordSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
	in.InitProvider.DeepCopyInto(&out.InitProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RecordSpec.
func (in *RecordSpec) DeepCopy() *RecordSpec {
	if in == nil {
		return nil
	}
	out := new(RecordSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RecordStatus) DeepCopyInto(out *RecordStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	in.AtProvider.DeepCopyInto(&out.AtProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RecordStatus.
func (in *RecordStatus) DeepCopy() *RecordStatus {
	if in == nil {
		return nil
	}
	out := new(RecordStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WeightedRoutingPolicyInitParameters) DeepCopyInto(out *WeightedRoutingPolicyInitParameters) {
	*out = *in
	if in.Weight != nil {
		in, out := &in.Weight, &out.Weight
		*out = new(float64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WeightedRoutingPolicyInitParameters.
func (in *WeightedRoutingPolicyInitParameters) DeepCopy() *WeightedRoutingPolicyInitParameters {
	if in == nil {
		return nil
	}
	out := new(WeightedRoutingPolicyInitParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WeightedRoutingPolicyObservation) DeepCopyInto(out *WeightedRoutingPolicyObservation) {
	*out = *in
	if in.Weight != nil {
		in, out := &in.Weight, &out.Weight
		*out = new(float64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WeightedRoutingPolicyObservation.
func (in *WeightedRoutingPolicyObservation) DeepCopy() *WeightedRoutingPolicyObservation {
	if in == nil {
		return nil
	}
	out := new(WeightedRoutingPolicyObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WeightedRoutingPolicyParameters) DeepCopyInto(out *WeightedRoutingPolicyParameters) {
	*out = *in
	if in.Weight != nil {
		in, out := &in.Weight, &out.Weight
		*out = new(float64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WeightedRoutingPolicyParameters.
func (in *WeightedRoutingPolicyParameters) DeepCopy() *WeightedRoutingPolicyParameters {
	if in == nil {
		return nil
	}
	out := new(WeightedRoutingPolicyParameters)
	in.DeepCopyInto(out)
	return out
}
