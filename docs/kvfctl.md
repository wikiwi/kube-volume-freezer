# kvfctl

## Usage
    Usage:
      kvfctl [OPTIONS] <command>

	Command-line client for kube-volume-freezer.

    Application Options:
          --address=   Address of kvf-apiserver (default: http://localhost:8080) [$KVF_ADDRESS]
          --namespace= Namespace of Pod (default: default) [$KVF_NAMESPACE]
      -t, --token=     Use given token for api user authentication [$KVF_TOKEN]
      -v, --verbose    Turn on verbose logging

    Help Options:
      -h, --help       Show this help message

    Available commands:
      freeze   Freeze Pod Volume
      list     List Pod Volumes
      thaw     Thaw Pod Volume
      version  Show version

