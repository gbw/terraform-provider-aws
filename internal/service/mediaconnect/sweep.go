// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build sweep
// +build sweep

package mediaconnect

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediaconnect"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep/awsv2"
)

func init() {
	resource.AddTestSweepers("aws_mediaconnect_flow", &resource.Sweeper{
		Name: "aws_mediaconnect_flow",
		F:    sweepFlows,
	})
}

func sweepFlows(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.MediaConnectClient(ctx)
	sweepResources := make([]sweep.Sweepable, 0)
	in := &mediaconnect.ListFlowsInput{}

	pages := mediaconnect.NewListFlowsPaginator(conn, in)

	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if awsv2.SkipSweepError(err) {
			log.Println("[WARN] Skipping MediaConnect Flows sweep for %s: %s", region, err)
			return nil
		}

		if err != nil {
			return fmt.Errorf("error retrieving MediaConnect Flows: %w", err)
		}

		for _, flow := range page.Flows {
			id := aws.ToString(flow.FlowArn)
			log.Printf("[INFO] Deleting MediaConnect Flows: %s", id)

			r := ResourceFlow()
			d := r.Data(nil)
			d.SetId(id)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}
	}

	if err := sweep.SweepOrchestrator(ctx, sweepResources); err != nil {
		return fmt.Errorf("error sweeping MediaConnect Flows for %s: %w", region, err)
	}

	return nil
}
