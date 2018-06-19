# Terraform Provider for Serialized.IO

## serialized_projection

Creates a projection for Serialized.IO. For more information see
[the official documentation](https://serialized.io/docs/apis/event-projection/) and
[API](https://serialized.io/api/#tag-Event-Projection-API).


### Example Usage

```hcl
resource "serialized_projection" "default" {
  name = "orders"
  feed = "order"

  handlers = [
    {
      event_type = "OrderCancelledEvent"

      functions = [
        {
          function        = "inc"
          target_selector = "$.projection.orders[?]"
          event_selector  = "$.event[?]"
          target_filter   = "@.orderId == $.event.orderId"
          event_filter    = "@.orderAmount > 4000"
        },
      ]
    },
  ]
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) A unique ID for the resource.
    Changing this forces a new resource to be created.

* `feed` - (Required) The ID of the feed.

* `handlers` - (Required) The event handlers.

## serialized_reaction

Creates a reaction for Serialized.IO. For more information see
[the official documentation](https://serialized.io/docs/apis/event-reaction/) and
[API](https://serialized.io/api/#tag-Event-Reaction-API).


### Example Usage

```hcl
resource "serialized_reaction" "default" {
  name      = "payment-processed-email-reaction"
  feed      = "payment"
  reacts_on = "PaymentProcessed"

  cancels_on = [
    "OrderCanceledEvent",
  ]

  trigger_time_field = "my.event.data.field"
  offset             = "PT1H"

  action {
    action_type = "HTTP_POST"
  }
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) A unique ID for the resource.
    Changing this forces a new resource to be created.

* `feed` - (Required) The ID of the feed.