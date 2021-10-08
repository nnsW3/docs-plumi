package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"

	"github.com/pulumi/docs/tools/resourcedocsgen/pkg"
)

var categoryNameMap = map[string]pkg.PackageCategory{
	"cloud":          pkg.PackageCategoryCloud,
	"database":       pkg.PackageCategoryDatabase,
	"infrastructure": pkg.PackageCategoryInfrastructure,
	"monitoring":     pkg.PackageCategoryMonitoring,
	"network":        pkg.PackageCategoryNetwork,
	"utility":        pkg.PackageCategoryUtility,
	"vcs":            pkg.PackageCategoryVCS,
}

var categoryLookup = map[string]pkg.PackageCategory{
	"aiven":                               pkg.PackageCategoryInfrastructure,
	"akamai":                              pkg.PackageCategoryNetwork,
	"alicloud":                            pkg.PackageCategoryCloud,
	"auth0":                               pkg.PackageCategoryInfrastructure,
	"aws":                                 pkg.PackageCategoryCloud,
	"aws-api-gateway":                     pkg.PackageCategoryCloud,
	"aws-miniflux":                        pkg.PackageCategoryCloud,
	"aws-native":                          pkg.PackageCategoryCloud,
	"aws-quickstart-aurora-mysql":         pkg.PackageCategoryCloud,
	"aws-quickstart-aurora-postgresql":    pkg.PackageCategoryCloud,
	"aws-quickstart-redshift":             pkg.PackageCategoryCloud,
	"aws-serverless":                      pkg.PackageCategoryCloud,
	"aws-quickstart-vpc":                  pkg.PackageCategoryCloud,
	"aws-s3-replicated-bucket":            pkg.PackageCategoryCloud,
	"azure":                               pkg.PackageCategoryCloud,
	"azure-native":                        pkg.PackageCategoryCloud,
	"azure-quickstart-acr-geo-replicated": pkg.PackageCategoryCloud,
	"azure-quickstart-aks":                pkg.PackageCategoryCloud,
	"azure-quickstart-compute":            pkg.PackageCategoryCloud,
	"azure-quickstart-sql":                pkg.PackageCategoryCloud,
	"azuread":                             pkg.PackageCategoryInfrastructure,
	"azuredevops":                         pkg.PackageCategoryInfrastructure,
	"civo":                                pkg.PackageCategoryCloud,
	"cloudamqp":                           pkg.PackageCategoryCloud,
	"cloudflare":                          pkg.PackageCategoryNetwork,
	"cloudinit":                           pkg.PackageCategoryUtility,
	"confluent":                           pkg.PackageCategoryInfrastructure,
	"consul":                              pkg.PackageCategoryInfrastructure,
	"coredns-helm":                        pkg.PackageCategoryNetwork,
	"datadog":                             pkg.PackageCategoryMonitoring,
	"digitalocean":                        pkg.PackageCategoryCloud,
	"dnsimple":                            pkg.PackageCategoryNetwork,
	"docker":                              pkg.PackageCategoryInfrastructure,
	"docker-buildkit":                     pkg.PackageCategoryInfrastructure,
	"eks":                                 pkg.PackageCategoryCloud,
	"equinix-metal":                       pkg.PackageCategoryCloud,
	"f5bigip":                             pkg.PackageCategoryNetwork,
	"fastly":                              pkg.PackageCategoryCloud,
	"gcp":                                 pkg.PackageCategoryCloud,
	"gcp-cloudrun-multi-region":           pkg.PackageCategoryCloud,
	"gcp-project-scaffold":                pkg.PackageCategoryCloud,
	"google-native":                       pkg.PackageCategoryCloud,
	"github":                              pkg.PackageCategoryVCS,
	"github-serverless-webhook":           pkg.PackageCategoryVCS,
	"gitlab":                              pkg.PackageCategoryVCS,
	"hcloud":                              pkg.PackageCategoryCloud,
	"istio-helm":                          pkg.PackageCategoryInfrastructure,
	"jaeger-helm":                         pkg.PackageCategoryMonitoring,
	"kafka":                               pkg.PackageCategoryInfrastructure,
	"keycloak":                            pkg.PackageCategoryInfrastructure,
	"kong":                                pkg.PackageCategoryInfrastructure,
	"kubernetes":                          pkg.PackageCategoryCloud,
	"libvirt":                             pkg.PackageCategoryUtility,
	"linode":                              pkg.PackageCategoryCloud,
	"mailgun":                             pkg.PackageCategoryInfrastructure,
	"minio":                               pkg.PackageCategoryInfrastructure,
	"mongodbatlas":                        pkg.PackageCategoryDatabase,
	"mysql":                               pkg.PackageCategoryDatabase,
	"newrelic":                            pkg.PackageCategoryMonitoring,
	"nginx-ingress-controller-helm":       pkg.PackageCategoryNetwork,
	"nomad":                               pkg.PackageCategoryInfrastructure,
	"ns1":                                 pkg.PackageCategoryNetwork,
	"okta":                                pkg.PackageCategoryInfrastructure,
	"openstack":                           pkg.PackageCategoryCloud,
	"opsgenie":                            pkg.PackageCategoryInfrastructure,
	"pagerduty":                           pkg.PackageCategoryInfrastructure,
	"postgresql":                          pkg.PackageCategoryDatabase,
	"prometheus-helm":                     pkg.PackageCategoryMonitoring,
	"rabbitmq":                            pkg.PackageCategoryInfrastructure,
	"rancher2":                            pkg.PackageCategoryInfrastructure,
	"random":                              pkg.PackageCategoryUtility,
	"rke":                                 pkg.PackageCategoryInfrastructure,
	"run-my-darn-container":               pkg.PackageCategoryInfrastructure,
	"shipa":                               pkg.PackageCategoryCloud,
	"signalfx":                            pkg.PackageCategoryMonitoring,
	"snowflake":                           pkg.PackageCategoryInfrastructure,
	"splunk":                              pkg.PackageCategoryInfrastructure,
	"spotinst":                            pkg.PackageCategoryInfrastructure,
	"sumologic":                           pkg.PackageCategoryMonitoring,
	"tls":                                 pkg.PackageCategoryUtility,
	"vault":                               pkg.PackageCategoryInfrastructure,
	"venafi":                              pkg.PackageCategoryInfrastructure,
	"vsphere":                             pkg.PackageCategoryCloud,
	"wavefront":                           pkg.PackageCategoryMonitoring,
	"yandex":                              pkg.PackageCategoryCloud,
}

