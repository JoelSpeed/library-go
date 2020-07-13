package resourceapply

import (
	"context"

	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog"
)

// applyUnstructured ensures the form of the specified usntructured is present in the API.
// If it does not exist, it will be created. If it does exist, the metadata of the required
// usntructured will be merged with the existing usntructured and an update performed if the
// usntructured spec and metadata differ from the previously required spec and metadata.
// For further detail, check the top-level comment.
func applyUnstructured(resourceClient dynamic.ResourceInterface, kind string, recorder events.Recorder,
	requiredOriginal *unstructured.Unstructured, expectedGeneration int64) (*unstructured.Unstructured, bool, error) {

	required := requiredOriginal.DeepCopy()

	if err := SetSpecHashAnnotation(required, required.Object["spec"]); err != nil {
		return nil, false, err
	}

	existing, err := resourceClient.Get(context.TODO(), required.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := resourceClient.Create(context.TODO(), required, metav1.CreateOptions{})
		reportCreateEvent(recorder, required, err)
		return actual, true, err
	}

	existingCopy := existing.DeepCopy()
	modified, err := ensureObjectMeta(existingCopy, required)
	if err != nil {
		return nil, false, err
	}

	if !modified && existingCopy.GetGeneration() == expectedGeneration {
		return existingCopy, false, nil
	}

	requiredSpec, exists, err := unstructured.NestedMap(required.Object, "spec")
	if err != nil {
		return nil, false, err
	}
	if !exists {
		// No spec to update
		return existingCopy, false, nil
	}

	// at this point we know that we're going to perform a write.  We're just trying to get the object correct
	toWrite := existingCopy // shallow copy so the code reads easier
	if err := unstructured.SetNestedField(toWrite.Object, requiredSpec, "spec"); err != nil {
		return nil, false, err
	}

	if klog.V(4) {
		klog.Infof("%s %q changes: %v", kind, required.GetNamespace()+"/"+required.GetName(), JSONPatchNoError(existing, toWrite))
	}

	actual, err := resourceClient.Update(context.TODO(), toWrite, metav1.UpdateOptions{})
	reportUpdateEvent(recorder, required, err)
	return actual, true, err
}

// This is based on resourcemerge.EnsureObjectMeta but uses the metav1.Object interface instead
// TODO: Update this to use resourcemerge.EnsureObjectMeta or update resourcemerge.EnsureObjectMeta to use the interface
func ensureObjectMeta(existing, required metav1.Object) (bool, error) {
	modified := resourcemerge.BoolPtr(false)

	namespace := existing.GetNamespace()
	name := existing.GetName()
	labels := existing.GetLabels()
	annotations := existing.GetAnnotations()

	resourcemerge.SetStringIfSet(modified, &namespace, required.GetNamespace())
	resourcemerge.SetStringIfSet(modified, &name, required.GetName())
	resourcemerge.MergeMap(modified, &labels, required.GetLabels())
	resourcemerge.MergeMap(modified, &annotations, required.GetAnnotations())

	existing.SetNamespace(namespace)
	existing.SetName(name)
	existing.SetLabels(labels)
	existing.SetAnnotations(annotations)

	return *modified, nil
}
