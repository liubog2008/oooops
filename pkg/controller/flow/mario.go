package flow

import (
	"fmt"
	"io/ioutil"
	"net/http"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	"github.com/liubog2008/oooops/pkg/client/clientset/scheme"
)

func (c *Controller) attachMario(flow *v1alpha1.Flow, jobMap map[string]*batchv1.Job) (bool, error) {
	// if mario is attched, don't attach again
	if flow.Spec.Mario != nil {
		return false, nil
	}
	// if mario job is missing, can't attach mario
	marioJob, ok := jobMap[v1alpha1.FlowStageMario]
	if !ok {
		return false, nil
	}
	selector, err := metav1.LabelSelectorAsSelector(marioJob.Spec.Selector)
	if err != nil {
		return false, err
	}
	pods, err := c.podLister.Pods(flow.Namespace).List(selector)
	if err != nil {
		return false, err
	}
	var mario *v1alpha1.Mario
	for _, pod := range pods {
		if IsPodReady(pod) && metav1.IsControlledBy(pod, marioJob) {
			m, err := c.fetchMario(pod.Status.PodIP)
			if err != nil {
				klog.Warningf("can't fetch mario from %s: %s", pod.Status.PodIP, err)
				continue
			}
			mario = m
		}
	}
	if mario == nil {
		// TODO(liubog2008): do something to recover
		return false, nil
	}
	flow.Spec.Mario = mario
	if _, err := c.extClient.MarioV1alpha1().Flows(flow.Namespace).Update(flow); err != nil {
		return false, err
	}

	return true, nil
}

func (c *Controller) fetchMario(ip string) (*v1alpha1.Mario, error) {
	req, err := http.NewRequest("GET", "http://"+ip+":8080", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("ContentType", "application/json")
	req.Header.Set("Authorization", "Bearer test")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't fetch mario [%s]: %s", resp.StatusCode, body)
	}

	m := v1alpha1.Mario{}

	decoder := scheme.Codecs.UniversalDecoder(v1alpha1.SchemeGroupVersion)

	if _, _, err := decoder.Decode(body, nil, &m); err != nil {
		return nil, err
	}

	return &m, nil
}
