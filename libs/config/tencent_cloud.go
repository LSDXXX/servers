package config

type TencentCloudConfig struct {
	SecretId      string `yaml:"secret_id"`
	SecretKey     string `yaml:"secret_key"`
	SmsSdkAppId   string `yaml:"sms_sdk_app_id"`
	VmsSdkAppId   string `yaml:"vms_sdk_app_id"`
	Region        string `yaml:"region"`
	SmsTemplateId string `yaml:"sms_template_id"`
	VmsTemplateId string `yaml:"vms_template_id"`
}
