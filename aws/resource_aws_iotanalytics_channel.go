package aws

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotanalytics"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func generateChannelCustomerManagedS3Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func generateChannelServiceManagedS3Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
}

func generateChannelStorageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"customer_managed_s3": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"storage.0.service_managed_s3"},
				Elem:          generateChannelCustomerManagedS3Schema(),
			},
			"service_managed_s3": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"storage.0.customer_managed_s3"},
				Elem:          generateChannelServiceManagedS3Schema(),
			},
		},
	}
}

func resourceAwsIotAnalyticsChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsIotAnalyticsChannelCreate,
		Read:   resourceAwsIotAnalyticsChannelRead,
		Update: resourceAwsIotAnalyticsChannelUpdate,
		Delete: resourceAwsIotAnalyticsChannelDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     generateChannelStorageSchema(),
			},
			"retention_period": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     generateRetentionPeriodSchema(),
			},
		},
	}
}

func parseChannelCustomerManagedS3(rawCustomerManagedS3 map[string]interface{}) *iotanalytics.CustomerManagedChannelS3Storage {
	bucket := rawCustomerManagedS3["bucket"].(string)
	roleArn := rawCustomerManagedS3["role_arn"].(string)
	customerManagedS3 := &iotanalytics.CustomerManagedChannelS3Storage{
		Bucket:  aws.String(bucket),
		RoleArn: aws.String(roleArn),
	}

	if v, ok := rawCustomerManagedS3["key_prefix"]; ok && len(v.(string)) >= 1 {
		customerManagedS3.KeyPrefix = aws.String(v.(string))
	}

	return customerManagedS3
}

func parseChannelServiceManagedS3(rawServiceManagedS3 map[string]interface{}) *iotanalytics.ServiceManagedChannelS3Storage {
	return &iotanalytics.ServiceManagedChannelS3Storage{}
}

func parseChannelStorage(rawChannelStorage map[string]interface{}) *iotanalytics.ChannelStorage {

	var customerManagedS3 *iotanalytics.CustomerManagedChannelS3Storage
	if list := rawChannelStorage["customer_managed_s3"].([]interface{}); len(list) > 0 {
		rawCustomerManagedS3 := list[0].(map[string]interface{})
		customerManagedS3 = parseChannelCustomerManagedS3(rawCustomerManagedS3)
	}

	var serviceManagedS3 *iotanalytics.ServiceManagedChannelS3Storage
	if list := rawChannelStorage["service_managed_s3"].([]interface{}); len(list) > 0 {
		rawServiceManagedS3 := list[0].(map[string]interface{})
		serviceManagedS3 = parseChannelServiceManagedS3(rawServiceManagedS3)
	}

	return &iotanalytics.ChannelStorage{
		CustomerManagedS3: customerManagedS3,
		ServiceManagedS3:  serviceManagedS3,
	}
}

func resourceAwsIotAnalyticsChannelCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotanalyticsconn

	params := &iotanalytics.CreateChannelInput{
		ChannelName: aws.String(d.Get("name").(string)),
	}

	channelStorageSet := d.Get("storage").(*schema.Set).List()
	if len(channelStorageSet) >= 1 {
		rawChannelStorage := channelStorageSet[0].(map[string]interface{})
		params.ChannelStorage = parseChannelStorage(rawChannelStorage)
	}

	retentionPeriodSet := d.Get("retention_period").(*schema.Set).List()
	if len(retentionPeriodSet) >= 1 {
		rawRetentionPeriod := retentionPeriodSet[0].(map[string]interface{})
		params.RetentionPeriod = parseRetentionPeriod(rawRetentionPeriod)
	}

	log.Printf("[DEBUG] Create IoTAnalytics Channel: %s", params)

	retrySecondsList := [6]int{1, 2, 5, 8, 10, 0}

	var err error

	// Primitive retry.
	// During testing channel, problem was detected.
	// When we try to create channel model and role arn that
	// will be assumed by channel during one apply we get:
	// 'Unable to assume role, role ARN' error. However if we run apply
	// second time(when all required resources are created) channel will be created successfully.
	// So we suppose that problem is that AWS return response of successful role arn creation before
	// process of creation is really ended, and then creation of channel model fails.
	for _, sleepSeconds := range retrySecondsList {
		err = nil

		_, err = conn.CreateChannel(params)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(sleepSeconds) * time.Second)
	}

	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return resourceAwsIotAnalyticsChannelRead(d, meta)
}

