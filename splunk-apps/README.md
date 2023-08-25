# Splunk apps for compliance-audit-router

This directory contains Splunk apps for use with the [compliance-audit-router](../README.md) tool.

## Building and deploying

To build the apps, run `make` in this directory. The resulting .spl file can be deployed to Splunk cloud from the web UI.

### auth_webhook

Provides a webhook alert action that supports authentication