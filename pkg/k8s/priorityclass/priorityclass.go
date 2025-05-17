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
		"eventual",
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(-1000000),
			Description: pulumi.String(
				"Low priority class, lower than low, default, normal or high, do not preempts other pods",
			),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String("eventual"),
				Labels: pulumi.StringMap{},
			},
			PreemptionPolicy: pulumi.String("Never"),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create eventual priority class: %v", err)
	}

	_, err = schedulingv1.NewPriorityClass(
		ctx,
		"low",
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(-1000),
			Description: pulumi.String(
				"Low priority class, lower than default, normal or high, preempts other pods",
			),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String("low"),
				Labels: pulumi.StringMap{},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create low priority class: %v", err)
	}

	_, err = schedulingv1.NewPriorityClass(
		ctx,
		"default",
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(0),
			Description: pulumi.String(
				"Default priority class, lower than normal or high, to be used by pods that do not specify a priority class",
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
		"normal",
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(1000),
			Description: pulumi.String(
				"Normal priority, higher than default, preempts other pods",
			),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String("normal"),
				Labels: pulumi.StringMap{},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create normal priority class: %v", err)
	}

	_, err = schedulingv1.NewPriorityClass(
		ctx,
		"high",
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(1000000),
			Description: pulumi.String(
				"High priority, higher normal, preempts other pods",
			),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String("high"),
				Labels: pulumi.StringMap{},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create high priority class: %v", err)
	}

	return nil
}
