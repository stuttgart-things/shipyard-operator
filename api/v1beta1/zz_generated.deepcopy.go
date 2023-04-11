//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023 patrick hermann.

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

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ShipyardTerraform) DeepCopyInto(out *ShipyardTerraform) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ShipyardTerraform.
func (in *ShipyardTerraform) DeepCopy() *ShipyardTerraform {
	if in == nil {
		return nil
	}
	out := new(ShipyardTerraform)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ShipyardTerraform) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ShipyardTerraformList) DeepCopyInto(out *ShipyardTerraformList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ShipyardTerraform, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ShipyardTerraformList.
func (in *ShipyardTerraformList) DeepCopy() *ShipyardTerraformList {
	if in == nil {
		return nil
	}
	out := new(ShipyardTerraformList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ShipyardTerraformList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ShipyardTerraformSpec) DeepCopyInto(out *ShipyardTerraformSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ShipyardTerraformSpec.
func (in *ShipyardTerraformSpec) DeepCopy() *ShipyardTerraformSpec {
	if in == nil {
		return nil
	}
	out := new(ShipyardTerraformSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ShipyardTerraformStatus) DeepCopyInto(out *ShipyardTerraformStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ShipyardTerraformStatus.
func (in *ShipyardTerraformStatus) DeepCopy() *ShipyardTerraformStatus {
	if in == nil {
		return nil
	}
	out := new(ShipyardTerraformStatus)
	in.DeepCopyInto(out)
	return out
}
