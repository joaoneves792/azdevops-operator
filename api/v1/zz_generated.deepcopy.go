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

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzDevopsAgentPool) DeepCopyInto(out *AzDevopsAgentPool) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzDevopsAgentPool.
func (in *AzDevopsAgentPool) DeepCopy() *AzDevopsAgentPool {
	if in == nil {
		return nil
	}
	out := new(AzDevopsAgentPool)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AzDevopsAgentPool) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzDevopsAgentPoolAutoscaling) DeepCopyInto(out *AzDevopsAgentPoolAutoscaling) {
	*out = *in
	out.Schedule = in.Schedule
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzDevopsAgentPoolAutoscaling.
func (in *AzDevopsAgentPoolAutoscaling) DeepCopy() *AzDevopsAgentPoolAutoscaling {
	if in == nil {
		return nil
	}
	out := new(AzDevopsAgentPoolAutoscaling)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzDevopsAgentPoolAutoscalingSchedule) DeepCopyInto(out *AzDevopsAgentPoolAutoscalingSchedule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzDevopsAgentPoolAutoscalingSchedule.
func (in *AzDevopsAgentPoolAutoscalingSchedule) DeepCopy() *AzDevopsAgentPoolAutoscalingSchedule {
	if in == nil {
		return nil
	}
	out := new(AzDevopsAgentPoolAutoscalingSchedule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzDevopsAgentPoolList) DeepCopyInto(out *AzDevopsAgentPoolList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AzDevopsAgentPool, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzDevopsAgentPoolList.
func (in *AzDevopsAgentPoolList) DeepCopy() *AzDevopsAgentPoolList {
	if in == nil {
		return nil
	}
	out := new(AzDevopsAgentPoolList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AzDevopsAgentPoolList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzDevopsAgentPoolSpec) DeepCopyInto(out *AzDevopsAgentPoolSpec) {
	*out = *in
	out.Project = in.Project
	out.Autoscaling = in.Autoscaling
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzDevopsAgentPoolSpec.
func (in *AzDevopsAgentPoolSpec) DeepCopy() *AzDevopsAgentPoolSpec {
	if in == nil {
		return nil
	}
	out := new(AzDevopsAgentPoolSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzDevopsAgentPoolStatus) DeepCopyInto(out *AzDevopsAgentPoolStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzDevopsAgentPoolStatus.
func (in *AzDevopsAgentPoolStatus) DeepCopy() *AzDevopsAgentPoolStatus {
	if in == nil {
		return nil
	}
	out := new(AzDevopsAgentPoolStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzDevopsProject) DeepCopyInto(out *AzDevopsProject) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzDevopsProject.
func (in *AzDevopsProject) DeepCopy() *AzDevopsProject {
	if in == nil {
		return nil
	}
	out := new(AzDevopsProject)
	in.DeepCopyInto(out)
	return out
}
