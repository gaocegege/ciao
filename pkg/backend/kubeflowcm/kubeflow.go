package kubeflowcm

import (
	"fmt"

	kubeflowbackend "github.com/caicloud/ciao/pkg/backend/kubeflow"
	s2iconfigmap "github.com/caicloud/ciao/pkg/s2i/configmap"
	pyttorchjobclient "github.com/kubeflow/pytorch-operator/pkg/client/clientset/versioned"
	tfv1alpha2 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1alpha2"
	tfjobclient "github.com/kubeflow/tf-operator/pkg/client/clientset/versioned"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	restclientset "k8s.io/client-go/rest"

	"github.com/caicloud/ciao/pkg/types"
)

const (
	baseImageTF      = ""
	baseImagePyTorch = ""
)

type Backend struct {
	kubeflowbackend.Backend
}

// New returns a new Backend.
func New(config *restclientset.Config) (*Backend, error) {
	tfJobClient, err := tfjobclient.NewForConfig(restclientset.AddUserAgent(config, kubeflowbackend.UserAgent))
	if err != nil {
		return nil, err
	}
	k8sClient, err := kubeclient.NewForConfig(restclientset.AddUserAgent(config, kubeflowbackend.UserAgent))
	if err != nil {
		return nil, err
	}
	pytorchClient, err := pyttorchjobclient.NewForConfig(restclientset.AddUserAgent(config, kubeflowbackend.UserAgent))
	if err != nil {
		return nil, err
	}

	return &Backend{
		Backend: kubeflowbackend.Backend{
			TFJobClient:      tfJobClient,
			K8sClient:        k8sClient,
			PyTorchJobClient: pytorchClient,
		},
	}, nil
}

func (b Backend) generateTFJob(parameter *types.Parameter) *tfv1alpha2.TFJob {
	psCount := int32(parameter.PSCount)
	workerCount := int32(parameter.WorkerCount)

	mountPath := fmt.Sprintf("/%s", parameter.Image)
	filename := fmt.Sprintf("/%s/%s", parameter.Image, s2iconfigmap.FileName)

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
									Name:  kubeflowbackend.DefaultContainerNameTF,
									Image: baseImageTF,
									Command: []string{
										"python",
										filename,
									},
									VolumeMounts: []v1.VolumeMount{
										v1.VolumeMount{
											Name:      parameter.Image,
											MountPath: mountPath,
										},
									},
								},
							},
							Volumes: []v1.Volume{
								v1.Volume{
									Name: parameter.Image,
									VolumeSource: v1.VolumeSource{
										ConfigMap: &v1.ConfigMapVolumeSource{
											LocalObjectReference: v1.LocalObjectReference{
												Name: parameter.Image,
											},
										},
									},
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
									Name:  kubeflowbackend.DefaultContainerNameTF,
									Image: baseImageTF,
									Command: []string{
										"python",
										filename,
									},
									VolumeMounts: []v1.VolumeMount{
										v1.VolumeMount{
											Name:      parameter.Image,
											MountPath: mountPath,
										},
									},
								},
							},
							Volumes: []v1.Volume{
								v1.Volume{
									Name: parameter.Image,
									VolumeSource: v1.VolumeSource{
										ConfigMap: &v1.ConfigMapVolumeSource{
											LocalObjectReference: v1.LocalObjectReference{
												Name: parameter.Image,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
