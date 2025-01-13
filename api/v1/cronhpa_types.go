/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CronHPASpec defines the desired state of CronHPA.
type CronHPASpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of CronHPA. Edit cronhpa_types.go to remove/update
	// jobs 定义多个扩缩容任务
	SyncUserGroup   []SyncUserGroup   `json:"syncUserGroup"`
	SyncAssetsGroup []SyncAssetsGroup `json:"syncAssetsGroup"`
}

// JobSpec 定义扩缩容任务的规格
type SyncUserGroup struct {
	// Name 表示扩缩容任务的名称
	Name string `json:"name"`

	// Schedule 表示 Cron 表达式，定义任务的调度时间
	Schedule string `json:"schedule"`

	// HostSourceURL
	UserSourceURL string `json:"userSourceUrl"`

	// JumpserverURL
	JumpserverURL string `json:"jumpserverUrl"`

	// JumpserverAK
	JumpserverAK string `json:"jumpserverAccessKey"`

	// JumpserverSK
	JumpserverSK string `json:"jumpserverSecurityKey"`

	// UserGroupName
	UserGroupName string `json:"userGroupName"`
}

type SyncAssetsGroup struct {
	// Name 表示扩缩容任务的名称
	Name string `json:"name"`

	// Schedule 表示 Cron 表达式，定义任务的调度时间
	Schedule string `json:"schedule"`

	// HostSourceURL
	HostSourceURL string `json:"hostSourceUrl"`

	// JumpserverURL
	JumpserverURL string `json:"jumpserverUrl"`

	// JumpserverAK
	JumpserverAK string `json:"jumpserverAccessKey"`

	// JumpserverSK
	JumpserverSK string `json:"jumpserverSecurityKey"`

	// AssetsPoint
	AssetsPointId string `json:"assetsPointId"`
}

// CronHPAStatus defines the observed state of CronHPA.
type CronHPAStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// LastScaleTime 表示最后一次扩缩容的时间
	LastScaleTime *metav1.Time `json:"lastScaleTime,omitempty"`
	// LastRunTimes 记录每个作业的最后运行时间
	LastRunTimes map[string]metav1.Time `json:"lastRunTimes,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CronHPA is the Schema for the cronhpas API.
type CronHPA struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CronHPASpec   `json:"spec,omitempty"`
	Status CronHPAStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CronHPAList contains a list of CronHPA.
type CronHPAList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CronHPA `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CronHPA{}, &CronHPAList{})
}
