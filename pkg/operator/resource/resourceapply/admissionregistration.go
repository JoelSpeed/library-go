package resourceapply

import (
	"context"

	"github.com/openshift/library-go/pkg/operator/events"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	gvr := admissionregistrationv1.SchemeGroupVersion.WithResource("mutatingwebhookconfigurations")
	resourcedClient := client.Resource(gvr)

	required := requiredOriginal.DeepCopy()
	if err := copyMutatingWebhookCABundle(resourcedClient, required); err != nil {
		return nil, false, err
	}

	// Explcitily specify type for requiredUnstr to get object meta accessor
	requiredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(required)
	if err != nil {
		return nil, false, err
	}
	requiredUnstr := &unstructured.Unstructured{Object: requiredObj}

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

func copyMutatingWebhookCABundle(resourceClient dynamic.ResourceInterface, required *admissionregistrationv1.MutatingWebhookConfiguration) error {
	existingUnstr, err := resourceClient.Get(context.TODO(), required.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	existing := &admissionregistrationv1.MutatingWebhookConfiguration{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(existingUnstr.Object, existing); err != nil {
		return err
	}

	existingMutatingWebhooks := make(map[string]admissionregistrationv1.MutatingWebhook)
	for _, mutatingWebhook := range existing.Webhooks {
		existingMutatingWebhooks[mutatingWebhook.Name] = mutatingWebhook
	}

	webhooks := make([]admissionregistrationv1.MutatingWebhook, len(required.Webhooks))
	for i, mutatingWebhook := range required.Webhooks {
		if webhook, ok := existingMutatingWebhooks[mutatingWebhook.Name]; ok {
			mutatingWebhook.ClientConfig.CABundle = webhook.ClientConfig.CABundle
		}
		webhooks[i] = mutatingWebhook
	}
	required.Webhooks = webhooks
	return nil
}

// ApplyValidatingWebhookConfiguration ensures the form of the specified
// validatingwebhookconfiguration is present in the API. If it does not exist,
// it will be created. If it does exist, the metadata of the required
// validatingwebhookconfiguration will be merged with the existing validatingwebhookconfiguration
// and an update performed if the validatingwebhookconfiguration spec and metadata differ from
// the previously required spec and metadata. For further detail, check the top-level comment.
func ApplyValidatingWebhookConfiguration(client dynamic.Interface, recorder events.Recorder,
	requiredOriginal *admissionregistrationv1.ValidatingWebhookConfiguration, expectedGeneration int64) (*admissionregistrationv1.ValidatingWebhookConfiguration, bool, error) {

	gvr := admissionregistrationv1.SchemeGroupVersion.WithResource("validatingwebhookconfigurations")
	resourcedClient := client.Resource(gvr)

	required := requiredOriginal.DeepCopy()
	if err := copyValidatingWebhookCABundle(resourcedClient, required); err != nil {
		return nil, false, err
	}

	// Explcitily specify type for requiredUnstr to get object meta accessor
	requiredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(requiredOriginal)
	if err != nil {
		return nil, false, err
	}

	requiredUnstr := &unstructured.Unstructured{Object: requiredObj}

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

func copyValidatingWebhookCABundle(resourceClient dynamic.ResourceInterface, required *admissionregistrationv1.ValidatingWebhookConfiguration) error {
	existingUnstr, err := resourceClient.Get(context.TODO(), required.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	existing := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(existingUnstr.Object, existing); err != nil {
		return err
	}

	existingValidatingWebhooks := make(map[string]admissionregistrationv1.ValidatingWebhook)
	for _, mutatingWebhook := range existing.Webhooks {
		existingValidatingWebhooks[mutatingWebhook.Name] = mutatingWebhook
	}

	webhooks := make([]admissionregistrationv1.ValidatingWebhook, len(required.Webhooks))
	for i, mutatingWebhook := range required.Webhooks {
		if webhook, ok := existingValidatingWebhooks[mutatingWebhook.Name]; ok {
			mutatingWebhook.ClientConfig.CABundle = webhook.ClientConfig.CABundle
		}
		webhooks[i] = mutatingWebhook
	}
	required.Webhooks = webhooks
	return nil
}
