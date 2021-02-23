# cs-cli

## Configuration

The configuration file stores the targets (vRA servers and credentials) that the CLI will use. By default cs-cli will use `$HOME/.cs-cli.yaml` as the config file. You can override this using the `--config` flag.

```bash
cs-cli --config /path/to/config.yaml get pipeline
```

### Working with targets

List available targets:
`cs-cli config get-target`

Add an target configuration:
```bash
❯ cs-cli config set-target --config test-config.yaml --name my-vra-server --password mypassword --username myuser --domain mydomain.com --server my-vra-server.mydomain.com
Creating new target my-vra-server
Use `cs-cli config use-target --name my-vra-server` to use this target
{
  "domain": "mydomain.com",
  "password": "mypassword",
  "server": "my-vra-server.mydomain.com",
  "username": "myuser"
}
```

```bash
#Set the active target
cs-cli config use-target --name my-vra-server --config test-config.yaml
#View the current active target
cs-cli config current-target --config test-config.yaml
```

## Working with Pipelines

Getting and filtering pipelines
```bash
# List all pipelines
cs-cli get pipeline
# List all pipelines in a project
cs-cli get pipeline --project "Field Demo"
# Get a pipeline by ID
cs-cli get pipeline --id 7a3b41af-0e49-4e3d-999b-6c4c5ec55956
# Get a pipeline by name
cs-cli get pipeline --name "vra-CreateVariable"
```

Exporting pipelines:
```bash
# Export a specific pipeline to current location
cs-cli get pipeline --name "vra-CreateVariable" --export
# Export a specific pipeline to a specific location
cs-cli get pipeline --name "vra-CreateVariable" --export --exportPath /path/to/my/folder
# Export all pipelines
cs-cli get pipeline --export
# Export all pipelines in a project
cs-cli get pipeline --project "Field Demo" --export
```

## Working with Variables

```bash
# Get all variables
cs-cli get variable
# Get a variable by ID
cs-cli get variable --id 50613ab6-6f25-4976-8b3e-5be7a4bc60eb
# Get a variable by name
cs-cli get variable --name cs-cli
# Create a new variable manually
cs-cli create variable --name cli-demo --project "Field Demo"  --type REGULAR --value "New variable..." --description "Now from the CLI\!"
```

# Working with Executions

```bash
# List all executions
cs-cli get execution
# View an execution by ID
cs-cli get execution --id 9cc5aedc-db48-4c02-a5e4-086de3160dc0
# View executions of a specific pipeline
get execution --name vra-authenticateUser
# View executions by status
cs-cli get execution --status Failed
```

Create a new execution of a pipeline:
```bash
# Get the input form of the pipeline to execute
cs-cli get pipeline --id 7a3b41af-0e49-4e3d-999b-6c4c5ec55956 --form
# Outputs:
# {
#   "vraFQDN": "",
#   "vraUserName": "",
#   "vraUserPassword": ""
# }

# Create a new execution with the input form from above
cs-cli create execution --id 7a3b41af-0e49-4e3d-999b-6c4c5ec55956 --inputs '{
  "vraFQDN": "vra8-test-ga.cmbu.local",
  "vraUserName": "fakeuser",
  "vraUserPassword": "fakeuser"
}' --comments "Executed from cs-cli"
# Outputs
# Execution /codestream/api/executions/9cc5aedc-db48-4c02-a5e4-086de3160dc0 created

# Inspect the new execution
cs-cli get execution --id 9cc5aedc-db48-4c02-a5e4-086de3160dc0
```



## Working with Endpoints
\\