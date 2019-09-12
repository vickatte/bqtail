### Export data on destination table modification 

### Scenario:

BqDispatch function is notified with all Big Query jobs completion, it matches actions to run
with the following route to export destination table to google storage gs://${config.Bucket}/export/dummy.json.gz


```json
  {
    "When": {
      "Dest": ".+:bqdispatch\\.dummy_v2",
      "Type": "QUERY"
    },
    "OnSuccess": [
      {
        "Action": "export",
        "Request": {
          "DestURL": "gs://${config.Bucket}/export/dummy.json.gz"
        }
      }
    ]
  }
```

### BqDispatch


#### Input:

* **Trigger**:
  - eventType: google.cloud.bigquery.job.complete
  - resource: projects/${projectID}/jobs/{jobId}
* **Configuration:** [gs://e2e-data/config/bqdispatch.json](../../../config/bqdispatch.json)
* **Data**:

```sql
    INSERT INTO bqdispatch.dummy_v2 AS SELECT * FROM bqdispatch.dummy_v1
```

#### Output

* **Logs:** 

- gs://${config.Bucket}/journal/

* **Data:**
- gs://${config.Bucket}/export/dummy.json.gz