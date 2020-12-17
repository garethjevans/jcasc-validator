# jcasc-validator

A simple tool to validate the generated jcasc configuration against a JCASC schema.

## To run with examples

Download/Fork this repo and build the binary.

```
make build
```

Run the application against the provided test resources

```
./build/jcasc-validator validate --schema-location test-resources/schema.json --template-location  test-resources/jcasc-config.yaml
```