// TODO[pulumi/pulumi#7813]: Remove this lookup once display name is available in
// the Pulumi schema.
//
// NOTE: For the time being this lookup map and the one used by the docs
// generator in `pulumi/pulumi` must be kept up-to-date.
//
// titleLookup is a map pf package name to the desired display name
// for display in the TOC menu under API Reference.
var titleLookup = map[string]string{
	"aiven":                               "Aiven",
	"akamai":                              "Akamai",
	"alicloud":                            "Alibaba Cloud",
	"auth0":                               "Auth0",
	"aws":                                 "AWS Classic",
	"aws-api-gateway":                     "AWS API Gateway",
	"aws-miniflux":                        "Miniflux",
	"aws-native":                          "AWS Native",
	"aws-quickstart-aurora-mysql":         "AWS QuickStart Aurora MySQL",
	"aws-quickstart-aurora-postgresql":    "AWS QuickStart Aurora PostgreSQL",
	"aws-quickstart-redshift":             "AWS QuickStart Redshift",
	"aws-serverless":                      "AWS Serverless",
	"aws-quickstart-vpc":                  "AWS QuickStart VPC",
	"aws-s3-replicated-bucket":            "AWS S3 Replicated Bucket",
	"azure":                               "Azure Classic",
	"azure-native":                        "Azure Native",
	"azure-quickstart-acr-geo-replicated": "Azure QuickStart ACR Geo Replicated",
	"azure-quickstart-aks":                "Azure QuickStart AKS",
	"azure-quickstart-compute":            "Azure QuickStart Compute",
	"azure-quickstart-sql":                "Azure QuickStart SQL",
	"azuread":                             "Azure Active Directory",
	"azuredevops":                         "Azure DevOps",
	"azuresel":                            "Azure",
	"civo":                                "Civo",
	"cloudamqp":                           "CloudAMQP",
	"cloudflare":                          "Cloudflare",
	"cloudinit":                           "cloud-init",
	"confluent":                           "Confluent Cloud",
	"consul":                              "Consul",
	"coredns-helm":                        "CoreDNS (Helm)",
	"datadog":                             "Datadog",
	"digitalocean":                        "DigitalOcean",
	"dnsimple":                            "DNSimple",
	"docker":                              "Docker",
	"docker-buildkit":                     "Docker BuildKit",
	"eks":                                 "Amazon EKS",
	"equinix-metal":                       "Equinix Metal",
	"f5bigip":                             "f5 BIG-IP",
	"fastly":                              "Fastly",
	"gcp":                                 "Google Cloud Classic",
	"gcp-cloudrun-multi-region":           "Google Cloud Run Multi-Region",
	"gcp-project-scaffold":                "Google Project Scaffolding",
	"google-native":                       "Google Cloud Native",
	"github":                              "GitHub",
	"github-serverless-webhook":           "GitHub Serverless Webhook",
	"gitlab":                              "GitLab",
	"hcloud":                              "Hetzner Cloud",
	"istio-helm":                          "Istio (Helm)",
	"jaeger-helm":                         "Jaeger (Helm)",
	"kafka":                               "Kafka",
	"keycloak":                            "Keycloak",
	"kong":                                "Kong",
	"kubernetes":                          "Kubernetes",
	"libvirt":                             "libvirt",
	"linode":                              "Linode",
	"mailgun":                             "Mailgun",
	"minio":                               "MinIO",
	"mongodbatlas":                        "MongoDB Atlas",
	"mysql":                               "MySQL",
	"newrelic":                            "New Relic",
	"nginx-ingress-controller-helm":       "NGINX Ingress Controller (Helm)",
	"nomad":                               "Nomad",
	"ns1":                                 "NS1",
	"okta":                                "Okta",
	"openstack":                           "OpenStack",
	"opsgenie":                            "Opsgenie",
	"packet":                              "Packet",
	"pagerduty":                           "PagerDuty",
	"postgresql":                          "PostgreSQL",
	"prometheus-helm":                     "Prometheus (Helm)",
	"rabbitmq":                            "RabbitMQ",
	"rancher2":                            "Rancher 2",
	"random":                              "random",
	"rke":                                 "Rancher RKE",
	"run-my-darn-container":               "Run My Darn Container",
	"shipa":                               "Shipa",
	"signalfx":                            "SignalFx",
	"snowflake":                           "Snowflake",
	"splunk":                              "Splunk",
	"spotinst":                            "Spotinst",
	"sumologic":                           "Sumo Logic",
	"tls":                                 "TLS",
	"vault":                               "Vault",
	"venafi":                              "Venafi",
	"vsphere":                             "vSphere",
	"wavefront":                           "Wavefront",
	"yandex":                              "Yandex",
}

