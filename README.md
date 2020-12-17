# jcasc-validator

A simple tool to validate the generated jcasc configuration against a JCASC schema.

## To run with examples

Download/Fork this repo and build the binary.

```
make build
```

Run the validator against the provided test resources

```
./build/jcasc-validator validate --schema-location test-resources/schema.json --template-location  test-resources/jcasc-config.yaml
```

To run the validator against a templated configfile, you need to template your config first some the validator has something to parse:

```
helm template <CHART_NAME> --output-dir generated
```

The JCASC config will be in the following location:

```
generated/<CHART_NAME>/charts/jenkins/templates/jcasc-config.yaml
```

