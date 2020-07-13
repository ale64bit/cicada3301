# Intro

This is a harness I built for myself for experimenting with the cicada puzzles. 
It's written in Golang. See [here](https://golang.org/doc/install) how to install it.
It has no external dependencies and should run in all supported runtimes (Windows, Mac, Linux, etc.).

# How to use it

There are two main binaries:

## done

The `done` binary shows the already decoded sections along with the information related to the decoder. This is useful to verify the already-solved pages and to try tweaks to existing decoders.

Example output for running the command `go run bin/done/main.go` under Linux:

![Alt text](/doc/done.png?raw=true "Done")

Additionally, you can restrict which sections to decode by passing a space-separated list of section IDs (see [data/data.go](/data/data.go)) as arguments to the binary.

## search
The `search` binary builds a set of decoders of various types and tries to decode the unsolved sections. The result is evaluated according to a dictionary search and scored accordingly.

Example output for running the command `go run bin/search/main.go` under Linux:

![Alt text](/doc/search.png?raw=true "Search")

Additionally, you can pass the following flags:
* `-prefix_len`: specify the length of the prefix to try to decode from each unsolved section. Longer prefixes result in slower searches.
* `-selected_sections`: a comma-separated list of the unsolved sections to try to decode. If unspecified, all unsolved sections are tried.
* `-match_score`: the minimum score that is considered a match (i.e. successfully decoded). Scores range from 0 to 1.
