load('ext://ko', 'ko_build')

allow_k8s_contexts('docker-desktop')

ko_build('todo-api', './cmd/app', deps=['./cmd', './internal'] )

k8s_yaml(kustomize('manifests/local'))
