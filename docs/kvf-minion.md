# kvf-minion

## Usage
    Usage:
      kvf-minion [OPTIONS]

    Runs the Minion server. This should be run in a Kubernetes Pod on each Node of the Cluster.

    Application Options:
          --listen=  Address to listen on (default: 0.0.0.0:8080) [$KVF_LISTEN]
      -t, --token=   Use given token for api user authentication [$KVF_TOKEN]
      -v, --verbose  Turn on verbose logging
          --version  Show version

    Help Options:
      -h, --help     Show this help message
