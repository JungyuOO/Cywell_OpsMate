package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (in *OpsMateConfig) DeepCopyInto(out *OpsMateConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec.RAG != nil {
		out.Spec.RAG = make([]RAGSpec, len(in.Spec.RAG))
		copy(out.Spec.RAG, in.Spec.RAG)
	}
	if in.Status.Conditions != nil {
		out.Status.Conditions = make([]metav1.Condition, len(in.Status.Conditions))
		copy(out.Status.Conditions, in.Status.Conditions)
	}
}

func (in *OpsMateConfig) DeepCopy() *OpsMateConfig {
	if in == nil {
		return nil
	}
	out := new(OpsMateConfig)
	in.DeepCopyInto(out)
	return out
}

func (in *OpsMateConfig) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}

func (in *OpsMateConfigList) DeepCopyInto(out *OpsMateConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		out.Items = make([]OpsMateConfig, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

func (in *OpsMateConfigList) DeepCopy() *OpsMateConfigList {
	if in == nil {
		return nil
	}
	out := new(OpsMateConfigList)
	in.DeepCopyInto(out)
	return out
}

func (in *OpsMateConfigList) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}
