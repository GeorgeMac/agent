---
title: otelcol.auth.sigv4
---

# otelcol.auth.sigv4

`otelcol.auth.sigv4` performs `sigv4` authentication for making requests to AWS services 
via `otelcol` components that support authentication extensions.

> **NOTE**: `otelcol.auth.sigv4` is a wrapper over the upstream OpenTelemetry
> Collector `sigv4auth` extension. Bug reports or feature requests will be
> redirected to the upstream repository, if necessary.

Multiple `otelcol.auth.sigv4` components can be specified by giving them
different labels.

## Usage

```river
otelcol.auth.sigv4 "LABEL" {
}
```

## Arguments

Name | Type | Description | Default | Required
---- | ---- | ----------- | ------- | --------
`region` | `string` | The AWS region for the service you are exporting to for AWS Sigv4. This is differentiated from sts_region to handle cross region authentication. | | no
`service` | `string` | The AWS service for AWS Sigv4. | | no

## Blocks

The following blocks are supported inside the definition of
`otelcol.auth.sigv4`:

Hierarchy | Block | Description | Required
--------- | ----- | ----------- | --------
assume_role | [assume_role][] | Custom header to attach to requests. | no

[assume_role]: #assume_role-block

### assume_role block

The `assume_role` block specifies the configuration needed to assume a role.

Name | Type | Description | Default | Required
---- | ---- | ----------- | ------- | --------
`arn` | `string` | The Amazon Resource Name (ARN) of a role to assume. | | no
`session_name` | `string` | The name of a role session. | | no
`sts_region` | `string` | The AWS region where STS is used to assumed the configured role. | | no

Note that if a role is intended to be assumed and `sts_region` is not provided, then `sts_region`
will default to the value for `region` if `region` is provided.

## Exported fields

The following fields are exported and can be referenced by other components:

Name | Type | Description
---- | ---- | -----------
`handler` | `capsule(otelcol.Handler)` | A value that other components can use to authenticate requests.

## Component health

`otelcol.auth.sigv4` is only reported as unhealthy if given an invalid
configuration.

## Debug information

`otelcol.auth.sigv4` does not expose any component-specific debug information.

## Example

This example configures [otelcol.exporter.otlp][] to use custom headers:

```river
otelcol.exporter.otlp "example" {
  client {
    endpoint = "my-otlp-grpc-server:4317"
    auth     = otelcol.auth.sigv4.creds.handler
  }
}

otelcol.auth.sigv4 "creds" {
}
```

[otelcol.exporter.otlp]: {{< relref "./otelcol.exporter.otlp.md" >}}
