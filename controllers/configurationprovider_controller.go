/*
Copyright 2022 Azure App Configuration.

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

package controllers

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appconfigv1alpha1 "appconfig/sync/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigurationProviderReconciler reconciles a ConfigurationProvider object
type ConfigurationProviderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=appconfig.kubernetes.config,resources=configurationproviders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=appconfig.kubernetes.config,resources=configurationproviders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=appconfig.kubernetes.config,resources=configurationproviders/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigurationProvider object against the actual cluster state, and then
// perforGZDGm operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *ConfigurationProviderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	reqLogger := log.Log.WithValues("namespace", req.Namespace, "ConfigurationProvider", req.Name)
	reqLogger.Info("====== Reconcil ConfigurationProvider")

	instance := &appconfigv1alpha1.ConfigurationProvider{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)

	if err != nil {
		return reconcile.Result{}, nil
	}

	if instance.Status.Phase == "" {
		instance.Status.Phase = appconfigv1alpha1.PhasePending
	}

	if instance.Status.Phase == appconfigv1alpha1.PhasePending {
		reqLogger.Info("====== Reconcil Start..")
	}

	appconfigname := instance.Spec.Endpoint

	endpoint := "https://" + appconfigname + ".azconfig.io"
	clientId := instance.Spec.ClientId
	clientSecret := instance.Spec.ClientId
	tenantId := instance.Spec.TenantId

	reqLogger.Info(endpoint)
	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	client, err := azappconfig.NewClient(endpoint, credential, nil)
	setting, err := client.SetSetting(context.TODO(), "message", nil, nil)

	if err != nil {
		reqLogger.Error(err, "error when get config", nil)
	}

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	//create configmap

	datas := make(map[string]string)
	datas[*setting.Key] = *setting.Value

	configMapToCreate := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            instance.Spec.ConfigMapName,
			Namespace:       "appconfig-dev",
			Labels:          instance.Labels,
			Annotations:     instance.Annotations,
			OwnerReferences: nil,
		},
		Data: datas,
	}

	configmapCreated, err := clientset.CoreV1().ConfigMaps("appconfig-dev").Create(context.TODO(), configMapToCreate, metav1.CreateOptions{})

	//secret, err := secretclient.Get(context.TODO(), "secret-to-be-created-wi", metav1.GetOptions{})

	if err != nil {
		panic(err)
	}

	fmt.Println(configmapCreated.GetName())

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigurationProviderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appconfigv1alpha1.ConfigurationProvider{}).
		Complete(r)
}
