package flow

// import (
// 	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
// 	batchv1 "k8s.io/api/batch/v1"
// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )
//
// const (
// 	repoPath = "/gitrepo"
// )
//
// func (c *Controller) RunGitAction(pipe *v1alpha1.Pipe, event *v1alpha1.Event) (string, error) {
// 	c.GeneratePVC(event.Name, pipe)
// 	c.GenerateMarioJob(pipe, event)
// 	return "", nil
// }
//
// func (c *Controller) GeneratePVC(name string, pipe *v1alpha1.Pipe) (*corev1.PersistentVolumeClaim, error) {
// 	pvc := corev1.PersistentVolumeClaim{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      name,
// 			Namespace: pipe.Namespace,
// 		},
// 	}
//
// 	return &pvc, nil
// }
//
// func (c *Controller) GenerateMarioJob(pipe *v1alpha1.Pipe, event *v1alpha1.Event) (*batchv1.Job, error) {
// 	job := batchv1.Job{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      event.Name,
// 			Namespace: event.Namespace,
// 		},
// 		Spec: batchv1.JobSpec{
// 			Template: corev1.PodTemplateSpec{
// 				ObjectMeta: metav1.ObjectMeta{},
// 				Spec: corev1.PodSpec{
// 					Containers: []corev1.Container{
// 						{
// 							Name:    "mario",
// 							Image:   c.marioImage,
// 							Command: c.genMarioCmd(event),
// 							VolumeMounts: []corev1.VolumeMount{
// 								{
// 									Name:      "data",
// 									MountPath: repoPath,
// 								},
// 							},
// 						},
// 					},
// 					Volumes: []corev1.Volume{
// 						{
// 							Name: "data",
// 							VolumeSource: corev1.VolumeSource{
// 								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
// 									ClaimName: event.Name,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	return &job, nil
// }
//
// func randomName(pipe *v1alpha1.Pipe) string {
// 	return pipe.Name
// }
//
// func (c *Controller) genMarioCmd(event *v1alpha1.Event) []string {
// 	cmd := []string{
// 		"/app/mario",
// 		"--ref",
// 		event.Spec.Version,
// 		"--remote-path",
// 		event.Spec.Repo,
// 		"--root-path",
// 		"/data",
// 	}
// 	return cmd
// }
