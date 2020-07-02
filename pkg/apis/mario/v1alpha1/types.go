package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// MarioFile defines file of Mario API in git project
	MarioFile = ".mario.yaml"
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
	// Imports defines import path of external mario action
	Imports []string `json:"imports,omitempty" protobuf:"bytes,1,rep,name=imports"`
	// Actions defines actions of the project
	// e.g. compile, test
	// +optional
	Actions []MarioAction `json:"actions,omitempty" protobuf:"bytes,2,rep,name=actions"`
}

// VersionDefinition defines info of git version
type VersionDefinition struct {
	// EnvName defines name of version env
	EnvName string `json:"envName" protobuf:"bytes,1,opt,name=envName"`
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
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Repo",type=string,JSONPath=`.spec.git.repo`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
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
	// Label selector for pods. Existing ReplicaSets whose pods are
	// selected by this will be the ones affected by this deployment.
	// It must match the pod template's labels.
	// +optional
	// +nullable
	Selector *metav1.LabelSelector `json:"selector" protobuf:"bytes,1,opt,name=selector"`
	// When defines when pipe will be triggered
	// +optional
	When []When `json:"when,omitempty" protobuf:"bytes,2,rep,name=when"`

	// Git defines git info
	Git Git `json:"git" protobuf:"bytes,3,opt,name=git"`
	// Stages defines pipe stages which will be run
	// +optional
	Stages []Stage `json:"stages,omitempty" protobuf:"bytes,4,rep,name=stages"`
}

// PipeStatus defines status of pipe
// TODO(liubog2008): add conditions  of pipe
type PipeStatus struct {
	// Phase defines phase of pipe
	Phase string `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Event defines event which can trigger pipe to generate flow
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Repo",type=string,JSONPath=`.spec.git.repo`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
type Event struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines desired props of Event
	// +optional
	Spec EventSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// EventSpec defines event which will trigger some pipes
type EventSpec struct {
	// Repo defines repo of git
	Repo string `json:"repo" protobuf:"bytes,1,opt,name=repo"`
	// When defines when the event triggered
	When When `json:"when" protobuf:"bytes,2,opt,name=when"`
	// Ref defines version of git repo
	// e.g. pull/11/head
	Ref string `json:"ref" protobuf:"bytes,3,opt,name=ref"`

	// Extra defines extra info of event
	// It can be used by action env
	// +optional
	Extra map[string]string `json:"extra" protobuf:"bytes,4,opt,name=extra"`
}

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
	// VolumeClaimTemplate defines template of volume to store git code
	// nolint: lll
	// +optional
	VolumeClaimTemplate *corev1.PersistentVolumeClaim `json:"volumeClaimTemplate,omitempty" protobuf:"bytes,4,opt,name=volumeClaimTemplate"`
}

// Stage defines stage of pipe
type Stage struct {
	// Name defines stage name
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Action defines action from mario
	Action string `json:"action" protobuf:"bytes,2,opt,name=action"`
}

const (
	// DefaultFlowRevisionLabelKey defines label key of flow revision
	DefaultFlowRevisionLabelKey = "flow.oooops.com/revision"

	// DefaultFlowStageLabelKey defines label key of flow stage label
	DefaultFlowStageLabelKey = "flow.oooops.com/stage"

	FlowStageGit   = "git"
	FlowStageMario = "mario"
)

const (
	UserJobPrefix = "user-"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Flow is a queue of jobs which will be run one by one
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
type Flow struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines desired props of flow
	// +optional
	Spec FlowSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status defines desired props of flow
	// +optional
	// +kubebuilder:default={phase:"Pending"}
	Status FlowStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// FlowSpec defines spec of flow
type FlowSpec struct {
	// Label selector for pods. Existing ReplicaSets whose pods are
	// selected by this will be the ones affected by this deployment.
	// It must match the pod template's labels.
	Selector *metav1.LabelSelector `json:"selector" protobuf:"bytes,1,opt,name=selector"`
	// Mario defines mario info of flow
	// +optional
	// +nullable
	Mario *Mario `json:"mario,omitempty" protobuf:"mario,2,opt,name=mario"`

	// Git defines git info of flow
	// +optional
	Git Git `json:"git" protobuf:"git,3,opt,name=git"`

	// Stages defines stages of flow
	// +optional
	Stages []Stage `json:"stages,omitempty" protobuf:"bytes,4,rep,name=stages"`
}

const (
	// FlowPending means flow is pending
	FlowPending = "Pending"
	// FlowRunning means flow is running
	FlowRunning = "Running"
	// FlowSucceed means flow has succeeded
	FlowSucceed = "Succeeded"
	// FlowFailed means flow has failed
	FlowFailed = "Failed"
)

// FlowStatus defines status of flow
// TODO(liubog2008): add conditions  of flow
type FlowStatus struct {
	// Phase of flow
	Phase string `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase"`
	// Stages of flow
	StageStatuses []StageStatus `json:"stageStatuses,omitempty" protobuf:"bytes,2,rep,name=stageStatuses"`
	// Conditions defines condition of flow
	Conditions []FlowCondition `json:"conditions,omitempty" protobuf:"bytes,3,rep,name=conditions"`
}

