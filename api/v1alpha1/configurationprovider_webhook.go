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

package v1alpha1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var configurationproviderlog = logf.Log.WithName("configurationprovider-resource")

func (r *ConfigurationProvider) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-appconfig-kubernetes-config-v1alpha1-configurationprovider,mutating=true,failurePolicy=fail,sideEffects=None,groups=appconfig.kubernetes.config,resources=configurationproviders,verbs=create;update,versions=v1alpha1,name=mconfigurationprovider.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &ConfigurationProvider{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ConfigurationProvider) Default() {
	configurationproviderlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-appconfig-kubernetes-config-v1alpha1-configurationprovider,mutating=false,failurePolicy=fail,sideEffects=None,groups=appconfig.kubernetes.config,resources=configurationproviders,verbs=create;update,versions=v1alpha1,name=vconfigurationprovider.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &ConfigurationProvider{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ConfigurationProvider) ValidateCreate() error {
	configurationproviderlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.

	if r.Spec.TenantId == "" {
		var allErrs field.ErrorList
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec"), "tenantId", "TenantId field is not specified"))
		return apierrors.NewInvalid(schema.GroupKind{}, r.Name, allErrs)
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ConfigurationProvider) ValidateUpdate(old runtime.Object) error {
	configurationproviderlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ConfigurationProvider) ValidateDelete() error {
	configurationproviderlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
