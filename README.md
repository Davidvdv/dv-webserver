# dv-webserver

A simple HTTP web server written in Go that serves static HTML content.

## Description

dv-webserver is a lightweight TCP-based web server that listens for HTTP connections and serves static HTML content from
the `public` directory. It supports basic HTTP request parsing and can handle multiple concurrent connections.

## Requirements

- Go 1.24 or higher
- A valid HTML file in the `public/index.html` location
