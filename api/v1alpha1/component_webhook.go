/*
Copyright 2022-2023 Red Hat, Inc.

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
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var componentlog = logf.Log.WithName("component-resource")

func (r *Component) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-appstudio-redhat-com-v1alpha1-component,mutating=true,failurePolicy=fail,sideEffects=None,groups=appstudio.redhat.com,resources=components,verbs=create;update,versions=v1alpha1,name=mcomponent.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Component{}

// Github is the only supported vendor right now
const SupportedGitRepo = "github.com"

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Component) Default() {

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-appstudio-redhat-com-v1alpha1-component,mutating=false,failurePolicy=fail,sideEffects=None,groups=appstudio.redhat.com,resources=components,verbs=create;update,versions=v1alpha1,name=vcomponent.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Component{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Component) ValidateCreate() error {
	componentlog.Info("validating the create request", "name", r.Name)

	// We use the DNS-1035 format for component names, so ensure it conforms to that specification
	if len(validation.IsDNS1035Label(r.Name)) != 0 {
		return fmt.Errorf(InvalidDNS1035Name, r.Name)
	}
	sourceSpecified := false

	if r.Spec.Source.GitSource != nil && r.Spec.Source.GitSource.URL != "" {
		if gitsourceURL, err := url.ParseRequestURI(r.Spec.Source.GitSource.URL); err != nil {
			return fmt.Errorf(err.Error() + InvalidSchemeGitSourceURL)
		} else if SupportedGitRepo != strings.ToLower(gitsourceURL.Host) {
			return fmt.Errorf(InvalidGithubVendorURL, gitsourceURL, SupportedGitRepo)
		}

		sourceSpecified = true
	} else if r.Spec.ContainerImage != "" {
		sourceSpecified = true
	}

	if !sourceSpecified {
		return errors.New(MissingGitOrImageSource)
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Component) ValidateUpdate(old runtime.Object) error {
	componentlog.Info("validating the update request", "name", r.Name)

	switch old := old.(type) {
	case *Component:

		if r.Spec.ComponentName != old.Spec.ComponentName {
			return fmt.Errorf(ComponentNameUpdateError, r.Spec.ComponentName)
		}

		if r.Spec.Application != old.Spec.Application {
			return fmt.Errorf(ApplicationNameUpdateError, r.Spec.Application)
		}

		if r.Spec.Source.GitSource != nil && old.Spec.Source.GitSource != nil && !reflect.DeepEqual(*(r.Spec.Source.GitSource), *(old.Spec.Source.GitSource)) {
			return fmt.Errorf(GitSourceUpdateError, *(r.Spec.Source.GitSource))
		}
	default:
		return errors.New(InvalidComponentError)
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Component) ValidateDelete() error {

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
