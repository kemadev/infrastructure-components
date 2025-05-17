package priorityclass

import (
	"fmt"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	schedulingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/scheduling/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateDefaultPriorityClasses(ctx *pulumi.Context) error {
	_, err := schedulingv1.NewPriorityClass(
		ctx,
		"default",
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(0),
			Description: pulumi.String(
				"Default priority class, used by pods that do not specify a priority class",
			),
			GlobalDefault: pulumi.Bool(true),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String("default"),
				Labels: pulumi.StringMap{},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create default priority class: %v", err)
	}

	_, err = schedulingv1.NewPriorityClass(
		ctx,
		"default",
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(0),
			Description: pulumi.String(
				"Default priority class, used by pods that do not specify a priority class",
			),
			GlobalDefault: pulumi.Bool(true),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String("default"),
				Labels: pulumi.StringMap{},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create default priority class: %v", err)
	}
	return nil
}

// -2147483648 to 1000000000
