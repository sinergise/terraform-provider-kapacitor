package kapacitor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/influxdata/kapacitor/client/v1"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAPACITOR_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAPACITOR_PASSWORD", nil),
			},
		},

		ConfigureContextFunc: configure,

		ResourcesMap: map[string]*schema.Resource{
			"kapacitor_task": taskResource(),
		},
	}
}

func configure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := client.Config{
		URL:       d.Get("url").(string),
		UserAgent: "Terraform",
	}

	if _, ok := d.GetOk("username"); ok {
		config.Credentials = &client.Credentials{
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
			Method:   client.UserAuthentication,
		}
	}

	conn, err := client.New(config)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	_, _, err = conn.Ping()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return conn, diags
}
