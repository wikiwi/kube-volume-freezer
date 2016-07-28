# kvf-apiserver

## Usage
    Usage:
      kvf-apiserver [OPTIONS]

    Runs the API server. Delegates freeze, and thaw requests to the Minion containing the Pod.

    Application Options:
          --listen=           Address to listen on (default: 0.0.0.0:8080) [$KVF_LISTEN]
      -t, --token=            Use given token for api user authentication [$KVF_TOKEN]
          --minion-token=     Use given token to authenticate to minion servers [$KVF_MINION_TOKEN]
          --minion-selector=  K8s label selector to find the Minion Pods [$KVF_MINION_SELECTOR]
          --minion-namespace= Namespace of Minion Pods (default: default) [$KVF_MINION_NAMESPACE]
          --minion-port=      Port of Minion Pods (default: 8080) [$KVF_MINION_PORT]
      -v, --verbose           Turn on verbose logging
          --version           Show version

    Help Options:
      -h, --help              Show this help message

