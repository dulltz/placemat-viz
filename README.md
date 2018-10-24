# placemat-viz

[![CircleCI](https://circleci.com/gh/cybozu-go/cke.svg?style=svg)](https://circleci.com/gh/cybozu-go/cke)

Visualization tool for Placemat.

## Requirements

- Graphviz (with `dot` CLI command)
- `cluster.yml` 
  - resource file for [Placemat](https://github.com/cybozu-go/placemat). The example is [here](https://raw.githubusercontent.com/cybozu-go/placemat-menu/master/testdata/cluster.yml).

## Usage 

Generate a figure from cluster.yml

```console
$ placemat-viz --input cluster.yml | dot -T svg > output.svg
```

