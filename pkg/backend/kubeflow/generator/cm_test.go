// Copyright 2018 Caicloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"testing"

	pytorchv1alpha2 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1alpha2"
	tfv1alpha2 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1alpha2"

	"github.com/caicloud/ciao/pkg/types"
)

func TestCMNewTFJob(t *testing.T) {
	cm := NewCM()

	expectedPSCount := 1
	expectedWorkerCount := 1
	expectedCM := "image"

	param := &types.Parameter{
		PSCount:     expectedPSCount,
		WorkerCount: expectedWorkerCount,
		Image:       expectedCM,
	}

	tfJob := cm.GenerateTFJob(param)
	actualPSCount := *tfJob.Spec.TFReplicaSpecs[tfv1alpha2.TFReplicaTypePS].Replicas
	actualWorkerCount := *tfJob.Spec.TFReplicaSpecs[tfv1alpha2.TFReplicaTypeWorker].Replicas
	actualCM := tfJob.Spec.TFReplicaSpecs[tfv1alpha2.TFReplicaTypePS].Template.Spec.Containers[0].VolumeMounts[0].Name
	if actualPSCount != int32(expectedPSCount) {
		t.Errorf("Expected %d ps, got %d", expectedPSCount, actualPSCount)
	}
	if actualWorkerCount != int32(expectedWorkerCount) {
		t.Errorf("Expected %d workers, got %d", expectedWorkerCount, actualWorkerCount)
	}
	if actualCM != expectedCM {
		t.Errorf("Expected configmap name %s, got %s", expectedCM, actualCM)
	}
}

func TestCMNewPyTorchJob(t *testing.T) {
	cm := NewCM()

	expectedMasterCount := 1
	expectedWorkerCount := 1
	expectedCM := "image"

	param := &types.Parameter{
		MasterCount: expectedMasterCount,
		WorkerCount: expectedWorkerCount,
		Image:       expectedCM,
	}

	pytorchJob := cm.GeneratePyTorchJob(param)
	actualMasterCount := *pytorchJob.Spec.PyTorchReplicaSpecs[pytorchv1alpha2.PyTorchReplicaTypeMaster].Replicas
	actualWorkerCount := *pytorchJob.Spec.PyTorchReplicaSpecs[pytorchv1alpha2.PyTorchReplicaTypeWorker].Replicas
	actualCM := pytorchJob.Spec.PyTorchReplicaSpecs[pytorchv1alpha2.PyTorchReplicaTypeMaster].Template.Spec.Containers[0].VolumeMounts[0].Name
	if actualMasterCount != int32(expectedMasterCount) {
		t.Errorf("Expected %d masters, got %d", expectedMasterCount, actualMasterCount)
	}
	if actualWorkerCount != int32(expectedWorkerCount) {
		t.Errorf("Expected %d workers, got %d", expectedWorkerCount, actualWorkerCount)
	}
	if actualCM != expectedCM {
		t.Errorf("Expected configmap name %s, got %s", expectedCM, actualCM)
	}
}
