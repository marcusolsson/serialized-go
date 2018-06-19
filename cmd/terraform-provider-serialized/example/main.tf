provider "serialized" {
  # NOTE: You can also use the SERIALIZED_ACCESS_KEY and       
  # SERIALIZE_SECRET_ACCESS_KEY environment variables to set credentials.  

  # access_key        = "..."    
  # secret_access_key = "..."
}

resource "serialized_projection" "example" {
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

resource "serialized_reaction" "example" {
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
