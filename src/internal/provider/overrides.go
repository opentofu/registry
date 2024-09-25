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
		"default/terraform-provider-aci":              "ciscodevnet/terraform-provider-aci",
		"default/terraform-provider-acme":             "vancluever/terraform-provider-acme",
		"default/terraform-provider-akamai":           "akamai/terraform-provider-akamai",
		"default/terraform-provider-alicloud":         "aliyun/terraform-provider-alicloud",
		"default/terraform-provider-aviatrix":         "aviatrixsystems/terraform-provider-aviatrix",
		"default/terraform-provider-avi":              "vmware/terraform-provider-avi",
		"default/terraform-provider-azuredevops":      "microsoft/terraform-provider-azuredevops",
		"default/terraform-provider-baiducloud":       "baidubce/terraform-provider-baiducloud",
		"default/terraform-provider-bigip":            "f5networks/terraform-provider-bigip",
		"default/terraform-provider-brightbox":        "brightbox/terraform-provider-brightbox",
		"default/terraform-provider-checkpoint":       "checkpointsw/terraform-provider-checkpoint",
		"default/terraform-provider-circonus":         "circonus-labs/terraform-provider-circonus",
		"default/terraform-provider-cloudflare":       "cloudflare/terraform-provider-cloudflare",
		"default/terraform-provider-cloudscale":       "cloudscale-ch/terraform-provider-cloudscale",
		"default/terraform-provider-constellix":       "constellix/terraform-provider-constellix",
		"default/terraform-provider-datadog":          "datadog/terraform-provider-datadog",
		"default/terraform-provider-digitalocean":     "digitalocean/terraform-provider-digitalocean",
		"default/terraform-provider-dme":              "dnsmadeeasy/terraform-provider-dme",
		"default/terraform-provider-dnsimple":         "dnsimple/terraform-provider-dnsimple", // manually detected from incorrect homepage
		"default/terraform-provider-dome9":            "dome9/terraform-provider-dome9",
		"default/terraform-provider-exoscale":         "exoscale/terraform-provider-exoscale",
		"default/terraform-provider-fastly":           "fastly/terraform-provider-fastly",
		"default/terraform-provider-flexibleengine":   "flexibleenginecloud/terraform-provider-flexibleengine",
		"default/terraform-provider-fortios":          "fortinetdev/terraform-provider-fortios",
		"default/terraform-provider-github":           "integrations/terraform-provider-github",
		"default/terraform-provider-gitlab":           "gitlabhq/terraform-provider-gitlab",
		"default/terraform-provider-grafana":          "grafana/terraform-provider-grafana",
		"default/terraform-provider-gridscale":        "gridscale/terraform-provider-gridscale",
		"default/terraform-provider-hcloud":           "hetznercloud/terraform-provider-hcloud",
		"default/terraform-provider-heroku":           "heroku/terraform-provider-heroku",
		"default/terraform-provider-huaweicloud":      "huaweicloud/terraform-provider-huaweicloud",
		"default/terraform-provider-huaweicloudstack": "huaweicloud/terraform-provider-huaweicloudstack",
		"default/terraform-provider-icinga2":          "icinga/terraform-provider-icinga2",
		"default/terraform-provider-launchdarkly":     "launchdarkly/terraform-provider-launchdarkly",
		"default/terraform-provider-linode":           "linode/terraform-provider-linode",
		"default/terraform-provider-logicmonitor":     "logicmonitor/terraform-provider-logicmonitor", // manually detected from incorrect homepage
		"default/terraform-provider-mongodbatlas":     "mongodb/terraform-provider-mongodbatlas",
		"default/terraform-provider-ncloud":           "navercloudplatform/terraform-provider-ncloud",
		"default/terraform-provider-newrelic":         "newrelic/terraform-provider-newrelic",
		"default/terraform-provider-ns1":              "ns1-terraform/terraform-provider-ns1",
		"default/terraform-provider-nsxt":             "vmware/terraform-provider-nsxt",
		"default/terraform-provider-nutanix":          "nutanix/terraform-provider-nutanix",
		"default/terraform-provider-oci":              "oracle/terraform-provider-oci",
		"default/terraform-provider-oktaasa":          "oktadeveloper/terraform-provider-oktaasa",
		"default/terraform-provider-okta":             "oktadeveloper/terraform-provider-okta",
		"default/terraform-provider-opennebula":       "opennebula/terraform-provider-opennebula",
		"default/terraform-provider-openstack":        "terraform-provider-openstack/terraform-provider-openstack",
		"default/terraform-provider-opentelekomcloud": "opentelekomcloud/terraform-provider-opentelekomcloud",
		"default/terraform-provider-opsgenie":         "opsgenie/terraform-provider-opsgenie",
		"default/terraform-provider-ovh":              "ovh/terraform-provider-ovh",
		"default/terraform-provider-packet":           "packethost/terraform-provider-packet",
		"default/terraform-provider-pagerduty":        "pagerduty/terraform-provider-pagerduty",
		"default/terraform-provider-panos":            "paloaltonetworks/terraform-provider-panos",
		"default/terraform-provider-powerdns":         "pan-net/terraform-provider-powerdns",
		"default/terraform-provider-prismacloud":      "paloaltonetworks/terraform-provider-prismacloud",
		"default/terraform-provider-profitbricks":     "ionos-cloud/terraform-provider-profitbricks",
		"default/terraform-provider-rancher2":         "rancher/terraform-provider-rancher2",
		"default/terraform-provider-rundeck":          "rundeck/terraform-provider-rundeck",
		"default/terraform-provider-scaleway":         "scaleway/terraform-provider-scaleway",
		"default/terraform-provider-selectel":         "selectel/terraform-provider-selectel",
		"default/terraform-provider-signalfx":         "splunk-terraform/terraform-provider-signalfx", // repo was moved "signalfx/terraform-provider-signalfx",
		"default/terraform-provider-skytap":           "skytap/terraform-provider-skytap",
		"default/terraform-provider-spotinst":         "spotinst/terraform-provider-spotinst",
		"default/terraform-provider-stackpath":        "stackpath/terraform-provider-stackpath",
		"default/terraform-provider-statuscake":       "statuscakedev/terraform-provider-statuscake",
		"default/terraform-provider-sumologic":        "sumologic/terraform-provider-sumologic",
		"default/terraform-provider-tencentcloud":     "tencentcloudstack/terraform-provider-tencentcloud",
		"default/terraform-provider-triton":           "joyent/terraform-provider-triton",
		"default/terraform-provider-turbot":           "turbot/terraform-provider-turbot",
		"default/terraform-provider-ucloud":           "ucloud/terraform-provider-ucloud",
		"default/terraform-provider-vcd":              "vmware/terraform-provider-vcd",
		"default/terraform-provider-venafi":           "venafi/terraform-provider-venafi",
		"default/terraform-provider-vmc":              "vmware/terraform-provider-vmc",
		"default/terraform-provider-vra7":             "vmware/terraform-provider-vra7",
		"default/terraform-provider-vultr":            "vultr/terraform-provider-vultr",
		"default/terraform-provider-wavefront":        "vmware/terraform-provider-wavefront",
		"default/terraform-provider-yandex":           "yandex-cloud/terraform-provider-yandex",
		// inaccessable "default/terraform-provider-genymotion":       "genymobile/terraform-provider-genymotion",
		// inaccessable "default/terraform-provider-cherryservers":    "arturasra/terraform-provider-cherryservers", // manually detected from incorrect homepage
	}
)
