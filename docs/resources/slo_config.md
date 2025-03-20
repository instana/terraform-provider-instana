# SLO Configuration

An SLO, or Service Level Objective, is a specific, measurable target that defines the expected level of performance, 
reliability, or availability of a service, agreed upon between a service provider and its users or customers.
For instance, an SLO could state that a specific SLI (Service Level Indicator), such as availability, must reach 99.9% 
over a set period.

API Documentation: <https://instana.github.io/openapi/#operation/createSloConfig>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

## Example Usage
Creating an application SLO with timebased latency indicator and rolling time window 

```hcl
resource "instana_slo_config" "slo_1" {
  name = "app_timebased_latency_rolling"
  target = 0.95
  tags = ["terraform", "app", "timebased", "latency", "fixed"]
  entity {
    application {
      application_id = "instana_application_config_id"
      boundary_scope = "ALL"
      include_internal = false
      include_synthetic = false
      filter_expression = "AND"
    }
  }
  indicator {
     time_based_latency {
       threshold = 13.1
       aggregation = "MEAN"
     }
  }
  time_window {
    rolling {
      duration = 1
      duration_unit = "day"
    }
  }
}
```

Creating a website SLO with timebased availability indicator and fixed time window 

```hcl
resource "instana_slo_config" "website_3" {
  name = "website_timebased_availability_fixed"
  target = 0.91
  tags = ["terraform", "web", "timebased", "availability", "fixed"]
  entity {
    website {
     website_id = "instana_website_monitoring_config_id"
     beacon_type = "httpRequest"
     filter_expression = "AND"
    }
  }
  indicator {
    time_based_availability {
      threshold = 14.7
      aggregation = "MEAN"
    }
  }
   time_window {
     fixed {
       duration = 1
       duration_unit = "day"
       start_timestamp = var.fixed_timewindow_start_timestamp
     }
   }
}
```

Creating a synthetic SLO with all traffic indicator and rolling time window 

```hcl
resource "instana_slo_config" "synthetic_r_6" {
  name = "synthetic_traffic_all_rolling"
  target = 0.91
  tags = ["terraform", "synthetic", "traffic", "all", "rolling-time-window"]
  entity {
     synthetic {
       synthetic_test_ids = ["DrMyeGl08w79poguQ3mhH", "sYDtb2slIIolfXhPBnodSz" ]
       filter_expression = "AND"
     } 
  }
  indicator {
    traffic {
      traffic_type = "all"
      threshold = 14
      aggregation = "SUM"
    }
  }
   time_window {
     rolling {
       duration = 1
       duration_unit = "day"
     }
   }
}
``` 

## Argument Reference
