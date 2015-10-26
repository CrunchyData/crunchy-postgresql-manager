// Package collect is the implementation of the CPM metrics
// collection process.  Metrics are collected and then
// written into a Prometheus data store for reporting.
// Healthcheck metrics when collected are written into the CPM administrative
// postgresql database.  The collect server runs and
// polls for metrics on a simple frequency.  Every few
// minutes, metrics are collected and written to a persistent
// store.
package collect
