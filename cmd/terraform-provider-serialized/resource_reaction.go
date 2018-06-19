package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	serialized "github.com/marcusolsson/serialized-go"
)

func resourceReaction() *schema.Resource {
	return &schema.Resource{
		Create: resourceReactionCreate,
		Read:   resourceReactionRead,
		Delete: resourceReactionDelete,

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
			"reacts_on": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cancels_on": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"trigger_time_field": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"offset": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"action": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     resourceAction(),
			},
		},
	}
}

func resourceAction() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_uri": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"body": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceReactionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*serialized.Client)

	var cancels []string
	for _, t := range d.Get("cancels_on").([]interface{}) {
		cancels = append(cancels, t.(string))
	}

	def := &serialized.ReactionDefinition{
		Name:               d.Get("name").(string),
		Feed:               d.Get("feed").(string),
		ReactOnEventType:   d.Get("reacts_on").(string),
		CancelOnEventTypes: cancels,
		TriggerTimeField:   d.Get("trigger_time_field").(string),
		Offset:             d.Get("offset").(string),
		Action:             actionFromResourceData(d.Get("action")),
	}

	if err := client.CreateReactionDefinition(context.Background(), def); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("/reactions/%s", def.Name))

	return nil
}

func actionFromResourceData(d interface{}) *serialized.Action {
	cfg := d.([]interface{})[0].(map[string]interface{})

	var (
		actionType = cfg["action_type"].(string)
		targetURI  = cfg["target_uri"].(string)
		body       = cfg["body"].(string)
	)

	return &serialized.Action{
		ActionType: serialized.ActionType(actionType),
		TargetURI:  targetURI,
		Body:       body,
	}
}

func resourceReactionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*serialized.Client)

	name := d.Get("name").(string)

	def, err := client.ReactionDefinition(context.Background(), name)
	if err != nil {
		return nil
	}

	d.Set("name", def.Name)
	d.Set("feed", def.Feed)
	d.Set("reacts_on", def.ReactOnEventType)
	d.Set("cancels_on", def.CancelOnEventTypes)
	d.Set("offset", def.Offset)
	d.Set("trigger_time_field", def.TriggerTimeField)
	d.Set("action", map[string]interface{}{
		"action_type": def.Action.ActionType,
		"target_uri":  def.Action.TargetURI,
		"body":        def.Action.Body,
	})

	return nil
}

func resourceReactionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*serialized.Client)

	name := d.Get("name").(string)

	return client.DeleteReactionDefinition(context.Background(), name)
}
