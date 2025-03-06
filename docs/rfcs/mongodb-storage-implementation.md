# Summary

This RFC proposes implementing MongoDB as an alternative storage solution for Bucketeer, specifically targeting open-source deployments and smaller companies. This change aims to simplify the deployment process using Docker Compose while maintaining core functionality. The document focuses on the implementation of pagination and discusses the trade-offs between different database solutions.

# Background

Currently, Bucketeer uses a combination of MySQL and BigQuery for data storage and analytics. While this architecture serves well for large-scale deployments, it presents complexity in setup and maintenance for smaller deployments. MongoDB offers a simpler alternative that can be easily containerized while providing sufficient functionality for smaller scale deployments.

# Goals

- Simplify deployment process for open-source users
- Maintain core functionality including efficient pagination
- Ensure consistent performance for small to medium-scale deployments


# Implementation Details

## MongoDB Schema Design

The MongoDB implementation will maintain similar data structures to the current MySQL schema, with adaptations to leverage MongoDB's document model. 

## SQL to MongoDB Query Migration

### Basic Query Patterns

1. **Feature Queries**
   ```sql
   -- SQL
   SELECT * FROM feature 
   WHERE environment_namespace = 'default' 
   AND archived = false 
   AND deleted = false;

   -- MongoDB
   db.feature.find({
     environment_namespace: "default",
     archived: false,
     deleted: false
   })
   ```

2. **Project with Organization**
   ```sql
   -- SQL
   SELECT p.*, o.name as org_name 
   FROM project p
   JOIN organization o ON p.organization_id = o.id
   WHERE p.disabled = false;

   -- MongoDB
   db.project.aggregate([
     {
       $match: { disabled: false }
     },
     {
       $lookup: {
         from: "organization",
         localField: "organization_id",
         foreignField: "id",
         as: "organization"
       }
     },
     {
       $unwind: "$organization"
     }
   ])
   ```

3. **Feature with Variations**
   ```sql
   -- SQL
   SELECT * FROM feature 
   WHERE environment_namespace = 'default'
   ORDER BY created_at DESC;

   -- MongoDB
   db.feature.find({
     environment_namespace: "default"
   }).sort({
     created_at: -1
   })
   ```

### Complex Query Patterns

1. **Experiment with Goals**
   ```sql
   -- SQL
   SELECT e.*, g.name as goal_name 
   FROM experiment e
   JOIN goal g ON e.goal_id = g.id 
   WHERE e.environment_namespace = 'default' 
   AND e.archived = false;

   -- MongoDB
   db.experiment.aggregate([
     {
       $match: {
         environment_namespace: "default",
         archived: false
       }
     },
     {
       $lookup: {
         from: "goal",
         localField: "goal_id",
         foreignField: "id",
         as: "goal"
       }
     }
   ])
   ```



2. **Segment Users Query**
   ```sql
   -- SQL
   SELECT s.*, COUNT(su.user_id) as user_count
   FROM segment s
   LEFT JOIN segment_user su ON s.id = su.segment_id
   WHERE s.environment_namespace = 'default'
   AND s.deleted = false
   GROUP BY s.id;

   -- MongoDB
   db.segment.aggregate([
     {
       $match: {
         environment_namespace: "default",
         deleted: false
       }
     },
     {
       $lookup: {
         from: "segment_user",
         localField: "id",
         foreignField: "segment_id",
         as: "users"
       }
     },
     {
       $project: {
         id: 1,
         name: 1,
         description: 1,
         rules: 1,
         user_count: { $size: "$users" }
       }
     }
   ])
   ```


### Join Queries

```sql
-- SQL
SELECT r.*, f.name as feature_name
FROM auto_ops_rule r
JOIN feature f ON r.feature_id = f.id
WHERE r.environment_namespace = 'default'
AND r.deleted = false;

-- MongoDB
db.auto_ops_rule.aggregate([
  {
    $match: {
      environment_namespace: "default",
      deleted: false
    }
  },
  {
    $lookup: {
      from: "feature",
      localField: "feature_id",
      foreignField: "id",
      as: "feature"
    }
  },
  {
    $unwind: "$feature"
  }
])
```

## Pagination Implementation

### Cursor-based Pagination

We will implement cursor-based pagination using MongoDB's native capabilities:

```javascript
{
  find: "collection_name",
  sort: { _id: 1 },
  limit: pageSize,
  filter: {
    _id: { $gt: lastSeenId }
  }
}
```

Benefits:
- Consistent performance regardless of offset
- Works well with real-time data
- Maintains performance with large datasets

### Offset-based Pagination (Alternative)

For simpler use cases, we'll also support traditional offset-based pagination:

```javascript
{
  find: "collection_name",
  skip: offset,
  limit: pageSize
}
```

## Performance Considerations

### Indexes

Essential indexes for pagination performance:

```javascript
// Primary indexes for pagination
db.features.createIndex({ created_at: -1, _id: 1 })
db.evaluations.createIndex({ feature_id: 1, created_at: -1 })
```

# Trade-offs

## Advantages

1. **Simplified Deployment**
   - Single database solution
   - Easy containerization
   - Reduced infrastructure complexity

2. **Schema Flexibility**
   - Easier schema evolution
   - Better handling of optional fields
   - Simplified document updates

3. **Developer Experience**
   - Familiar JSON-like document model
   - Rich query language
   - Strong community support

## Disadvantages

1. **Limited Analytics Capabilities**
   - No direct replacement for BigQuery analytics
   - Complex aggregations may be slower
   - Limited support for complex joins

2. **Eventual Consistency**
   - Default eventual consistency model
   - May require careful configuration for strong consistency

3. **Resource Requirements**
   - Higher memory usage compared to MySQL
   - Needs careful index planning

