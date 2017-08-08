# Fast and flexible URL shortener

Shorten is a URL shortening service with a pluggable backend and customizable short URL generation strategies,
with the ability to record hits (ie. clicks).

The service is written in Go and has built-in support for the Redis backend.

## Quickstart

Install the package:

    go get github.com/senko/shorten

Shorten a link using a command-line tool:

    shorten -shorten http://example.com/

Expand a link using a command-line tool:

    shorten -expand <short-key>

Start a HTTP microservice for link expansion:

    shortserver -listen :9000

Visit the shortened URL and verify the redirection to full URL:

    curl -v http://localhost:9000/<short-key>

## Copyright and license

See the LICENSE.txt file for details
