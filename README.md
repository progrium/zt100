# zt100

## Dependencies
* go 1.15+
* make
* Docker (optional)
* npm (optional, for building tailwind)

## Setup
Run this upon cloning:
```
$ make setup
```
Use `make help` for more make commands.

## Development
Run this to startup for development:
```
$ make dev
```

## Deployment
Run this to create a deployable Docker image as `okta/zt100`:
```
$ make image
```
You can easily run this image locally with `make docker`. 

### Storage

A volume mounted at `/zt100` in the container will persist state.

## Integration

* handoff (make demo, delegate auth)
* cmd endpoint
* demo endpoint