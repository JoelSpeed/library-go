package resourceapply

import (
	"github.com/openshift/library-go/pkg/operator/events"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
)

// ApplyMutatingWebhookConfiguration ensures the form of the specified
// mutatingwebhookconfiguration is present in the API. If it does not exist,
// it will be created. If it does exist, the metadata of the required
// mutatingwebhookconfiguration will be merged with the existing mutatingwebhookconfiguration
// and an update performed if the mutatingwebhookconfiguration spec and metadata differ from
// the previously required spec and metadata. For further detail, check the top-level comment.
func ApplyMutatingWebhookConfiguration(client dynamic.Interface, recorder events.Recorder,
	requiredOriginal *admissionregistrationv1.MutatingWebhookConfiguration, expectedGeneration int64) (*admissionregistrationv1.MutatingWebhookConfiguration, bool, error) {

	// Explcitily specify type for requiredUnstr to get object meta accessor
	requiredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(requiredOriginal)
	if err != nil {
		return nil, false, err
	}

	gvr := admissionregistrationv1.SchemeGroupVersion.WithResource("mutatingwebhookconfigurations")
	requiredUnstr := &unstructured.Unstructured{Object: requiredObj}
	resourcedClient := client.Resource(gvr)
	actualUnstr, modified, err := applyUnstructured(resourcedClient, "MutatingWebhookConfiguration", recorder, requiredUnstr, expectedGeneration)
	if err != nil {
		return nil, modified, err
	}

	actual := &admissionregistrationv1.MutatingWebhookConfiguration{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(actualUnstr.Object, actual); err != nil {
		return nil, modified, err
	}
	return actual, modified, nil
}

// ApplyValidatingWebhookConfiguration ensures the form of the specified
// validatingwebhookconfiguration is present in the API. If it does not exist,
// it will be created. If it does exist, the metadata of the required
// validatingwebhookconfiguration will be merged with the existing validatingwebhookconfiguration
// and an update performed if the validatingwebhookconfiguration spec and metadata differ from
// the previously required spec and metadata. For further detail, check the top-level comment.
func ApplyValidatingWebhookConfiguration(client dynamic.Interface, recorder events.Recorder,
	requiredOriginal *admissionregistrationv1.ValidatingWebhookConfiguration, expectedGeneration int64) (*admissionregistrationv1.ValidatingWebhookConfiguration, bool, error) {

	// Explcitily specify type for requiredUnstr to get object meta accessor
	requiredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(requiredOriginal)
	if err != nil {
		return nil, false, err
	}

	gvr := admissionregistrationv1.SchemeGroupVersion.WithResource("validatingwebhookconfigurations")
	requiredUnstr := &unstructured.Unstructured{Object: requiredObj}
	resourcedClient := client.Resource(gvr)
	actualUnstr, modified, err := applyUnstructured(resourcedClient, "ValidatingWebhookConfiguration", recorder, requiredUnstr, expectedGeneration)
	if err != nil {
		return nil, modified, err
	}

	actual := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(actualUnstr.Object, actual); err != nil {
		return nil, modified, err
	}
	return actual, modified, nil
}
