package aws

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

func resourceAwsKinesisFirehoseMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AWS Kinesis Firehose Delivery Stream State v0; migrating to v2")
		newIs, err := migrateKinesisFirehoseV0toV1(is)
		if err != nil {
			return nil, err
		}

		return migrateKinesisFirehoseV1toV2(newIs)
	case 1:
		log.Println("[INFO] Found AWS Kinesis Firehose Delivery Stream State v1; migrating to v2")
		return migrateKinesisFirehoseV1toV2(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateKinesisFirehoseV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty Kinesis Firehose Delivery State; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] Attributes before migration from V0 to V1: %#v", is.Attributes)

	// migrate flate S3 configuration to a s3_configuration block
	// grab initial values
	is.Attributes["s3_configuration.#"] = "1"
	// Required parameters
	is.Attributes["s3_configuration.0.role_arn"] = is.Attributes["role_arn"]
	is.Attributes["s3_configuration.0.bucket_arn"] = is.Attributes["s3_bucket_arn"]

	// Optional parameters
	if is.Attributes["s3_buffer_size"] != "" {
		is.Attributes["s3_configuration.0.buffer_size"] = is.Attributes["s3_buffer_size"]
	}
	if is.Attributes["s3_data_compression"] != "" {
		is.Attributes["s3_configuration.0.compression_format"] = is.Attributes["s3_data_compression"]
	}
	if is.Attributes["s3_buffer_interval"] != "" {
		is.Attributes["s3_configuration.0.buffer_interval"] = is.Attributes["s3_buffer_interval"]
	}
	if is.Attributes["s3_prefix"] != "" {
		is.Attributes["s3_configuration.0.prefix"] = is.Attributes["s3_prefix"]
	}

	delete(is.Attributes, "role_arn")
	delete(is.Attributes, "s3_bucket_arn")
	delete(is.Attributes, "s3_buffer_size")
	delete(is.Attributes, "s3_data_compression")
	delete(is.Attributes, "s3_buffer_interval")
	delete(is.Attributes, "s3_prefix")

	log.Printf("[DEBUG] Attributes after migration from V0 to V1: %#v", is.Attributes)
	return is, nil
}

func migrateKinesisFirehoseV1toV2(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty Kinesis Firehose Delivery State; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] Attributes before migration from V1 to V2: %#v", is.Attributes)

	// basic s3_configuration will be deprecated, changes are now reflected in extended_s3_configuration
	// keep everything under extended_s3_configuration to avoid confusion

	for key, value := range is.Attributes {
		if strings.HasPrefix(key, "s3_configuration") {
			newKey := strings.Replace(key, "s3_configuration", "extended_s3_configuration", 1)
			is.Attributes[newKey] = value
			delete(is.Attributes, key)
		}
	}

	// processor parameters are now stored individually
	// convert parameter_name/parameter_value pairs
	paramNameRe := regexp.MustCompile(`^((extended_s3_configuration\.\d+\.processing_configuration\.\d+\.processors\.\d+\.parameters\.)\d+\.)parameter_name$`)
	for key, value := range is.Attributes {
		if parsed := paramNameRe.FindStringSubmatch(key); len(parsed) == 3 {
			prefix, prefixWithIndex, parameterName := parsed[2], parsed[1], value

			var newKey string
			switch parameterName {
			case "LambdaArn":
				newKey = prefix + "lambda_arn"
			case "NumberOfRetries":
				newKey = prefix + "number_of_retries"
			default:
				return nil, fmt.Errorf("[ERROR] Unexpected parameter name for processor: %v", parameterName)
			}

			valueKey := prefixWithIndex + "parameter_value"
			is.Attributes[newKey] = is.Attributes[valueKey]
			delete(is.Attributes, key)
			delete(is.Attributes, valueKey)
		}
	}

	// cleanup orphan values and list counter
	paramValueRe := regexp.MustCompile(`^extended_s3_configuration\.\d+\.processing_configuration\.\d+\.processors\.\d+\.parameters\.\d+\.parameter_value$`)
	counterRe := regexp.MustCompile(`^extended_s3_configuration\.\d+\.processing_configuration\.\d+\.processors\.\d+\.parameters\.#$`)
	for key := range is.Attributes {
		if counterRe.MatchString(key) ||
			paramValueRe.MatchString(key) {
			delete(is.Attributes, key)
		}
	}

	log.Printf("[DEBUG] Attributes after migration from V1 to V2: %#v", is.Attributes)
	return is, nil
}
