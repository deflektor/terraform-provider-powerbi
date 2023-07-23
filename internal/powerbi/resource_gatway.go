package powerbi

import (
	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ResourceGateways represent managment of DataSources on On-premises gatways.
func ResourceGateways() *schema.Resource {
	return &schema.Resource{
		Read: getGateway,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"gatewayId": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The gateway ID. When using a gateway cluster, the gateway ID refers to the primary (first) gateway in the cluster. In such cases, gateway ID is similar to gateway cluster ID.",
			},
			"gateway": {
				Type:        schema.TypeSet,
				Description: "Gatway Definition",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "The gateway ID. When using a gateway cluster, the gateway ID refers to the primary (first) gateway in the cluster and is similar to the gateway cluster ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The gateway name",
						},
						"type": {
							Type:        schema.TypeString,
							Description: "The gateway type",
						},
						"gatewayStatus": {
							Type:        schema.TypeString,
							Description: "The gateway connectivity status",
						},
						"gatewayAnnotation": {
							Type:        schema.TypeString,
							Description: "Gateway metadata in JSON format",
						},
					},
				},
			},
		},

		//Create:   createDatasource,
		//Read:     getDatasourceStatus,
		//Delete:   deleteDatasource,
		//Importer: getDatasources,
		///DeleteDatasourceUser
		///AddDatasourceUser
		//GetGateways
		//GetGateway
		////// from here the importer
		//GetDatasources
		//GetDatasource
		//GetDatasourceStatus
		//GetDatasourceUsers
	}
}

func getGateway(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	gatewayId := d.Get("gatewayId").(string)
	gateway, err := client.GetGateway(gatewayId)
	if err != nil {
		return err
	}

	if gateway == nil {
		d.SetId("")
	} else {
		d.SetId(gateway.ID)
		d.Set("name", gateway.Name)
		d.Set("gatewayAnnotation", gateway.GatewayAnnotation)
		d.Set("type", gateway.Type)
		d.Set("gatewayStatus", gateway.GatewayStatus)
		//todo
		/*		if gateway.PublicKey {
					d.Set("exponent", gateway.PublicKey.exponent)
					d.Set("modulus", gateway.PublicKey.modulus)
				} else {
					d.Set("exponent", "")
					d.Set("modulus", "")
				}
		*/
	}

	return nil
}

// todo
func getGateways(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	gateways, err := client.GetGateways()
	if err != nil {
		return err
	}

	if gateways == nil {

	}
	return nil
}
