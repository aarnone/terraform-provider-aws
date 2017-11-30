package aws

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsKinesisFirehoseDeliveryStreamImportState(
	d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	fmt.Println(d)

	// Extract firehose stream name from arn
	re := regexp.MustCompile(`arn:aws:firehose:[-a-z0-9]+:[0-9]+:deliverystream/(.+)`)
	if !re.MatchString(d.Id()) {
		return nil, fmt.Errorf("[ERROR] The Firehose Delivery Stream ID must be a well formed arn")
	}

	d.Set("name", re.FindStringSubmatch(d.Id())[1])
	return []*schema.ResourceData{d}, nil
}
