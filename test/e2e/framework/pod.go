package framework

import (
	"context"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// WaitForPodStatus waits for given podName to be available with a matching PodPhase
//                  or it returns a timeout.
func (f *Framework) WaitForPodStatus(podName string, status v1.PodPhase, timeout time.Duration, interval time.Duration) (*v1.Pod, error) {

	var pod *v1.Pod
	var err error

	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()
	err = RetryWithContext(ctx, interval, func() (bool, error) {
		pod, err = f.KubeClient.CoreV1().Pods(f.Namespace).Get(podName, metav1.GetOptions{})
		if err != nil {
			// pod does not exist yet
			return false, nil
		}
		return pod.Status.Phase == status, nil
	})

	return pod, err
}

// GetInterconnectPods returns all pods for the given interconnect instance
func (f *Framework) GetInterconnectPods(ic *v1alpha1.Interconnect) ([]v1.Pod, error) {
	options := metav1.ListOptions{LabelSelector: "application=" + ic.Name + ",interconnect_cr=" + ic.Name}
	podList, err := f.KubeClient.CoreV1().Pods(ic.Namespace).List(options)
	return podList.Items, err
}

// GetInterconnectPodNames returns all pod names for the given interconnect instance
func (f *Framework) GetInterconnectPodNames(ic *v1alpha1.Interconnect) ([]string, error) {
	var podNames []string
	pods, err := f.GetInterconnectPods(ic)
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames, err
}