// FlowConditionType defines type of flow condition
type FlowConditionType string

const (
	// FlowGitVolumeReady means git volume is created and bounded
	FlowGitVolumeReady FlowConditionType = "GitVolumeReady"

	// FlowMarioReady means repo code is fetched and mario is attached
	FlowMarioReady FlowConditionType = "MarioReady"
)

const (
	// FlowReasonGitVolumeClaiming means git volume claim is creating
	FlowReasonGitVolumeClaiming = "GitVolumeClaiming"
	// FlowReasonGitVolumePending means git volume is pending
	FlowReasonGitVolumePending = "GitVolumePending"
	// FlowReasonGitVolumeBound means git volume is bound
	FlowReasonGitVolumeBound = "GitVolumeBound"
	// FlowReasonGitVolumeLost means git volume is lost
	FlowReasonGitVolumeLost = "GitVolumeLost"
	// FlowReasonGitVolumeUnknown means git volume status is unknown
	FlowReasonGitVolumeUnknown = "GitVolumeUnknwon"

	// FlowReasonMarioFailed means mario can not be attached
	FlowReasonMarioFailed = "MarioFailed"
	// FlowReasonMarioPending means mario is waiting for attching
	FlowReasonMarioPending = "MarioPending"
	// FlowReasonMarioReady means mario is ready
	FlowReasonMarioReady = "MarioReady"

	// FlowReasonGitFailed means git job is failed
	FlowReasonGitFailed = "GitFailed"
	// FlowReasonGitPending means git job is waiting for starting
	FlowReasonGitPending = "GitPending"
)

// FlowCondition defines condition of flow
type FlowCondition struct {
	// Type is the type of the condition.
	Type FlowConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=PodConditionType"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	// Last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty" protobuf:"bytes,3,opt,name=lastProbeTime"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

const (
	// StageJobMissing means job is missing, which will be generated when
	StageJobMissing = "JobMissing"
	// StageJobComplete means job is completed
	StageJobComplete = "JobComplete"
	// StageJobFailed means job is failed
	StageJobFailed = "JobFailed"
	// StageJobRunning means job is running
	StageJobRunning = "JobRunning"
)

// StageStatus means status of each stage of flow
type StageStatus struct {
	// Job of current stage
	Job string `json:"job,omitempty" protobuf:"bytes,1,opt,name=job"`
	// Phase of stage
	Phase string `json:"phase,omitempty" protobuf:"bytes,2,opt,name=phase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Action defines an external action which can be imported by mario
type Action struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines desired props of Action
	// +optional
	Spec ActionSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type ActionSpec struct {
	Template *ActionTemplate `json:"template"`

	Args []ActionArg `json:"args"`
}

type ActionArg struct {
	Name        string `json:"name"`
	Optional    bool   `json:"optional"`
	Description string `json:"description"`
}

type MarioAction struct {
	// Name defines name of action
	Name string `json:"name" protobuf:"bytes,1,name=name"`
	// Template defines action template, if action is an imported one, this field will be ignored
	Template *ActionTemplate `json:"template,omitempty" protobuf:"bytes,2,opt,name=template"`

	Env []ActionEnvVar `json:"envs,omitempty" protobuf:"rep,3,opt,name=envs"`

	Secrets []ActionSecret `json:"secrets,omitempty" protobuf:"rep,4,opt,name=version"`

	ServiceAccountName string `json:"serviceAccountName,omitempty" protobuf:"bytes,5,opt,name=serviceAccountName"`
}

type ActionTemplate struct {
	// +optional
	Image string `json:"image,omitempty" protobuf:"bytes,2,opt,name=image"`
	// +optional
	Command []string `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`
	// +optional
	Args []string `json:"args,omitempty" protobuf:"bytes,4,rep,name=args"`

	// WorkingDir defines dir to run action, it will always be the git project
	// root dir
	// +optional
	WorkingDir string `json:"workingDir,omitempty" protobuf:"bytes,5,opt,name=workingDir"`
	// Version defines info of git version
	// +optional
	Version VersionDefinition `json:"version,omitempty" protobuf:"bytes,6,opt,name=version"`
}

// ActionEnvVar defines env variable of action
type ActionEnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ActionSecret struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}
