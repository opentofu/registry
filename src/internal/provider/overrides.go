package provider

var (
	_ = `
	rm providers-hashicorp.json
	url="https://api.github.com/orgs/hashicorp/repos?type=all&per_page=100&page=1"
	while [ ! -z "$url" ]; do
		echo "Fetching $url"
		curl -H "Authorization: Bearer $GH_TOKEN" -H "Accept: application/vnd.github.text-match+json" "$url" -o /tmp/github.json -w '%{header_json}' > /tmp/headers.json
		cat /tmp/github.json >> providers-hashicorp.json
		url=$(cat /tmp/headers.json | jq '.link[0]' | tr ',' '\n' | grep next | sed -e 's/.*<\(.*\)>.*/\1/')
		sleep 1
	done
	cat providers-hashicorp..json | jq -r '.[] | select(.fork) | select(.full_name | contains("hashicorp/terraform-provider")) | "\"\(.full_name)\": \"\(.homepage)\","' | sed -e 's,https://github.com/,,'
`
	ArchivedOverrides = map[string]string{
		"hashicorp/terraform-provider-aci":              "CiscoDevNet/terraform-provider-aci",
		"hashicorp/terraform-provider-acme":             "vancluever/terraform-provider-acme",
		"hashicorp/terraform-provider-akamai":           "akamai/terraform-provider-akamai",
		"hashicorp/terraform-provider-alicloud":         "aliyun/terraform-provider-alicloud",
		"hashicorp/terraform-provider-aviatrix":         "AviatrixSystems/terraform-provider-aviatrix",
		"hashicorp/terraform-provider-avi":              "vmware/terraform-provider-avi",
		"hashicorp/terraform-provider-azuredevops":      "microsoft/terraform-provider-azuredevops",
		"hashicorp/terraform-provider-baiducloud":       "baidubce/terraform-provider-baiducloud",
		"hashicorp/terraform-provider-bigip":            "F5Networks/terraform-provider-bigip",
		"hashicorp/terraform-provider-brightbox":        "brightbox/terraform-provider-brightbox",
		"hashicorp/terraform-provider-checkpoint":       "CheckPointSW/terraform-provider-checkpoint",
		"hashicorp/terraform-provider-circonus":         "circonus-labs/terraform-provider-circonus",
		"hashicorp/terraform-provider-cloudflare":       "cloudflare/terraform-provider-cloudflare",
		"hashicorp/terraform-provider-cloudscale":       "cloudscale-ch/terraform-provider-cloudscale",
		"hashicorp/terraform-provider-constellix":       "Constellix/terraform-provider-constellix",
		"hashicorp/terraform-provider-datadog":          "DataDog/terraform-provider-datadog",
		"hashicorp/terraform-provider-digitalocean":     "digitalocean/terraform-provider-digitalocean",
		"hashicorp/terraform-provider-dme":              "DNSMadeEasy/terraform-provider-dme",
		"hashicorp/terraform-provider-dnsimple":         "dnsimple/terraform-provider-dnsimple", // Manually detected from incorrect homepage
		"hashicorp/terraform-provider-dome9":            "dome9/terraform-provider-dome9",
		"hashicorp/terraform-provider-exoscale":         "exoscale/terraform-provider-exoscale",
		"hashicorp/terraform-provider-fastly":           "fastly/terraform-provider-fastly",
		"hashicorp/terraform-provider-flexibleengine":   "FlexibleEngineCloud/terraform-provider-flexibleengine",
		"hashicorp/terraform-provider-fortios":          "fortinetdev/terraform-provider-fortios",
		"hashicorp/terraform-provider-github":           "integrations/terraform-provider-github",
		"hashicorp/terraform-provider-gitlab":           "gitlabhq/terraform-provider-gitlab",
		"hashicorp/terraform-provider-grafana":          "grafana/terraform-provider-grafana",
		"hashicorp/terraform-provider-gridscale":        "gridscale/terraform-provider-gridscale",
		"hashicorp/terraform-provider-hcloud":           "hetznercloud/terraform-provider-hcloud",
		"hashicorp/terraform-provider-heroku":           "heroku/terraform-provider-heroku",
		"hashicorp/terraform-provider-huaweicloud":      "huaweicloud/terraform-provider-huaweicloud",
		"hashicorp/terraform-provider-huaweicloudstack": "huaweicloud/terraform-provider-huaweicloudstack",
		"hashicorp/terraform-provider-icinga2":          "Icinga/terraform-provider-icinga2",
		"hashicorp/terraform-provider-launchdarkly":     "launchdarkly/terraform-provider-launchdarkly",
		"hashicorp/terraform-provider-linode":           "linode/terraform-provider-linode",
		"hashicorp/terraform-provider-logicmonitor":     "logicmonitor/terraform-provider-logicmonitor", // Manually detected from incorrect homepage
		"hashicorp/terraform-provider-mongodbatlas":     "mongodb/terraform-provider-mongodbatlas",
		"hashicorp/terraform-provider-ncloud":           "NaverCloudPlatform/terraform-provider-ncloud",
		"hashicorp/terraform-provider-newrelic":         "newrelic/terraform-provider-newrelic",
		"hashicorp/terraform-provider-ns1":              "ns1-terraform/terraform-provider-ns1",
		"hashicorp/terraform-provider-nsxt":             "vmware/terraform-provider-nsxt",
		"hashicorp/terraform-provider-nutanix":          "nutanix/terraform-provider-nutanix",
		"hashicorp/terraform-provider-oktaasa":          "oktadeveloper/terraform-provider-oktaasa",
		"hashicorp/terraform-provider-okta":             "oktadeveloper/terraform-provider-okta",
		"hashicorp/terraform-provider-opennebula":       "OpenNebula/terraform-provider-opennebula",
		"hashicorp/terraform-provider-openstack":        "terraform-provider-openstack/terraform-provider-openstack",
		"hashicorp/terraform-provider-opentelekomcloud": "opentelekomcloud/terraform-provider-opentelekomcloud",
		"hashicorp/terraform-provider-opsgenie":         "opsgenie/terraform-provider-opsgenie",
		"hashicorp/terraform-provider-ovh":              "ovh/terraform-provider-ovh",
		"hashicorp/terraform-provider-packet":           "packethost/terraform-provider-packet",
		"hashicorp/terraform-provider-pagerduty":        "PagerDuty/terraform-provider-pagerduty",
		"hashicorp/terraform-provider-panos":            "PaloAltoNetworks/terraform-provider-panos",
		"hashicorp/terraform-provider-powerdns":         "pan-net/terraform-provider-powerdns",
		"hashicorp/terraform-provider-prismacloud":      "PaloAltoNetworks/terraform-provider-prismacloud",
		"hashicorp/terraform-provider-profitbricks":     "ionos-cloud/terraform-provider-profitbricks",
		"hashicorp/terraform-provider-rancher2":         "rancher/terraform-provider-rancher2",
		"hashicorp/terraform-provider-rundeck":          "rundeck/terraform-provider-rundeck",
		"hashicorp/terraform-provider-scaleway":         "scaleway/terraform-provider-scaleway",
		"hashicorp/terraform-provider-selectel":         "selectel/terraform-provider-selectel",
		"hashicorp/terraform-provider-signalfx":         "splunk-terraform/terraform-provider-signalfx", // Repo was moved "signalfx/terraform-provider-signalfx",
		"hashicorp/terraform-provider-skytap":           "skytap/terraform-provider-skytap",
		"hashicorp/terraform-provider-spotinst":         "spotinst/terraform-provider-spotinst",
		"hashicorp/terraform-provider-stackpath":        "stackpath/terraform-provider-stackpath",
		"hashicorp/terraform-provider-statuscake":       "StatusCakeDev/terraform-provider-statuscake",
		"hashicorp/terraform-provider-sumologic":        "SumoLogic/terraform-provider-sumologic",
		"hashicorp/terraform-provider-tencentcloud":     "tencentcloudstack/terraform-provider-tencentcloud",
		"hashicorp/terraform-provider-triton":           "joyent/terraform-provider-triton",
		"hashicorp/terraform-provider-turbot":           "turbot/terraform-provider-turbot",
		"hashicorp/terraform-provider-ucloud":           "ucloud/terraform-provider-ucloud",
		"hashicorp/terraform-provider-vcd":              "vmware/terraform-provider-vcd",
		"hashicorp/terraform-provider-venafi":           "Venafi/terraform-provider-venafi",
		"hashicorp/terraform-provider-vmc":              "vmware/terraform-provider-vmc",
		"hashicorp/terraform-provider-vra7":             "vmware/terraform-provider-vra7",
		"hashicorp/terraform-provider-vultr":            "vultr/terraform-provider-vultr",
		"hashicorp/terraform-provider-wavefront":        "vmware/terraform-provider-wavefront",
		"hashicorp/terraform-provider-yandex":           "yandex-cloud/terraform-provider-yandex",
		// Inaccessable "hashicorp/terraform-provider-genymotion":       "Genymobile/terraform-provider-genymotion",
		// Inaccessable "hashicorp/terraform-provider-cherryservers":    "ArturasRa/terraform-provider-cherryservers", // Manually detected from incorrect homepage
	}
)
