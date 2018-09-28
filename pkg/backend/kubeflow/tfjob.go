package kubeflow

import (
	tfv1alpha2 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1alpha2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caicloud/ciao/pkg/types"
)

const (
	DefaultContainerNameTF = "tensorflow"
)

func (b Backend) createTFJob(parameter *types.Parameter) (*types.Job, error) {
	tfJob := b.generateTFJob(parameter)
	tfJob, err := b.TFJobClient.KubeflowV1alpha2().TFJobs(metav1.NamespaceDefault).Create(tfJob)
	if err != nil {
		return nil, err
	}
	return &types.Job{
		Name:      tfJob.Name,
		Framework: types.FrameworkTypeTensorFlow,
		PS:        parameter.PSCount,
		Worker:    parameter.WorkerCount,
	}, nil
}

func (b Backend) generateTFJob(parameter *types.Parameter) *tfv1alpha2.TFJob {
	psCount := int32(parameter.PSCount)
	workerCount := int32(parameter.WorkerCount)

	return &tfv1alpha2.TFJob{
		TypeMeta: metav1.TypeMeta{
			Kind: tfv1alpha2.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      parameter.GenerateName,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: tfv1alpha2.TFJobSpec{
			TFReplicaSpecs: map[tfv1alpha2.TFReplicaType]*tfv1alpha2.TFReplicaSpec{
				tfv1alpha2.TFReplicaTypePS: &tfv1alpha2.TFReplicaSpec{
					Replicas: &psCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Name:  DefaultContainerNameTF,
									Image: parameter.Image,
								},
							},
						},
					},
				},
				tfv1alpha2.TFReplicaTypeWorker: &tfv1alpha2.TFReplicaSpec{
					Replicas: &workerCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Name:  DefaultContainerNameTF,
									Image: parameter.Image,
								},
							},
						},
					},
				},
			},
		},
	}
}
