// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package cloudfront

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []func(context.Context) (datasource.DataSourceWithConfigure, error) {
	return []func(context.Context) (datasource.DataSourceWithConfigure, error){}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []func(context.Context) (resource.ResourceWithConfigure, error) {
	return []func(context.Context) (resource.ResourceWithConfigure, error){}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) map[string]func() *schema.Resource {
	return map[string]func() *schema.Resource{
		"aws_cloudfront_cache_policy":                   DataSourceCachePolicy,
		"aws_cloudfront_distribution":                   DataSourceDistribution,
		"aws_cloudfront_function":                       DataSourceFunction,
		"aws_cloudfront_log_delivery_canonical_user_id": DataSourceLogDeliveryCanonicalUserID,
		"aws_cloudfront_origin_access_identities":       DataSourceOriginAccessIdentities,
		"aws_cloudfront_origin_access_identity":         DataSourceOriginAccessIdentity,
		"aws_cloudfront_origin_request_policy":          DataSourceOriginRequestPolicy,
		"aws_cloudfront_realtime_log_config":            DataSourceRealtimeLogConfig,
		"aws_cloudfront_response_headers_policy":        DataSourceResponseHeadersPolicy,
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) map[string]func() *schema.Resource {
	return map[string]func() *schema.Resource{
		"aws_cloudfront_cache_policy":                   ResourceCachePolicy,
		"aws_cloudfront_distribution":                   ResourceDistribution,
		"aws_cloudfront_field_level_encryption_config":  ResourceFieldLevelEncryptionConfig,
		"aws_cloudfront_field_level_encryption_profile": ResourceFieldLevelEncryptionProfile,
		"aws_cloudfront_function":                       ResourceFunction,
		"aws_cloudfront_key_group":                      ResourceKeyGroup,
		"aws_cloudfront_monitoring_subscription":        ResourceMonitoringSubscription,
		"aws_cloudfront_origin_access_control":          ResourceOriginAccessControl,
		"aws_cloudfront_origin_access_identity":         ResourceOriginAccessIdentity,
		"aws_cloudfront_origin_request_policy":          ResourceOriginRequestPolicy,
		"aws_cloudfront_public_key":                     ResourcePublicKey,
		"aws_cloudfront_realtime_log_config":            ResourceRealtimeLogConfig,
		"aws_cloudfront_response_headers_policy":        ResourceResponseHeadersPolicy,
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.CloudFront
}

var ServicePackage = &servicePackage{}
