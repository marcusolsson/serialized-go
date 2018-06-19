package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	serialized "github.com/marcusolsson/serialized-go"
)

func resourceProjection() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectionCreate,
		Read:   resourceProjectionRead,
		Delete: resourceProjectionDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"feed": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"handlers": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     resourceEventHandler(),
			},
		},
	}
}

func resourceEventHandler() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"event_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"functions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     resourceFunction(),
			},
		},
	}
}

func resourceFunction() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"function": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_selector": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"event_selector": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"event_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"raw_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceProjectionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*serialized.Client)

	var handlers []*serialized.EventHandler
	for _, h := range d.Get("handlers").([]interface{}) {
		handlers = append(handlers, handlerFromResourceData(h))
	}

	def := &serialized.ProjectionDefinition{
		Name:     d.Get("name").(string),
		Feed:     d.Get("feed").(string),
		Handlers: handlers,
	}

	if err := client.CreateProjectionDefinition(context.Background(), def); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("/projections/definitions/%s", def.Name))

	return nil
}

func handlerFromResourceData(d interface{}) *serialized.EventHandler {
	hh := d.(map[string]interface{})

	var fns []*serialized.Function
	for _, f := range hh["functions"].([]interface{}) {
		ff := f.(map[string]interface{})
		fns = append(fns, &serialized.Function{
			Function:       ff["function"].(string),
			TargetSelector: ff["target_selector"].(string),
			EventSelector:  ff["event_selector"].(string),
			TargetFilter:   ff["target_filter"].(string),
			EventFilter:    ff["event_filter"].(string),
			RawData:        ff["raw_data"].(string),
		})
	}

	return &serialized.EventHandler{
		EventType: hh["event_type"].(string),
		Functions: fns,
	}
}

func resourceProjectionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*serialized.Client)

	name := d.Get("name").(string)

	def, err := client.ProjectionDefinition(context.Background(), name)
	if err != nil {
		return nil
	}

	d.Set("name", def.Name)
	d.Set("feed", def.Feed)

	var hs []interface{}
	for _, h := range def.Handlers {
		hdata := make(map[string]interface{})
		hdata["event_type"] = h.EventType

		var fs []interface{}
		for _, f := range h.Functions {
			fdata := make(map[string]interface{})
			fdata["function"] = f.Function
			fdata["target_selector"] = f.TargetSelector
			fdata["event_selector"] = f.EventSelector
			fdata["target_filter"] = f.TargetFilter
			fdata["event_filter"] = f.EventFilter
			fdata["raw_data"] = f.RawData

			fs = append(fs, fdata)
		}
		hdata["functions"] = fs

		hs = append(hs, hdata)
	}

	d.Set("handlers", hs)

	return nil
}

func resourceProjectionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*serialized.Client)

	name := d.Get("name").(string)

	return client.DeleteProjectionDefinition(context.Background(), name)
}
