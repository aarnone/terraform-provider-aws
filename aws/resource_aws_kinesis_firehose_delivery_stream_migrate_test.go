package aws

import (
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestAWSKinesisFirehoseMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		Attributes   map[string]string
		Expected     map[string]string
		Meta         interface{}
	}{
		"v0 to v2": {
			StateVersion: 0,
			Attributes: map[string]string{
				// EBS
				"role_arn":            "arn:aws:iam::somenumber:role/tf_acctest_4271506651559170635",
				"s3_bucket_arn":       "arn:aws:s3:::tf-test-bucket",
				"s3_buffer_interval":  "400",
				"s3_buffer_size":      "10",
				"s3_data_compression": "GZIP",
			},
			Expected: map[string]string{
				"extended_s3_configuration.#":                    "1",
				"extended_s3_configuration.0.bucket_arn":         "arn:aws:s3:::tf-test-bucket",
				"extended_s3_configuration.0.buffer_interval":    "400",
				"extended_s3_configuration.0.buffer_size":        "10",
				"extended_s3_configuration.0.compression_format": "GZIP",
				"extended_s3_configuration.0.role_arn":           "arn:aws:iam::somenumber:role/tf_acctest_4271506651559170635",
			},
		},
		"v0 to v2, partial": {
			StateVersion: 0,
			Attributes: map[string]string{
				// EBS
				"role_arn":      "arn:aws:iam::somenumber:role/tf_acctest_4271506651559170635",
				"s3_bucket_arn": "arn:aws:s3:::tf-test-bucket",
			},
			Expected: map[string]string{
				"extended_s3_configuration.#":            "1",
				"extended_s3_configuration.0.bucket_arn": "arn:aws:s3:::tf-test-bucket",
				"extended_s3_configuration.0.role_arn":   "arn:aws:iam::somenumber:role/tf_acctest_4271506651559170635",
			},
		},
		"v1 to v2": {
			StateVersion: 1,
			Attributes: map[string]string{
				"destination":                                                              "s3",
				"s3_configuration.#":                                                       "1",
				"s3_configuration.0.bucket_arn":                                            "arn:aws:s3:::tf-test-bucket-7748474670850130673",
				"s3_configuration.0.buffer_interval":                                       "300",
				"s3_configuration.0.buffer_size":                                           "5",
				"s3_configuration.0.cloudwatch_logging_options.#":                          "1",
				"s3_configuration.0.cloudwatch_logging_options.2513562885.enabled":         "true",
				"s3_configuration.0.cloudwatch_logging_options.2513562885.log_group_name":  "groupname",
				"s3_configuration.0.cloudwatch_logging_options.2513562885.log_stream_name": "streamname",
				"s3_configuration.0.compression_format":                                    "UNCOMPRESSED",
				"s3_configuration.0.kms_key_arn":                                           "kmskeyarn",
				"s3_configuration.0.prefix":                                                "s3prefix",
				"s3_configuration.0.role_arn":                                              "arn:aws:iam::666:role/tf_acctest_firehose_delivery_role_7748474670850130673",
			},
			Expected: map[string]string{
				"destination":                                                                       "s3",
				"extended_s3_configuration.#":                                                       "1",
				"extended_s3_configuration.0.bucket_arn":                                            "arn:aws:s3:::tf-test-bucket-7748474670850130673",
				"extended_s3_configuration.0.buffer_interval":                                       "300",
				"extended_s3_configuration.0.buffer_size":                                           "5",
				"extended_s3_configuration.0.cloudwatch_logging_options.#":                          "1",
				"extended_s3_configuration.0.cloudwatch_logging_options.2513562885.enabled":         "true",
				"extended_s3_configuration.0.cloudwatch_logging_options.2513562885.log_group_name":  "groupname",
				"extended_s3_configuration.0.cloudwatch_logging_options.2513562885.log_stream_name": "streamname",
				"extended_s3_configuration.0.compression_format":                                    "UNCOMPRESSED",
				"extended_s3_configuration.0.kms_key_arn":                                           "kmskeyarn",
				"extended_s3_configuration.0.prefix":                                                "s3prefix",
				"extended_s3_configuration.0.role_arn":                                              "arn:aws:iam::666:role/tf_acctest_firehose_delivery_role_7748474670850130673",
			},
		},
		"v1 to v2, processors": {
			StateVersion: 1,
			Attributes: map[string]string{
				"extended_s3_configuration.#":                                                                      "1",
				"extended_s3_configuration.0.bucket_arn":                                                           "arn:aws:s3:::tf-test-bucket-4747185083350722827",
				"extended_s3_configuration.0.buffer_interval":                                                      "300",
				"extended_s3_configuration.0.buffer_size":                                                          "5",
				"extended_s3_configuration.0.compression_format":                                                   "UNCOMPRESSED",
				"extended_s3_configuration.0.kms_key_arn":                                                          "kmskeyarn",
				"extended_s3_configuration.0.prefix":                                                               "s3prefix",
				"extended_s3_configuration.0.processing_configuration.#":                                           "1",
				"extended_s3_configuration.0.processing_configuration.0.enabled":                                   "false",
				"extended_s3_configuration.0.processing_configuration.0.processors.#":                              "1",
				"extended_s3_configuration.0.processing_configuration.0.processors.0.parameters.#":                 "2",
				"extended_s3_configuration.0.processing_configuration.0.processors.0.parameters.0.parameter_name":  "LambdaArn",
				"extended_s3_configuration.0.processing_configuration.0.processors.0.parameters.0.parameter_value": "arn:aws:lambda:us-west-2:666:function:aws_kinesis_firehose_delivery_stream_test_rvab6",
				"extended_s3_configuration.0.processing_configuration.0.processors.0.parameters.1.parameter_name":  "NumberOfRetries",
				"extended_s3_configuration.0.processing_configuration.0.processors.0.parameters.1.parameter_value": "3",
				"extended_s3_configuration.0.processing_configuration.0.processors.0.type":                         "Lambda",
				"extended_s3_configuration.0.role_arn":                                                             "arn:aws:iam::666:role/tf_acctest_firehose_delivery_role_4747185083350722827",
			},
			Expected: map[string]string{
				"extended_s3_configuration.#":                                                                      "1",
				"extended_s3_configuration.0.bucket_arn":                                                           "arn:aws:s3:::tf-test-bucket-4747185083350722827",
				"extended_s3_configuration.0.buffer_interval":                                                      "300",
				"extended_s3_configuration.0.buffer_size":                                                          "5",
				"extended_s3_configuration.0.compression_format":                                                   "UNCOMPRESSED",
				"extended_s3_configuration.0.kms_key_arn":                                                          "kmskeyarn",
				"extended_s3_configuration.0.prefix":                                                               "s3prefix",
				"extended_s3_configuration.0.processing_configuration.#":                                           "1",
				"extended_s3_configuration.0.processing_configuration.0.enabled":                                   "false",
				"extended_s3_configuration.0.processing_configuration.0.processors.#":                              "1",
				"extended_s3_configuration.0.processing_configuration.0.processors.0.parameters.lambda_arn":        "arn:aws:lambda:us-west-2:666:function:aws_kinesis_firehose_delivery_stream_test_rvab6",
				"extended_s3_configuration.0.processing_configuration.0.processors.0.parameters.number_of_retries": "3",
				"extended_s3_configuration.0.processing_configuration.0.processors.0.type":                         "Lambda",
				"extended_s3_configuration.0.role_arn":                                                             "arn:aws:iam::666:role/tf_acctest_firehose_delivery_role_4747185083350722827",
			},
		},
		"v1 to v2, orphans": {
			StateVersion: 1,
			Attributes: map[string]string{
				"extended_s3_configuration.0.processing_configuration.#":                                           "5",
				"extended_s3_configuration.0.processing_configuration.1.processors.0.parameters.0.parameter_value": "arn:aws:lambda:us-west-2:666:function:aws_kinesis_firehose_delivery_stream_test_rvab6",
				"extended_s3_configuration.0.processing_configuration.2.processors.0.parameters.1.parameter_value": "3",
			},
			Expected: map[string]string{},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         "i-abc123",
			Attributes: tc.Attributes,
		}
		is, err := resourceAwsKinesisFirehoseMigrateState(
			tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.Expected {
			if is.Attributes[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					tn, k, v, k, is.Attributes[k], is.Attributes)
			}
		}
	}
}

func TestAWSKinesisFirehoseMigrateState_empty(t *testing.T) {
	var is *terraform.InstanceState
	var meta interface{}

	// should handle nil
	is, err := resourceAwsKinesisFirehoseMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
	if is != nil {
		t.Fatalf("expected nil instancestate, got: %#v", is)
	}

	// should handle non-nil but empty
	is = &terraform.InstanceState{}
	is, err = resourceAwsInstanceMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
}
