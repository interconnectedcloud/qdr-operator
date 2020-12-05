package interconnect

import (
	"context"
	"github.com/RHsyseng/operator-utils/pkg/olm"
	"reflect"
	"strconv"
	"strings"

	rhseutils "github.com/RHsyseng/operator-utils/pkg/utils/openshift"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/resources/certificates"
	"github.com/interconnectedcloud/qdr-operator/pkg/resources/configmaps"
	"github.com/interconnectedcloud/qdr-operator/pkg/resources/deployments"
	"github.com/interconnectedcloud/qdr-operator/pkg/resources/ingresses"
	"github.com/interconnectedcloud/qdr-operator/pkg/resources/rolebindings"
	"github.com/interconnectedcloud/qdr-operator/pkg/resources/roles"
	"github.com/interconnectedcloud/qdr-operator/pkg/resources/routes"
	"github.com/interconnectedcloud/qdr-operator/pkg/resources/serviceaccounts"
	"github.com/interconnectedcloud/qdr-operator/pkg/resources/services"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/configs"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/random"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/selectors"
	cmv1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extv1b1 "k8s.io/api/extensions/v1beta1"
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

var (
	log                = logf.Log.WithName("controller_interconnect")
	openshift_detected *bool
)

const maxConditions = 6

func isOpenShift() bool {
	if openshift_detected == nil {
		isos, err := rhseutils.IsOpenShift(nil)
		if err != nil {
			log.Error(err, "Failed to detect cluster type")
		}
		openshift_detected = &isos
	}
	return *openshift_detected
}

// Add creates a new Interconnect Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	// TODO(ansmith): verify this is still needed if cert-manager is fully installed
	scheme := mgr.GetScheme()
	if certificates.DetectCertmgrIssuer() {
		utilruntime.Must(cmv1alpha1.AddToScheme(scheme))
		utilruntime.Must(scheme.SetVersionPriority(cmv1alpha1.SchemeGroupVersion))
	}

	if isOpenShift() {
		utilruntime.Must(routev1.AddToScheme(scheme))
		utilruntime.Must(scheme.SetVersionPriority(routev1.SchemeGroupVersion))
	}
	return &ReconcileInterconnect{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("interconnect-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Interconnect
	err = c.Watch(&source.Kind{Type: &v1alpha1.Interconnect{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Deployment and requeue the owner Interconnect
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Interconnect{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Interconnect{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Service and requeue the owner Interconnect
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Interconnect{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource ServiceAccount and requeue the owner Interconnect
	err = c.Watch(&source.Kind{Type: &corev1.ServiceAccount{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Interconnect{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource RoleBinding and requeue the owner Interconnect
	err = c.Watch(&source.Kind{Type: &rbacv1.RoleBinding{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Interconnect{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Secret and requeue the owner Interconnect
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Interconnect{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource ConfigMap and requeue the owner Interconnect
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Interconnect{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner Interconnect
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Interconnect{},
	})
	if err != nil {
		return err
	}

	if certificates.DetectCertmgrIssuer() {
		// Watch for changes to secondary resource Issuer and requeue the owner Interconnect
		err = c.Watch(&source.Kind{Type: &cmv1alpha1.Issuer{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &v1alpha1.Interconnect{},
		})

		// Watch for changes to secondary resource Certificates and requeue the owner Interconnect
		err = c.Watch(&source.Kind{Type: &cmv1alpha1.Certificate{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &v1alpha1.Interconnect{},
		})
	}

	if isOpenShift() {
		// Watch for changes to secondary resource Route and requeue the owner Interconnect
		err = c.Watch(&source.Kind{Type: &routev1.Route{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &v1alpha1.Interconnect{},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileInterconnect{}

// ReconcileInterconnect reconciles a Interconnect object
type ReconcileInterconnect struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

func addCondition(conditions []v1alpha1.InterconnectCondition, condition v1alpha1.InterconnectCondition) []v1alpha1.InterconnectCondition {
	size := len(conditions) + 1
	first := 0
	if size > maxConditions {
		first = size - maxConditions
	}
	return append(conditions, condition)[first:size]
}

// Reconcile reads that state of the cluster for a Interconnect object and makes changes based on the state read
// and what is in the Interconnect.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileInterconnect) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Interconnect")

	// Fetch the Interconnect instance
	instance := &v1alpha1.Interconnect{}
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
		condition := v1alpha1.InterconnectCondition{
			Type:           v1alpha1.InterconnectConditionProvisioning,
			Reason:         "provision spec to desired state",
			TransitionTime: metav1.Now(),
		}
		instance.Status.Conditions = addCondition(instance.Status.Conditions, condition)
		r.client.Status().Update(context.TODO(), instance)
	}

	requestCert, updateDefaults := configs.SetInterconnectDefaults(instance, certificates.DetectCertmgrIssuer())
	if updateDefaults {
		reqLogger.Info("Updating interconnect instance defaults")
		r.client.Update(context.TODO(), instance)
	}

	// Check if role already exists, if not create a new one
	roleFound := &rbacv1.Role{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, roleFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new role
		role := roles.NewRoleForCR(instance)
		controllerutil.SetControllerReference(instance, role, r.scheme)
		reqLogger.Info("Creating a new Role", "role", role)
		err = r.client.Create(context.TODO(), role)
		if err != nil && !errors.IsAlreadyExists(err) {
			reqLogger.Error(err, "Failed to create new Role")
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Role")
		return reconcile.Result{}, err
	}

	// Check if rolebinding already exists, if not create a new one
	rolebindingFound := &rbacv1.RoleBinding{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, rolebindingFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new rolebinding
		rolebinding := rolebindings.NewRoleBindingForCR(instance)
		controllerutil.SetControllerReference(instance, rolebinding, r.scheme)
		reqLogger.Info("Creating a new RoleBinding", "RoleBinding", rolebinding)
		err = r.client.Create(context.TODO(), rolebinding)
		if err != nil && !errors.IsAlreadyExists(err) {
			reqLogger.Error(err, "Failed to create new RoleBinding")
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get RoleBinding")
		return reconcile.Result{}, err
	}

	// Check if serviceaccount already exists, if not create a new one
	svcAccntFound := &corev1.ServiceAccount{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, svcAccntFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new serviceaccount
		svcaccnt := serviceaccounts.NewServiceAccountForCR(instance)
		controllerutil.SetControllerReference(instance, svcaccnt, r.scheme)
		reqLogger.Info("Creating a new ServiceAccount", "ServiceAccount", svcaccnt)
		err = r.client.Create(context.TODO(), svcaccnt)
		if err != nil && !errors.IsAlreadyExists(err) {
			reqLogger.Error(err, "Failed to create new ServiceAccount")
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get ServiceAccount")
		return reconcile.Result{}, err
	}

	if requestCert {
		// If no spec.Issuer, set up a self-signed issuer
		caSecret := instance.Spec.DeploymentPlan.Issuer
		if instance.Spec.DeploymentPlan.Issuer == "" {
			selfSignedIssuerFound := &cmv1alpha1.Issuer{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-selfsigned", Namespace: instance.Namespace}, selfSignedIssuerFound)
			if err != nil && errors.IsNotFound(err) {
				// Define a new selfsigned issuer
				newIssuer := certificates.NewSelfSignedIssuerForCR(instance)
				controllerutil.SetControllerReference(instance, newIssuer, r.scheme)
				reqLogger.Info("Creating a new self signed issuer %s%s\n", newIssuer.Namespace, newIssuer.Name)
				err = r.client.Create(context.TODO(), newIssuer)
				if err != nil && !errors.IsAlreadyExists(err) {
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
				if err != nil && !errors.IsAlreadyExists(err) {
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
			if err != nil && !errors.IsAlreadyExists(err) {
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
			if instance.Spec.SslProfiles[i].GenerateCaCert {
				caCertFound := &cmv1alpha1.Certificate{}
				err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.SslProfiles[i].CaCert, Namespace: instance.Namespace}, caCertFound)
				if err != nil && errors.IsNotFound(err) {
					// Create a new ca certificate
					cert := certificates.NewCACertificateForCR(instance, instance.Spec.SslProfiles[i].CaCert)
					controllerutil.SetControllerReference(instance, cert, r.scheme)
					reqLogger.Info("Creating a new ca cert %s%s\n", cert.Namespace, cert.Name)
					err = r.client.Create(context.TODO(), cert)
					if err != nil && !errors.IsAlreadyExists(err) {
						reqLogger.Info("Failed to create new ca cert", "error", err)
						return reconcile.Result{}, err
					}
				} else if err != nil {
					reqLogger.Info("Failed to create ca cert", "error", err)
					return reconcile.Result{}, err
				}
			}
			if instance.Spec.SslProfiles[i].GenerateCredentials {
				certFound := &cmv1alpha1.Certificate{}
				err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.SslProfiles[i].Credentials, Namespace: instance.Namespace}, certFound)
				if err != nil && errors.IsNotFound(err) {
					var issuerName string
					//if MutualAuth is specified use the CA for the profile to generate the credentials, else use the top level issuer
					if instance.Spec.SslProfiles[i].MutualAuth {
						//ensure we have the necessary issuer
						issuerName = instance.Spec.SslProfiles[i].CaCert
						err = ensureCAIssuer(r, instance, issuerName, instance.Namespace, issuerName)
						if err != nil {
							reqLogger.Info("Failed to reconcile CA issuer", "error", err)
							return reconcile.Result{}, err
						}
						reqLogger.Info("Reconciled CA issuer", "profile", instance.Spec.SslProfiles[i].Name)
					}

					// Create a new certificate
					cert := certificates.NewCertificateForCR(instance, instance.Spec.SslProfiles[i].Name, instance.Spec.SslProfiles[i].Credentials, issuerName)
					controllerutil.SetControllerReference(instance, cert, r.scheme)
					reqLogger.Info("Creating a new cert %s%s\n", cert.Namespace, cert.Name, "issuer", issuerName)
					err = r.client.Create(context.TODO(), cert)
					if err != nil && !errors.IsAlreadyExists(err) {
						reqLogger.Info("Failed to create new cert", "error", err)
						return reconcile.Result{}, err
					}
					// Cert created successfully - set credential return and requeue
					return reconcile.Result{Requeue: true}, nil
				} else if err != nil {
					reqLogger.Info("Failed to create cert", "error", err)
					return reconcile.Result{}, err
				}
			}
		}
	}

	// Check if the deployment already exists, if not create a new one
	if instance.Spec.DeploymentPlan.Placement == v1alpha1.PlacementAny ||
		instance.Spec.DeploymentPlan.Placement == v1alpha1.PlacementAntiAffinity {
		depFound := &appsv1.Deployment{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, depFound)
		if err != nil && errors.IsNotFound(err) {
			// Define a new deployment
			dep := deployments.NewDeploymentForCR(instance)
			controllerutil.SetControllerReference(instance, dep, r.scheme)
			reqLogger.Info("Creating a new Deployment", "Deployment", dep)
			err = r.client.Create(context.TODO(), dep)
			if err != nil && !errors.IsAlreadyExists(err) {
				reqLogger.Error(err, "Failed to create new Deployment")
				return reconcile.Result{}, err
			}
			// update status
			condition := v1alpha1.InterconnectCondition{
				Type:           v1alpha1.InterconnectConditionDeployed,
				Reason:         "deployment created",
				TransitionTime: metav1.Now(),
			}
			instance.Status.Conditions = addCondition(instance.Status.Conditions, condition)
			r.client.Status().Update(context.TODO(), instance)
			// Deployment created successfully - return and requeue
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to get Deployment")
			return reconcile.Result{}, err
		}

		// Ensure the deployment count is the same as the spec size
		// TODO(ansmith): for now, when deployment does not match,
		// delete to recreate pod instances
		size := instance.Spec.DeploymentPlan.Size
		if size != 0 && *depFound.Spec.Replicas != size {
			ct := v1alpha1.InterconnectConditionScalingUp
			if *depFound.Spec.Replicas > size {
				ct = v1alpha1.InterconnectConditionScalingDown
			}
			*depFound.Spec.Replicas = size
			r.client.Update(context.TODO(), depFound)
			// update status
			condition := v1alpha1.InterconnectCondition{
				Type:           ct,
				Reason:         "Instance spec count updated",
				TransitionTime: metav1.Now(),
			}
			instance.Status.Conditions = addCondition(instance.Status.Conditions, condition)
			r.client.Status().Update(context.TODO(), instance)
			return reconcile.Result{Requeue: true}, nil
		} else if !deployments.CheckDeployedContainer(&depFound.Spec.Template, instance) {
			reqLogger.Info("Container config has changed")
			ct := v1alpha1.InterconnectConditionProvisioning
			err = r.client.Update(context.TODO(), depFound)
			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "error", err)
				return reconcile.Result{}, err
			} else {
				reqLogger.Info("Deployment updated", "name", depFound.Name)
			}
			// update status
			condition := v1alpha1.InterconnectCondition{
				Type:           ct,
				Reason:         "Configuration changed",
				TransitionTime: metav1.Now(),
			}
			instance.Status.Conditions = addCondition(instance.Status.Conditions, condition)
			r.client.Status().Update(context.TODO(), instance)
			return reconcile.Result{Requeue: true}, nil
		}
	} else if instance.Spec.DeploymentPlan.Placement == v1alpha1.PlacementEvery {
		dsFound := &appsv1.DaemonSet{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, dsFound)
		if err != nil && errors.IsNotFound(err) {
			// Define a new daemon set
			ds := deployments.NewDaemonSetForCR(instance)
			controllerutil.SetControllerReference(instance, ds, r.scheme)
			reqLogger.Info("Creating a new DaemonSet", "DaemonSet", ds)
			err = r.client.Create(context.TODO(), ds)
			if err != nil && !errors.IsAlreadyExists(err) {
				reqLogger.Error(err, "Failed to create new DaemonSet")
				return reconcile.Result{}, err
			}
			// update status
			condition := v1alpha1.InterconnectCondition{
				Type:           v1alpha1.InterconnectConditionDeployed,
				Reason:         "daemonset created",
				TransitionTime: metav1.Now(),
			}
			instance.Status.Conditions = addCondition(instance.Status.Conditions, condition)
			r.client.Update(context.TODO(), instance)
			// DaemonSet created successfully - return and requeue
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to get DaemonSet")
			return reconcile.Result{}, err
		} else if !deployments.CheckDeployedContainer(&dsFound.Spec.Template, instance) {
			reqLogger.Info("Container config has changed")
			ct := v1alpha1.InterconnectConditionProvisioning
			r.client.Update(context.TODO(), dsFound)
			// update status
			condition := v1alpha1.InterconnectCondition{
				Type:           ct,
				Reason:         "Configuration changed",
				TransitionTime: metav1.Now(),
			}
			instance.Status.Conditions = addCondition(instance.Status.Conditions, condition)
			r.client.Status().Update(context.TODO(), instance)
			return reconcile.Result{Requeue: true}, nil
		}
	} //end of placement is every

	// Check if the service for the deployment already exists, if not create a new one
	svcFound := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, svcFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		svc := services.NewServiceForCR(instance, requestCert)
		controllerutil.SetControllerReference(instance, svc, r.scheme)
		reqLogger.Info("Creating service for interconnect deployment", "Service", svc)
		err = r.client.Create(context.TODO(), svc)
		if err != nil && !errors.IsAlreadyExists(err) {
			reqLogger.Error(err, "Failed to create new Service")
			return reconcile.Result{}, err
		}
		// Service created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	}

	// create route for exposed listeners
	exposedListeners := configs.GetInterconnectExposedListeners(instance)
	for _, listener := range exposedListeners {
		target := listener.Name
		if target == "" {
			target = strconv.Itoa(int(listener.Port))
		}
		if isOpenShift() {
			routeFound := &routev1.Route{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-" + target, Namespace: instance.Namespace}, routeFound)
			if err != nil && errors.IsNotFound(err) {
				// Define a new route
				if listener.SslProfile == "" && !listener.Http {
					// create the route but issue warning
					reqLogger.Info("Warning an exposed listener should be http or ssl enabled", "listener", listener)
				}
				route := routes.NewRouteForCR(instance, listener)
				controllerutil.SetControllerReference(instance, route, r.scheme)
				reqLogger.Info("Creating route for interconnect deployment", "listener", listener)
				err = r.client.Create(context.TODO(), route)
				if err != nil && !errors.IsAlreadyExists(err) {
					reqLogger.Error(err, "Failed to create new Route")
					return reconcile.Result{}, err
				}
				// Route created successfully - return and requeue
				return reconcile.Result{Requeue: true}, nil
			} else if err != nil {
				reqLogger.Error(err, "Failed to get Route")
				return reconcile.Result{}, err
			}
		} else {
			ingressFound := &extv1b1.Ingress{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-" + target, Namespace: instance.Namespace}, ingressFound)
			if err != nil && errors.IsNotFound(err) {
				// Define a new Ingress
				if listener.SslProfile == "" && !listener.Http {
					// create the ingress but issue warning
					reqLogger.Info("Warning an exposed listener should be http or ssl enabled", "listener", listener)
				}
				ingress := ingresses.NewIngressForCR(instance, listener)
				controllerutil.SetControllerReference(instance, ingress, r.scheme)
				reqLogger.Info("Creating Ingress for interconnect deployment", "listener", listener)
				err = r.client.Create(context.TODO(), ingress)
				if err != nil && !errors.IsAlreadyExists(err) {
					reqLogger.Error(err, "Failed to create new Ingress")
					return reconcile.Result{}, err
				}
				// Ingress created successfully - return and requeue
				return reconcile.Result{Requeue: true}, nil
			} else if err != nil {
				reqLogger.Error(err, "Failed to get Ingress")
				return reconcile.Result{}, err
			}
		}
	}

	if !(strings.EqualFold(instance.Spec.Users, "none") || strings.EqualFold(instance.Spec.Users, "false")) {
		// Check if the secret containing users already exists, if not create it
		userSecret := &corev1.Secret{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.Users, Namespace: instance.Namespace}, userSecret)
		if err != nil && errors.IsNotFound(err) {
			labels := selectors.LabelsForInterconnect(instance.Name)
			users := &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Labels:    labels,
					Name:      instance.Spec.Users,
					Namespace: instance.Namespace,
				},
				StringData: map[string]string{"guest": random.GenerateRandomString(8)},
			}

			controllerutil.SetControllerReference(instance, users, r.scheme)
			reqLogger.Info("Creating user secret for interconnect deployment", "Secret", users)
			err = r.client.Create(context.TODO(), users)
			if err != nil && !errors.IsAlreadyExists(err) {
				reqLogger.Error(err, "Failed to create user secret")
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to get Secret")
			return reconcile.Result{}, err
		}

		saslConfig := &corev1.ConfigMap{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-sasl-config", Namespace: instance.Namespace}, saslConfig)
		if err != nil && errors.IsNotFound(err) {
			saslConfig = configmaps.NewConfigMapForSaslConfig(instance)
			controllerutil.SetControllerReference(instance, saslConfig, r.scheme)
			reqLogger.Info("Creating ConfigMap for sasl config", "ConfigMap", saslConfig)
			err = r.client.Create(context.TODO(), saslConfig)
			if err != nil && !errors.IsAlreadyExists(err) {
				reqLogger.Error(err, "Failed to create ConfigMap for sasl config")
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to get ConfigMap for sasl config")
			return reconcile.Result{}, err
		}
	}

	// List the pods for this deployment
	reqLogger.Info("Collecting Pod Status")
	dStatus := getPodStatus(r, instance)
	// Update status.Deployments if needed
	if !reflect.DeepEqual(dStatus, instance.Status.PodStatus) {

		instance.Status.PodStatus = dStatus
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update pod names")
			return reconcile.Result{}, err
		}
		reqLogger.Info("Pod names updated")
		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}

func ensureCAIssuer(r *ReconcileInterconnect, instance *v1alpha1.Interconnect, name string, namespace string, secret string) error {
	caIssuerFound := &cmv1alpha1.Issuer{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, caIssuerFound)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating CA Issuer")
		newIssuer := certificates.NewCAIssuer(name, namespace, secret)
		controllerutil.SetControllerReference(instance, newIssuer, r.scheme)
		return r.client.Create(context.TODO(), newIssuer)
	} else if err != nil {
		return err
	} else if caIssuerFound.Spec.IssuerConfig.CA.SecretName != secret {
		log.Info("Updating CA Issuer")
		caIssuerFound.Spec.IssuerConfig.CA.SecretName = secret
		return r.client.Update(context.TODO(), caIssuerFound)
	}
	return nil
}

func getPodStatus(r *ReconcileInterconnect, instance *v1alpha1.Interconnect) olm.DeploymentStatus {
	// List the pods for this deployment
	var status olm.DeploymentStatus

	if instance.Spec.DeploymentPlan.Placement == v1alpha1.PlacementEvery {
		dsFound := &appsv1.DaemonSet{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, dsFound)
		if err == nil {
			status = olm.GetSingleDaemonSetStatus(*dsFound)
		} else if err != nil {
			log.Error(err, "failed to collect daemonset pod status")
		}
	} else if instance.Spec.DeploymentPlan.Placement == v1alpha1.PlacementAny ||
		instance.Spec.DeploymentPlan.Placement == v1alpha1.PlacementAntiAffinity {
		depFound := &appsv1.Deployment{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, depFound)
		if err == nil {
			status = olm.GetSingleDeploymentStatus(*depFound)
		} else if err != nil {
			log.Error(err, "failed to collect deployment pod status")
		}
	} else {
		log.Info("Unknown placement, unable to determine the pod status", "", instance.Spec.DeploymentPlan.Placement)
	}
	return status
}
