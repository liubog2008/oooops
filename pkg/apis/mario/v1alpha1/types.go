package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// MarioFile defines file of Mario API in git project
	MarioFile = ".mario"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MarioList defines list of pipe
type MarioList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items defines an array of mario
	Items []Mario `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipeList defines list of pipe
type PipeList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items defines an array of pipe
	Items []Pipe `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FlowList defines list of flow
type FlowList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items defines an array of flow
	Items []Flow `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventList defines list of event
type EventList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items defines an array of event
	Items []Event `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Mario is API for user to define project action such as compile and build
type Mario struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines desired props of Mario
	// +optional
	Spec MarioSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// MarioSpec defines spec of Mario
type MarioSpec struct {
	// Actions defines actions of the project
	// e.g. compile, test
	// +optional
	Actions []Action `json:"actions,omitempty" protobuf:"bytes,1,rep,name=actions"`
	// VersionEnv defines env name whose value will be the version
	// +optional
	VersionEnv string `json:"versionEnv,omitempty" protobuf:"bytes,2,opt,name=versionEnv"`
}

// Action defines custom action defined by users
type Action struct {
	// Name defines action name
	Name string `json:"name" protobuf:"bytes,1,name=name"`
	// Template defines template of action job
	// +optional
	Template *JobTemplateSpec `json:"template,omitempty" protobuf:"bytes,2,opt,name=template"`
	// WorkingDir defines dir to run action, it will always be the git project
	// root dir
	// +optional
	WorkingDir string `json:"workingDir,omitempty" protobuf:"bytes,3,opt,name=workingDir"`
}

const (
	// SystemActionPrefix defines name prefix of system action
	// Now system action contains
	// - system::build-image
	// - system::deploy
	SystemActionPrefix = "system::"
)

// When defines when event triggered
type When string

const (
	// Push means event when git push
	Push When = "git:push"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Pipe defines a pipe which will be triggered by event and generate flow to
// run many jobs
type Pipe struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines desired props of Pipe
	// +optional
	Spec PipeSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status defines status of Pipe
	Status PipeStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// PipeSpec defines spec of Pipe
type PipeSpec struct {
	// Git defines git info
	Git Git `json:"git" protobuf:"bytes,1,opt,name=git"`
	// On defines when pipe will be triggered
	// +optional
	On []When `json:"on,omitempty" protobuf:"bytes,2,rep,name=on"`
	// Stages defines pipe stages which will be run
	// +optional
	Stages []Stage `json:"stages,omitempty" protobuf:"bytes,3,rep,name=stages"`

	// VolumeClaimTemplate defines template of volume to store git code
	// +optional
	VolumeClaimTemplate *corev1.PersistentVolumeClaim `json:"volumeClaimTemplate,omitempty" protobuf:"bytes,4,opt,name=volumeClaimTemplate"`
}

// PipeStatus defines status of pipe
// TODO(liubog2008): add conditions  of pipe
type PipeStatus struct {
	Phase string
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Event defines event which can trigger pipe to generate flow
type Event struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines desired props of Event
	// +optional
	Spec EventSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status defines status of event
	Status EventStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// EventSpec defines event which will trigger some pipes
type EventSpec struct {
	// Repo defines repo of git
	Repo string `json:"repo" protobuf:"bytes,1,opt,name=repo"`
	// When defines when the event triggered
	When When `json:"when" protobuf:"bytes,2,opt,name=when"`
	// Version defines version of git repo
	// It is git ref in fact
	Ref string `json:"version" protobuf:"bytes,3,opt,name=ref"`

	// Extra defines extra info of event
	// It can be used by action env
	// +optional
	Extra map[string]string `json:"extra" protobuf:"bytes,4,opt,name=extra"`
}

// EventStatus defines status of event
type EventStatus struct {
	Phase string `json:"phase"`
}

const (
	// EventPending means event is waiting to be consumed
	EventPending = "Pending"

	// EventConsuming means event is being consuming
	EventConsuming = "Consuming"

	// EventConsumed means event is consumed
	EventConsumed = "Consumed"
)

// Git defines git info
type Git struct {
	// Repo defines git repo
	Repo string `json:"repo" protobuf:"bytes,1,opt,name=repo"`
	// Ref defines git repo ref
	// +optional
	Ref string `json:"ref" protobuf:"bytes,2,opt,name=ref"`
	// GitPullSecret defines secret for git to pull code
	// +optional
	GitPullSecret corev1.LocalObjectReference `json:"gitPullSecret" protobuf:"bytes,3,opt,name=gitPullSecret"`
}

// Stage defines stage of pipe
type Stage struct {
	// Name defines stage name
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Action defines action from mario
	Action string `json:"action" protobuf:"bytes,2,opt,name=action"`
}

const (
	// DefaultFlowRevisionLabel defines label key of flow revision
	DefaultFlowRevisionLabel = "flow.oooops.com/revision"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Flow is a queue of jobs which will be run one by one
type Flow struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines desired props of flow
	// +optional
	Spec FlowSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status defines desired props of flow
	// +optional
	Status FlowStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// FlowSpec defines spec of flow
type FlowSpec struct {
	// Mario defines mario info of flow
	// +optional
	Mario *Mario `json:"mario" protobuf:"mario,1,opt,name=mario"`

	// Git defines git info of flow
	// +optional
	Git Git `json:"git" protobuf:"git,2,opt,name=git"`

	// Stages defines stages of flow
	// +optional
	Stages []Stage `json:"stages,omitempty" protobuf:"bytes,3,rep,name=stages"`

	// VolumeClaim defines pvc referenced by this flow
	VolumeClaim string `json:"volumeClaim,omitempty" protobuf:"bytes,4,opt,name=volumeClaim"`
}

const (
	// FlowPending means flow is pending
	FlowPending = "Pending"
	// FlowRunning means flow is running
	FlowRunning = "Running"
)

// FlowStatus defines status of flow
// TODO(liubog2008): add conditions  of flow
type FlowStatus struct {
	Phase string `josn:"phase,omitempty" protobuf:"bytes,1,opt,name=phase"`
}

// JobTemplateSpec defines template of mario job
type JobTemplateSpec struct {
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines job sepc
	// +optional
	Spec JobSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// JobSpec defines spec of mario job
type JobSpec struct {
	// Containers defines container of job
	Containers []Container `json:"containers,omitempty" protobuf:"bytes,1,rep,name=containers"`
}

// Container defines container of job
type Container struct {
	// Name defines unique key in containers name
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// +optional
	Image string `json:"image,omitempty" protobuf:"bytes,2,opt,name=image"`
	// +optional
	Command []string `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`
	// +optional
	Args []string `json:"args,omitempty" protobuf:"bytes,4,rep,name=args"`
}
