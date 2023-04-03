This feature is inspired by https://launchdarkly.com/blog/take-the-maintenance-out-of-maintenance-windows/.

# Controversial topics

1. Creating a new command vs Inserting a new field to DatetimeClause

Creating a new command

```proto
message AddScheduledDatetimeClauseCommand {
  DatetimeClause datetime_clause = 1;
  int64 duration = 2;
}

message ChangeScheduledDatetimeClauseCommand {
  string id = 1;
  DatetimeClause datetime_clause = 2;
  int64 duration = 3;
}
```

Inserting a new field to DatetimeClause

```proto
message DatetimeClause {
  int64 time = 1;
  int64 duration = 2;
}
```

2. seconds vs minutes vs hours

Which notation is appropriate for the duration field?

# Server-side changes

## Auto ops

Basically we change pkg/autoops/domain/auto_ops_rule.go to impelemt this feature.
