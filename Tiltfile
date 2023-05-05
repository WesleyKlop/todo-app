load('ext://ko', 'ko_build')
load('ext://restart_process', 'docker_build_with_restart')

allow_k8s_contexts('docker-desktop')

local_resource(
  'todo-api-bin',
  'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/todo-api ./cmd/app',
  deps=['./cmd', './internal'])

# ko_build('todo-api', './cmd/app', deps=['./cmd', './internal'])
docker_build_with_restart(
  'todo-api-image',
  '.',
  entrypoint=['/app/bin/todo-api'],
  only=[
    './bin'
  ],
  live_update=[
    sync('./bin', '/app/bin')
  ]
)

k8s_yaml(kustomize('manifests/local'))
k8s_resource('todo-api', resource_deps=['todo-api-bin'])
