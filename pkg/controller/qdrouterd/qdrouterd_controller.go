package qdrouterd

import (
	"context"
	"reflect"

	v1alpha1 "github.com/interconnectedcloud/qdrouterd-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/resources/certificates"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/resources/configmaps"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/resources/deployments"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/resources/rolebindings"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/resources/roles"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/resources/serviceaccounts"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/resources/services"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/utils/configs"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/utils/selectors"
	cmv1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
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

const maxConditions = 6

// Add creates a new Qdrouterd Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	// TODO(ansmith): verify this is still needed if cert-manager is fully installed
	scheme := mgr.GetScheme()
	utilruntime.Must(cmv1alpha1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(cmv1alpha1.SchemeGroupVersion))

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

	// Watch for changes to secondary resource ServiceAccount and requeue the owner Qdrouterd
	err = c.Watch(&source.Kind{Type: &corev1.ServiceAccount{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Qdrouterd{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource RoleBinding and requeue the owner Qdrouterd
	err = c.Watch(&source.Kind{Type: &rbacv1.RoleBinding{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Qdrouterd{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Secreet and requeue the owner Qdrouterd
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Qdrouterd{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource ConfigMap and requeue the owner Qdrouterd
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Qdrouterd{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner Qdrouterd
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Qdrouterd{},
	})
	if err != nil {
		return err
	}

	// TODO(ansmith): Check if there is a cert-manager crd instance, handle err
	// Watch for changes to secondary resource Issuer and requeue the owner Qdrouterd
	err = c.Watch(&source.Kind{Type: &cmv1alpha1.Issuer{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Qdrouterd{},
	})

	// Watch for changes to secondary resource Certificates and requeue the owner Qdrouterd
	err = c.Watch(&source.Kind{Type: &cmv1alpha1.Certificate{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Qdrouterd{},
	})

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

func addCondition(conditions []v1alpha1.QdrouterdCondition, condition v1alpha1.QdrouterdCondition) []v1alpha1.QdrouterdCondition {
	size := len(conditions) + 1
	first := 0
	if size > maxConditions {
		first = size - maxConditions
	}
	return append(conditions, condition)[first:size]
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

	// Assign the generated resource version to the status
	if instance.Status.RevNumber == "" {
		instance.Status.RevNumber = instance.ObjectMeta.ResourceVersion
		// update status
		condition := v1alpha1.QdrouterdCondition{
			Type:           v1alpha1.QdrouterdConditionProvisioning,
			Reason:         "provision spec to desired state",
			TransitionTime: metav1.Now(),
		}
		instance.Status.Conditions = addCondition(instance.Status.Conditions, condition)
		r.client.Update(context.TODO(), instance)
	}

	requestCert := configs.SetQdrouterdDefaults(instance)

	// Check if role already exists, if not create a new one
	roleFound := &rbacv1.Role{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, roleFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new role
		role := roles.NewRoleForCR(instance)
		controllerutil.SetControllerReference(instance, role, r.scheme)
		reqLogger.Info("Creating a new Role %s%s\n", role.Namespace, role.Name)
		err = r.client.Create(context.TODO(), role)
		if err != nil {
			reqLogger.Info("Failed to create new Role: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Info("Failed to get Role: %v\n", err)
		return reconcile.Result{}, err
	}

	// Check if rolebinding already exists, if not create a new one
	rolebindingFound := &rbacv1.RoleBinding{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, rolebindingFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new rolebinding
		rolebinding := rolebindings.NewRoleBindingForCR(instance)
		controllerutil.SetControllerReference(instance, rolebinding, r.scheme)
		reqLogger.Info("Creating a new RoleBinding %s%s\n", rolebinding.Namespace, rolebinding.Name)
		err = r.client.Create(context.TODO(), rolebinding)
		if err != nil {
			reqLogger.Info("Failed to create new RoleBinding: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Info("Failed to get RoleBinding: %v\n", err)
		return reconcile.Result{}, err
	}

	// Check if serviceaccount already exists, if not create a new one
	svcAccntFound := &corev1.ServiceAccount{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, svcAccntFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new serviceaccount
		svcaccnt := serviceaccounts.NewServiceAccountForCR(instance)
		controllerutil.SetControllerReference(instance, svcaccnt, r.scheme)
		reqLogger.Info("Creating a new ServiceAccount %s%s\n", svcaccnt.Namespace, svcaccnt.Name)
		err = r.client.Create(context.TODO(), svcaccnt)
		if err != nil {
			reqLogger.Info("Failed to create new ServiceAccount: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Info("Failed to get ServiceAccount: %v\n", err)
		return reconcile.Result{}, err
	}

	if requestCert {
		// If no spec.Issuer, set up a self-signed issuer
		caSecret := instance.Spec.Issuer
		if instance.Spec.Issuer == "" {
			selfSignedIssuerFound := &cmv1alpha1.Issuer{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-selfsigned", Namespace: instance.Namespace}, selfSignedIssuerFound)
			if err != nil && errors.IsNotFound(err) {
				// Define a new selfsigned issuer
				newIssuer := certificates.NewSelfSignedIssuerForCR(instance)
				controllerutil.SetControllerReference(instance, newIssuer, r.scheme)
				reqLogger.Info("Creating a new self signed issuer %s%s\n", newIssuer.Namespace, newIssuer.Name)
				err = r.client.Create(context.TODO(), newIssuer)
				if err != nil {
					reqLogger.Info("Failed to create new self signed issuer", "error", err)
					return reconcile.Result{}, err
				}
				// Issuer created successfully - return and requeue
				return reconcile.Result{Requeue: true}, nil
			} else if err != nil {
				reqLogger.Info("Failed to get self signed issuer", "error", err)
				return reconcile.Result{}, err
			}

			selfSignedCertFound := &cmv1alpha1.Certificate{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-selfsigned", Namespace: instance.Namespace}, selfSignedCertFound)
			if err != nil && errors.IsNotFound(err) {
				// Create a new self signed certificate
				cert := certificates.NewSelfSignedCACertificateForCR(instance)
				controllerutil.SetControllerReference(instance, cert, r.scheme)
				reqLogger.Info("Creating a new self signed cert %s%s\n", cert.Namespace, cert.Name)
				err = r.client.Create(context.TODO(), cert)
				if err != nil {
					reqLogger.Info("Failed to create new self signed cert", "error", err)
					return reconcile.Result{}, err
				}
				// Cert created successfully - return and requeue
				return reconcile.Result{Requeue: true}, nil
			} else if err != nil {
				reqLogger.Info("Failed to create self signed cert", "error", err)
				return reconcile.Result{}, err
			}
			caSecret = selfSignedCertFound.Name
		}

		// Check if CA issuer exists and if not create one
		caIssuerFound := &cmv1alpha1.Issuer{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-ca", Namespace: instance.Namespace}, caIssuerFound)
		if err != nil && errors.IsNotFound(err) {
			// Define a new ca issuer
			newIssuer := certificates.NewCAIssuerForCR(instance, caSecret)
			controllerutil.SetControllerReference(instance, newIssuer, r.scheme)
			reqLogger.Info("Creating a new ca issuer %s%s\n", newIssuer.Namespace, newIssuer.Name)
			err = r.client.Create(context.TODO(), newIssuer)
			if err != nil {
				reqLogger.Info("Failed to create new ca issuer", "error", err)
				return reconcile.Result{}, err
			}
			// Issuer created successfully - return and requeue
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Info("Failed to get ca issuer", "error", err)
			return reconcile.Result{}, err
		}

		// As needed, create certs for SslProfiles
		for i := range instance.Spec.SslProfiles {
			if instance.Spec.SslProfiles[i].Credentials == "" {
				certFound := &cmv1alpha1.Certificate{}
				err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-" + instance.Spec.SslProfiles[i].Name + "-tls", Namespace: instance.Namespace}, certFound)
				if err != nil && errors.IsNotFound(err) {
					// Create a new certificate
					cert := certificates.NewCertificateForCR(instance, instance.Spec.SslProfiles[i].Name)
					controllerutil.SetControllerReference(instance, cert, r.scheme)
					reqLogger.Info("Creating a new cert %s%s\n", cert.Namespace, cert.Name)
					err = r.client.Create(context.TODO(), cert)
					if err != nil {
						reqLogger.Info("Failed to create new cert", "error", err)
						return reconcile.Result{}, err
					}
					// Cert created successfully - set credential return and requeue
					instance.Spec.SslProfiles[i].Credentials = instance.Name + "-" + instance.Spec.SslProfiles[i].Name + "-tls"
					r.client.Update(context.TODO(), instance)
					return reconcile.Result{Requeue: true}, nil
				} else if err != nil {
					reqLogger.Info("Failed to create cert", "error", err)
					return reconcile.Result{}, err
				}
			}
			if instance.Spec.SslProfiles[i].RequireClientCerts && instance.Spec.SslProfiles[i].CaCert == "" {
				caCertFound := &cmv1alpha1.Certificate{}
				err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-" + instance.Spec.SslProfiles[i].Name + "-ca", Namespace: instance.Namespace}, caCertFound)
				if err != nil && errors.IsNotFound(err) {
					// Create a new ca certificate
					cert := certificates.NewCACertificateForCR(instance, instance.Spec.SslProfiles[i].Name)
					controllerutil.SetControllerReference(instance, cert, r.scheme)
					reqLogger.Info("Creating a new ca cert %s%s\n", cert.Namespace, cert.Name)
					err = r.client.Create(context.TODO(), cert)
					if err != nil {
						reqLogger.Info("Failed to create new ca cert", "error", err)
						return reconcile.Result{}, err
					}
					// ca cert created successfully - set cacert return and requeue
					instance.Spec.SslProfiles[i].CaCert = instance.Name + "-" + instance.Spec.SslProfiles[i].Name + "-ca"
					r.client.Update(context.TODO(), instance)
					return reconcile.Result{Requeue: true}, nil
				} else if err != nil {
					reqLogger.Info("Failed to create ca cert", "error", err)
					return reconcile.Result{}, err
				}
			}
		}
	}

	// Check if configmap already exists, if not create a new one
	cfgmapFound := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, cfgmapFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new configmap
		cfgmap := configmaps.NewConfigMapForCR(instance)
		controllerutil.SetControllerReference(instance, cfgmap, r.scheme)
		reqLogger.Info("Creating a new ConfigMap %s%s\n", cfgmap.Namespace, cfgmap.Name)
		err = r.client.Create(context.TODO(), cfgmap)
		if err != nil {
			reqLogger.Info("Failed to create new ConfigMap: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Info("Failed to get ConfigMap: %v\n", err)
		return reconcile.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	depFound := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, depFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := deployments.NewDeploymentForCR(instance)
		controllerutil.SetControllerReference(instance, dep, r.scheme)
		reqLogger.Info("Creating a new Deployment %s%s\n", dep.Namespace, dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Info("Failed to create new Deployment: %v\n", err)
			return reconcile.Result{}, err
		}
		// update status
		condition := v1alpha1.QdrouterdCondition{
			Type:           v1alpha1.QdrouterdConditionDeployed,
			Reason:         "deployment created",
			TransitionTime: metav1.Now(),
		}
		instance.Status.Conditions = addCondition(instance.Status.Conditions, condition)
		r.client.Update(context.TODO(), instance)
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Info("Failed to get Deployment: %v\n", err)
		return reconcile.Result{}, err
	}

	// Ensure the deployment count is the same as the spec size
	// TODO(ansmith): for now, when deployment does not match,
	// delete to recreate pod instances
	count := instance.Spec.Count
	if count != 0 && *depFound.Spec.Replicas != count {
		ct := v1alpha1.QdrouterdConditionScalingUp
		if *depFound.Spec.Replicas > count {
			ct = v1alpha1.QdrouterdConditionScalingDown
		}
		*depFound.Spec.Replicas = count
		r.client.Update(context.TODO(), depFound)
		// update status
		condition := v1alpha1.QdrouterdCondition{
			Type:           ct,
			Reason:         "Instance spec count updated",
			TransitionTime: metav1.Now(),
		}
		instance.Status.Conditions = addCondition(instance.Status.Conditions, condition)
		instance.Status.PodNames = instance.Status.PodNames[:0]
		r.client.Update(context.TODO(), instance)
		return reconcile.Result{Requeue: true}, nil
	}

	// Check if the service for the deployment already exists, if not create a new one
	svcFound := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, svcFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		svc := services.NewServiceForCR(instance, requestCert)
		controllerutil.SetControllerReference(instance, svc, r.scheme)
		reqLogger.Info("Creating service for qdrouterd deployment")
		err = r.client.Create(context.TODO(), svc)
		if err != nil {
			reqLogger.Info("Failed to create new Service: %v\n", err)
			return reconcile.Result{}, err
		}
		// Service created successfully - return and requeue
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
		if pod.GetObjectMeta().GetDeletionTimestamp() == nil {
			podNames = append(podNames, pod.Name)
		}
	}
	return podNames
}
