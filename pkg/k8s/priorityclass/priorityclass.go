package priorityclass

import (
	"fmt"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	schedulingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/scheduling/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	// Very low priority class, lower than low, does not preempt other pods
	PriorityClassEventual = "eventual"
	// Low priority class, lower than default, higher than eventual, preempts other pods
	PriorityClassLow = "low"
	// Default priority class, lower than normal, higher than low, preempts other pods, used by pods that do not specify a priority class, should be used explicitly
	PriorityClassDefault = "default"
	// Normal priority, lower than moderate, higher than default, preempts other pods, should be used as a default
	PriorityClassNormal = "normal"
	// Moderate priority, lower than high, higher than normal, preempts other pods
	PriorityClassModerate = "moderate"
	// High priority, higher than moderate, preempts other pods
	PriorityClassHigh = "high"
)

func CreateDefaultPriorityClasses(ctx *pulumi.Context) error {
	_, err := schedulingv1.NewPriorityClass(
		ctx,
		PriorityClassEventual,
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(-1000000),
			Description: pulumi.String(
				"Very low priority class, lower than low, does not preempt other pods",
			),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String(PriorityClassEventual),
				Labels: pulumi.StringMap{},
			},
			PreemptionPolicy: pulumi.String("Never"),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create "+PriorityClassEventual+" priority class: %v", err)
	}

	_, err = schedulingv1.NewPriorityClass(
		ctx,
		PriorityClassLow,
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(-1000),
			Description: pulumi.String(
				"Low priority class, lower than default, higher than eventual, preempts other pods",
			),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String(PriorityClassLow),
				Labels: pulumi.StringMap{},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create "+PriorityClassLow+" priority class: %v", err)
	}

	_, err = schedulingv1.NewPriorityClass(
		ctx,
		PriorityClassDefault,
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(0),
			Description: pulumi.String(
				"Default priority class, lower than normal, higher than low, preempts other pods, used by pods that do not specify a priority class, should be used explicitly",
			),
			GlobalDefault: pulumi.Bool(true),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String(PriorityClassDefault),
				Labels: pulumi.StringMap{},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create "+PriorityClassDefault+" priority class: %v", err)
	}

	_, err = schedulingv1.NewPriorityClass(
		ctx,
		PriorityClassNormal,
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(1000),
			Description: pulumi.String(
				"Normal priority, lower than moderate, higher than default, preempts other pods, should be used as a default",
			),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String(PriorityClassNormal),
				Labels: pulumi.StringMap{},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create "+PriorityClassNormal+" priority class: %v", err)
	}

	_, err = schedulingv1.NewPriorityClass(
		ctx,
		PriorityClassModerate,
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(500000),
			Description: pulumi.String(
				"Moderate priority, lower than high, higher than normal, preempts other pods",
			),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String(PriorityClassModerate),
				Labels: pulumi.StringMap{},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create "+PriorityClassModerate+" priority class: %v", err)
	}

	_, err = schedulingv1.NewPriorityClass(
		ctx,
		PriorityClassHigh,
		&schedulingv1.PriorityClassArgs{
			Value: pulumi.Int(1000000),
			Description: pulumi.String(
				"High priority, higher than moderate, preempts other pods",
			),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String(PriorityClassHigh),
				Labels: pulumi.StringMap{},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create "+PriorityClassHigh+" priority class: %v", err)
	}

	return nil
}
