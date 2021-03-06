### Data aggregation with multi step

### Scenario:

This scenario test data aggregation with multi step transformation.

In this scenario incoming batched data is loaded to transient tables, and copied to dest.table
For each successful batch, batch data is aggregated and exported to GS trigger bucket 
to be finally ingested with [performance.yaml](rule/performance.yaml) rule.

It uses the following rule:

- [@transaction.json](rule/transaction.yaml)
```yaml
When:
  Prefix: "/data/case${parentIndex}"
  Suffix: ".json"
Async: true
Batch:
  Window:
    DurationInSec: 10
Dest:
  Table: bqtail.transactions_v${parentIndex}_$Mod(2)
  Transient:
    Dataset: temp
    Alias: t
  Transform:
    charge: (CASE WHEN type_id = 1 THEN t.payment + f.value WHEN type_id = 2 THEN t.payment * (1 + f.value) END)
  SideInputs:
    - Table: bqtail.fees
      Alias: f
      'On': t.fee_id = f.id
OnSuccess:
  - Action: query
    Request:
      SQL: SELECT
        DATE(timestamp) AS date,
        sku_id,
        supply_entity_id,
        MAX($EventID) AS batch_id,
        SUM( payment) payment,
        SUM((CASE WHEN type_id = 1 THEN t.payment + f.value WHEN type_id = 2 THEN t.payment * (1 + f.value) END)) charge,
        SUM(COALESCE(qty, 1.0)) AS qty
        FROM $TempTable t
        LEFT JOIN bqtail.fees f ON f.id = t.fee_id
        GROUP BY 1, 2, 3
      Dest: ${ProjectID}.temp.supply_performance_v${EventID}
      Append: false
    OnSuccess:
      - Action: export
        Request:
          Source: ${ProjectID}.temp.supply_performance_v${EventID}
          DestURL: gs://${TriggerBucket}/data/supply_performance/transient-*.avro
          Format: AVRO
          UseAvroLogicalTypes: true
        OnSuccess:
          - Action: delete #delete batch files
          - Action: drop
            Request:
              Table: ${ProjectID}.temp.supply_performance_v${EventID}
```

- [@performance.json](rule/performance.yaml)
```yaml
When:
  Prefix: "/data/supply_performance"
  Suffix: ".avro"
Async: true
Batch:
  Window:
    DurationInSec: 10
Dest:
  Table: bqtail.supply_performance_v${parentIndex}
OnSuccess:
  - Action: delete
```

Note that storage trigger $EventID is used as batch id.


