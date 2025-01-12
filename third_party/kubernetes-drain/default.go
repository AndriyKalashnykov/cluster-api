/*
Copyright 2019 The Kubernetes Authors.

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

package drain

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// This file contains default implementations of how to
// drain/cordon/uncordon nodes.  These functions may be called
// directly, or their functionality copied into your own code, for
// example if you want different output behaviour.

// RunNodeDrain shows the canonical way to drain a node.
// You should first cordon the node, e.g. using RunCordonOrUncordon
// Deprecated: This method has been deprecated in favor of using the
// same method from the upstream kubectl codebase.
func RunNodeDrain(ctx context.Context, drainer *Helper, nodeName string) error {
	// TODO(justinsb): Ensure we have adequate e2e coverage of this function in library consumers
	list, errs := drainer.GetPodsForDeletion(ctx, nodeName)
	if errs != nil {
		return utilerrors.NewAggregate(errs)
	}
	if warnings := list.Warnings(); warnings != "" {
		fmt.Fprintf(drainer.ErrOut, "WARNING: %s\n", warnings)
	}

	if err := drainer.DeleteOrEvictPods(ctx, list.Pods()); err != nil {
		// Maybe warn about non-deleted pods here
		return err
	}
	return nil
}

// RunCordonOrUncordon demonstrates the canonical way to cordon or uncordon a Node
// Deprecated: This method has been deprecated in favor of using the
// same method from the upstream kubectl codebase.
func RunCordonOrUncordon(ctx context.Context, drainer *Helper, node *corev1.Node, desired bool) error {
	// TODO(justinsb): Ensure we have adequate e2e coverage of this function in library consumers
	c := NewCordonHelper(node)

	if updateRequired := c.UpdateIfRequired(desired); !updateRequired {
		// Already done
		return nil
	}

	err, patchErr := c.PatchOrReplace(ctx, drainer.Client)
	if patchErr != nil {
		return patchErr
	}
	if err != nil {
		return err
	}

	return nil
}
