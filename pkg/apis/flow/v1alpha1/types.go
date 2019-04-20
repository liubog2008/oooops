package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FlowList defines list of flow
type FlowList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items defines an array of flow
	Items []Flow `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Flow defines a CI/CD flow
type Flow struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of CI/CD flow
	// +optional
	Spec FlowSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// Status defines the current status of CI/CD flow
	// +optional
	Status FlowStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=spec"`
}

// FlowSpec defines spec of CI/CD flow
type FlowSpec struct {
	// Selector defines label selector for pod
	Selector *metav1.LabelSelector

	// Sources defines code source
	Sources []CodeSource `json:"sources"`

	// Images defines container images' config
	Images []Image `json:"images"`

	// Destination defines deploy config
	Destination Destination `json:"destination"`

	// Stage defines flow flow stage
	Stages []Stage `json:"stages"`

	// VolumeClaimTemplate defines whether
	VolumeClaimTemplate v1.PersistentVolumeClaim `json:"volumeClaimTemplate"`
}

// FlowStatus defines status of CI/CD flow
type FlowStatus struct {
	// Phase used to mark a simple phase of flow
	Phase string `json:"phase"`

	// Stage defines current stage of flow
	Stage string `json:"stage"`

	// Conditions defines flow conditions
	Conditions []FlowCondition `json:"conditions"`
}

// CodeSource defines source of code, e.g. github
type CodeSource struct {
	// Name defines name of the code source
	Name string `json:"name"`
	// Git defines git source
	Git GitSource `json:"git,omitempty"`
}

// Image defines container image config
type Image struct {
	// Name defines name of the image
	// Name can be referred by deploy template
	Name string `json:"name"`
	// Repository defines image repository
	Repository string `json:"repository"`
	// Path defines the path of image config
	// e.g. Dockerfile
	// TODO(liubog2008): support relative path of code source
	Path string `json:"path"`

	// WorkDir defines work dir to build image
	WorkDir string `json:"workDir"`
}

// When defines when to run the stage
type When string

const (
	// WhenCodeReady defines time point after code is ready
	WhenCodeReady When = "CodeReady"
	// WhenImageReady defines time point after iamge is ready
	WhenImageReady When = "ImageReady"
	// WhenApplicationReady defines time point after application is ready
	WhenApplicationReady When = "ApplicationReady"
)

// Stage defines CI/CD stage config
type Stage struct {
	// Name defines name of enviroment
	Name string `json:"name"`
	// When defines when to run the stage
	When When `json:"when"`
	// TemplateName defines pod template of enviroment
	TemplateName string `json:"template"`
	// Commands defines commands
	Commands []string `json:"commands"`
}

// Destination defines config
// TODO(liubog2008): add different config
type Destination struct {
	// Path defines dir which contains deploy yaml templates
	Path string `json:"path"`
}

// GitSource is a config of git
type GitSource struct {
	// Repository defines github repo
	Repository string `json:"repository"`

	// Type defines git source from
	// e.g. branch, pr, release, revision
	Type GitSourceType `json:"type"`

	// Matches defines matcher of git source
	// Now only support equals
	// If matcher is empty, branch will be master,
	// pr and release will be the lastest
	// TODO(liubog2008): support regexp matcher
	Matches string `json:"matches"`
}

// GitSourceType defines type of git source
type GitSourceType string

const (
	// GitBranch defines git source from specified branch
	// Default branch is master
	GitBranch GitSourceType = "Branch"
	// GitPullRequest defines git source from a pull request
	// Default is the latest pull request
	GitPullRequest GitSourceType = "PullRequest"
	// GitRelease defines git source from a release tag
	// TODO(liubog2008): default should use latest release
	GitRelease GitSourceType = "Release"
	// GitRevision defines git source from a specified revision
	// Default is the latest commit of master branch
	GitRevision GitSourceType = "Revision"
)

// FlowConditionType defines condition type of flow
type FlowConditionType string

// These are valid conditions of a flow
const (
	// FlowComplete means the flow has completed its execution.
	FlowComplete FlowConditionType = "Complete"
	// FlowFailed means the flow has failed its execution.
	FlowFailed FlowConditionType = "Failed"
	// FlowWaiting means the flow is waiting for next trigger
	FlowWaiting FlowConditionType = "Waiting"
)

// FlowCondition describes current state of a flow.
type FlowCondition struct {
	// Type of flow condition, Complete or Failed.
	Type FlowConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=FlowConditionType"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`
	// Last time the condition was checked.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty" protobuf:"bytes,3,opt,name=lastProbeTime"`
	// Last time the condition transit from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// (brief) reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// Human readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

const (
	// LabelStage used to defines stage of flow
	LabelStage = "alpha.oooops.com/stage"

	// LabelFlowHash used to defines flow definitions
	LabelFlowHash = "alpha.oooops.com/hash"
)

const (
	// StageSCM defines internal stage to fetch codes
	StageSCM = "scm"

	// StageImage defines internal stage to build and push image
	StageImage = "image"

	// StageDeploy defines internal stage to deploy application
	StageDeploy = "deploy"
)