func packageMetadataCmd() *cobra.Command {
	var metadataOutDir string
	var categoryStr string
	var component bool
	var featured bool
	var publisher string
	var title string
	var updatedOn int64

	cmd := &cobra.Command{
		Use:   "metadata <metadataOutDir> [featured]",
		Short: "Generate package metadata from Pulumi schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			status := pkg.PackageStatusGA
			if strings.HasPrefix(version, "v0.") {
				status = pkg.PackageStatusPublicPreview
			}

			var category pkg.PackageCategory
			// If a category was passed-in, use that as the override.
			if categoryStr != "" {
				if n, ok := categoryNameMap[categoryStr]; !ok {
					return errors.New(fmt.Sprintf("invalid category name %s", categoryStr))
				} else {
					category = n
				}
			} else if c, ok := categoryLookup[mainSpec.Name]; ok {
				// Otherwise, try to lookup the package in our existing category
				// lookup.
				category = c
			} else {
				// The default is to categorize the package as "Cloud".
				category = pkg.PackageCategoryCloud
			}

			if title == "" {
				title = mainSpec.Name
			}
			if v, ok := titleLookup[mainSpec.Name]; ok {
				title = v
			}

			// TODO[pulumi/pulumi#7813]: Use tags in the schema spec to determine
			// if the package is a native provider or not.
			native := mainSpec.Attribution == ""
			// Due to the way we currently detect if a package is a native provider
			// or not (see above) a package cannot be both a native provider AND a
			// component so if a package is a component, override the native property
			// to false.
			if component {
				native = false
			}
			pm := pkg.PackageMeta{
				Name:          mainSpec.Name,
				UpdatedOn:     updatedOn,
				Publisher:     publisher,
				Title:         title,
				Description:   mainSpec.Description,
				Category:      category,
				PackageStatus: status,
				Component:     component,
				Featured:      featured,
				Native:        native,
				Version:       version,
				LogoURL:       mainSpec.LogoURL,
			}
			b, err := yaml.Marshal(pm)
			if err != nil {
				return errors.Wrap(err, "generating package metadata")
			}

			metadataFileName := fmt.Sprintf("%s.yaml", mainSpec.Name)
			if err := emitFile(metadataOutDir, metadataFileName, b); err != nil {
				return errors.Wrap(err, "writing metadata file")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&metadataOutDir, "metadataOutDir", "", "The directory path to where the docs will be written to")
	cmd.Flags().StringVar(&categoryStr, "category", "", fmt.Sprintf("The category for the package. Value must match one of the keys in the map: %v", categoryNameMap))
	cmd.Flags().BoolVar(&featured, "featured", false, "Whether or not this package should be marked as featured in its metadata")
	cmd.Flags().StringVar(&publisher, "publisher", "Pulumi", "The publisher's display name to be shown in the package")
	cmd.Flags().StringVar(&title, "title", "", "The display name of the package. If ommitted, the name of the package will be used")
	cmd.Flags().Int64Var(&updatedOn, "updatedOn", time.Now().Unix(), "The timestamp (epoch) to use for when the package was last updated")
	cmd.Flags().BoolVar(&component, "component", false, "Whether or not this package is a component and not a provider")

	cmd.MarkFlagRequired("metadataOutDir")

	return cmd
}
