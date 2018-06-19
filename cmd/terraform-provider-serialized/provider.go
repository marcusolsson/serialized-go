package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	serialized "github.com/marcusolsson/serialized-go"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERIALIZED_ACCESS_KEY", nil),
				Description: "Serialized.IO Access Key",
			},
			"secret_access_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERIALIZED_SECRET_ACCESS_KEY", nil),
				Description: "Serialized.IO Secret Access Key",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"serialized_projection": resourceProjection(),
			"serialized_reaction":   resourceReaction(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	return serialized.NewClient(
		serialized.WithAccessKey(d.Get("access_key").(string)),
		serialized.WithSecretAccessKey(d.Get("secret_access_key").(string)),
	), nil
}
