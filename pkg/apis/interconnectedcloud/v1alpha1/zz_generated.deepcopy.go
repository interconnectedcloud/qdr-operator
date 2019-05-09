// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Address) DeepCopyInto(out *Address) {
	*out = *in
	if in.IngressPhase != nil {
		in, out := &in.IngressPhase, &out.IngressPhase
		*out = new(int32)
		**out = **in
	}
	if in.EgressPhase != nil {
		in, out := &in.EgressPhase, &out.EgressPhase
		*out = new(int32)
		**out = **in
	}
	if in.Priority != nil {
		in, out := &in.Priority, &out.Priority
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Address.
func (in *Address) DeepCopy() *Address {
	if in == nil {
		return nil
	}
	out := new(Address)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoLink) DeepCopyInto(out *AutoLink) {
	*out = *in
	if in.Phase != nil {
		in, out := &in.Phase, &out.Phase
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoLink.
func (in *AutoLink) DeepCopy() *AutoLink {
	if in == nil {
		return nil
	}
	out := new(AutoLink)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Connector) DeepCopyInto(out *Connector) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Connector.
func (in *Connector) DeepCopy() *Connector {
	if in == nil {
		return nil
	}
	out := new(Connector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploymentPlanType) DeepCopyInto(out *DeploymentPlanType) {
	*out = *in
	in.Resources.DeepCopyInto(&out.Resources)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploymentPlanType.
func (in *DeploymentPlanType) DeepCopy() *DeploymentPlanType {
	if in == nil {
		return nil
	}
	out := new(DeploymentPlanType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LinkRoute) DeepCopyInto(out *LinkRoute) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LinkRoute.
func (in *LinkRoute) DeepCopy() *LinkRoute {
	if in == nil {
		return nil
	}
	out := new(LinkRoute)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Listener) DeepCopyInto(out *Listener) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Listener.
func (in *Listener) DeepCopy() *Listener {
	if in == nil {
		return nil
	}
	out := new(Listener)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Qdr) DeepCopyInto(out *Qdr) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Qdr.
func (in *Qdr) DeepCopy() *Qdr {
	if in == nil {
		return nil
	}
	out := new(Qdr)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Qdr) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QdrCondition) DeepCopyInto(out *QdrCondition) {
	*out = *in
	in.TransitionTime.DeepCopyInto(&out.TransitionTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QdrCondition.
func (in *QdrCondition) DeepCopy() *QdrCondition {
	if in == nil {
		return nil
	}
	out := new(QdrCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QdrList) DeepCopyInto(out *QdrList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Qdr, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QdrList.
func (in *QdrList) DeepCopy() *QdrList {
	if in == nil {
		return nil
	}
	out := new(QdrList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *QdrList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QdrSpec) DeepCopyInto(out *QdrSpec) {
	*out = *in
	in.DeploymentPlan.DeepCopyInto(&out.DeploymentPlan)
	if in.Listeners != nil {
		in, out := &in.Listeners, &out.Listeners
		*out = make([]Listener, len(*in))
		copy(*out, *in)
	}
	if in.InterRouterListeners != nil {
		in, out := &in.InterRouterListeners, &out.InterRouterListeners
		*out = make([]Listener, len(*in))
		copy(*out, *in)
	}
	if in.EdgeListeners != nil {
		in, out := &in.EdgeListeners, &out.EdgeListeners
		*out = make([]Listener, len(*in))
		copy(*out, *in)
	}
	if in.SslProfiles != nil {
		in, out := &in.SslProfiles, &out.SslProfiles
		*out = make([]SslProfile, len(*in))
		copy(*out, *in)
	}
	if in.Addresses != nil {
		in, out := &in.Addresses, &out.Addresses
		*out = make([]Address, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AutoLinks != nil {
		in, out := &in.AutoLinks, &out.AutoLinks
		*out = make([]AutoLink, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.LinkRoutes != nil {
		in, out := &in.LinkRoutes, &out.LinkRoutes
		*out = make([]LinkRoute, len(*in))
		copy(*out, *in)
	}
	if in.Connectors != nil {
		in, out := &in.Connectors, &out.Connectors
		*out = make([]Connector, len(*in))
		copy(*out, *in)
	}
	if in.InterRouterConnectors != nil {
		in, out := &in.InterRouterConnectors, &out.InterRouterConnectors
		*out = make([]Connector, len(*in))
		copy(*out, *in)
	}
	if in.EdgeConnectors != nil {
		in, out := &in.EdgeConnectors, &out.EdgeConnectors
		*out = make([]Connector, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QdrSpec.
func (in *QdrSpec) DeepCopy() *QdrSpec {
	if in == nil {
		return nil
	}
	out := new(QdrSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QdrStatus) DeepCopyInto(out *QdrStatus) {
	*out = *in
	if in.PodNames != nil {
		in, out := &in.PodNames, &out.PodNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]QdrCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QdrStatus.
func (in *QdrStatus) DeepCopy() *QdrStatus {
	if in == nil {
		return nil
	}
	out := new(QdrStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SslProfile) DeepCopyInto(out *SslProfile) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SslProfile.
func (in *SslProfile) DeepCopy() *SslProfile {
	if in == nil {
		return nil
	}
	out := new(SslProfile)
	in.DeepCopyInto(out)
	return out
}
