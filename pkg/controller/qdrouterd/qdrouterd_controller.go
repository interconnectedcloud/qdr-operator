package qdrouterd

import (
	"context"
	"reflect"

	v1alpha1 "github.com/interconnectedcloud/qdrouterd-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/resources/deployments"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/resources/services"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/utils/configs"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/utils/selectors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_qdrouterd")

// Add creates a new Qdrouterd Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileQdrouterd{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("qdrouterd-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Qdrouterd
	err = c.Watch(&source.Kind{Type: &v1alpha1.Qdrouterd{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Deployment and requeue the owner Qdrouterd
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Qdrouterd{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Service and requeue the owner Qdrouterd
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Qdrouterd{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileQdrouterd{}

// ReconcileQdrouterd reconciles a Qdrouterd object
type ReconcileQdrouterd struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Qdrouterd object and makes changes based on the state read
// and what is in the Qdrouterd.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileQdrouterd) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Qdrouterd")

	// Fetch the Qdrouterd instance
	instance := &v1alpha1.Qdrouterd{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	requestCert := configs.SetQdrouterdDefaults(instance)

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := deployments.NewDeploymentForCR(instance)
		reqLogger.Info("Creating a new Deployment %s%s\n", dep.Namespace, dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Info("Failed to create new Deployment: %v\n", err)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		controllerutil.SetControllerReference(instance, dep, r.scheme)
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Info("Failed to get Deployment: %v\n", err)
		return reconcile.Result{}, err
	}

	// Ensure the deployment count is the same as the spec size
	count := instance.Spec.Count
	if *found.Spec.Replicas != count {
		found.Spec.Replicas = &count
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			reqLogger.Info("Failed to update Deployment: %v\n", err)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		reqLogger.Info("Deployment count updated")
		return reconcile.Result{Requeue: true}, nil
	}

	// TODO(ansmith): until the qdrouterd can re-read config, this can wait
	// Ensure the deployment container matches the instance spec

	// Check if the service for the deployment already exists, if not create a new one
	locate := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, locate)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		svc := services.NewServiceForCR(instance, requestCert)
		reqLogger.Info("Creating service for qdrouterd deployment")
		err = r.client.Create(context.TODO(), svc)
		if err != nil {
			reqLogger.Info("Failed to create new Service: %v\n", err)
			return reconcile.Result{}, err
		}
		// Service created successfully - return and requeue
		controllerutil.SetControllerReference(instance, svc, r.scheme)
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Info("Failed to get Service: %v\n", err)
		return reconcile.Result{}, err
	}

	// List the pods for this deployment
	podList := &corev1.PodList{}
	labelSelector := selectors.ResourcesByQdrouterdName(instance.Name)
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
	if err != nil {
		reqLogger.Info("Failed to list pods: %v\n", err)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.PodNames if needed
	if !reflect.DeepEqual(podNames, instance.Status.PodNames) {
		instance.Status.PodNames = podNames
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info("Failed to update pod names: %v\n", err)
			return reconcile.Result{}, err
		}
		reqLogger.Info("Pod names updated")
		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