func flattenChannelCustomerManagedS3(customerManagedS3 *iotanalytics.CustomerManagedChannelS3Storage) map[string]interface{} {
	if customerManagedS3 == nil {
		return nil
	}

	rawCustomerManagedS3 := make(map[string]interface{})

	rawCustomerManagedS3["bucket"] = aws.StringValue(customerManagedS3.Bucket)
	rawCustomerManagedS3["role_arn"] = aws.StringValue(customerManagedS3.RoleArn)

	if customerManagedS3.KeyPrefix != nil {
		rawCustomerManagedS3["key_prefix"] = aws.StringValue(customerManagedS3.KeyPrefix)
	}

	return rawCustomerManagedS3
}

func flattenChannelServiceManagedS3(serviceManagedS3 *iotanalytics.ServiceManagedChannelS3Storage) map[string]interface{} {
	if serviceManagedS3 == nil {
		return nil
	}

	rawServiceManagedS3 := make(map[string]interface{})
	return rawServiceManagedS3
}

func flattenChannelStorage(channelStorage *iotanalytics.ChannelStorage) map[string]interface{} {
	customerManagedS3 := flattenChannelCustomerManagedS3(channelStorage.CustomerManagedS3)
	serviceManagedS3 := flattenChannelServiceManagedS3(channelStorage.ServiceManagedS3)

	if customerManagedS3 == nil && serviceManagedS3 == nil {
		return nil
	}

	rawStorage := make(map[string]interface{})
	rawStorage["customer_managed_s3"] = wrapMapInList(customerManagedS3)
	rawStorage["service_managed_s3"] = wrapMapInList(serviceManagedS3)
	return rawStorage
}

func resourceAwsIotAnalyticsChannelRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotanalyticsconn

	params := &iotanalytics.DescribeChannelInput{
		ChannelName: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Reading IoTAnalytics Channel: %s", params)

	out, err := conn.DescribeChannel(params)

	if err != nil {
		return err
	}

	d.Set("name", out.Channel.Name)
	storage := flattenChannelStorage(out.Channel.Storage)
	d.Set("storage", wrapMapInList(storage))
	retentionPeriod := flattenRetentionPeriod(out.Channel.RetentionPeriod)
	d.Set("retention_period", wrapMapInList(retentionPeriod))

	return nil
}

func resourceAwsIotAnalyticsChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotanalyticsconn

	params := &iotanalytics.UpdateChannelInput{
		ChannelName: aws.String(d.Get("name").(string)),
	}

	channelStorageSet := d.Get("storage").(*schema.Set).List()
	if len(channelStorageSet) >= 1 {
		rawChannelStorage := channelStorageSet[0].(map[string]interface{})
		params.ChannelStorage = parseChannelStorage(rawChannelStorage)
	}

	retentionPeriodSet := d.Get("retention_period").(*schema.Set).List()
	if len(retentionPeriodSet) >= 1 {
		rawRetentionPeriod := retentionPeriodSet[0].(map[string]interface{})
		params.RetentionPeriod = parseRetentionPeriod(rawRetentionPeriod)
	}

	log.Printf("[DEBUG] Updating IoTAnalytics Channel: %s", params)

	retrySecondsList := [6]int{1, 2, 5, 8, 10, 0}

	var err error

	// Primitive retry.
	// Full explanation can be found in function `resourceAwsIotAnalyticsChannelCreate`.
	// We suppose that such error can appear during update also, if you update
	// role arn.
	for _, sleepSeconds := range retrySecondsList {
		err = nil

		_, err = conn.UpdateChannel(params)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(sleepSeconds) * time.Second)
	}

	if err != nil {
		return err
	}

	return resourceAwsIotAnalyticsChannelRead(d, meta)
}

func resourceAwsIotAnalyticsChannelDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotanalyticsconn

	params := &iotanalytics.DeleteChannelInput{
		ChannelName: aws.String(d.Id()),
	}
	log.Printf("[DEBUG] Delete IoTAnalytics Channel: %s", params)
	_, err := conn.DeleteChannel(params)

	return err
}
