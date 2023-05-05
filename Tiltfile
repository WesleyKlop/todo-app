load("ext://ko", "ko_build")
load("ext://restart_process", "docker_build_with_restart")
load("ext://helm_resource", "helm_repo", "helm_resource")

allow_k8s_contexts("docker-desktop")

helm_repo("jaegertracing", "https://jaegertracing.github.io/helm-charts")
helm_resource(
    "jaeger",
    "jaegertracing/jaeger",
    deps = ["jaegertracing"],
    flags = [
        "--set",
        "allInOne.enabled=true",
        "--set",
        "provisionDataStore.cassandra=false",
        "--set",
        "storage.type=none",
        "--set",
        "agent.enabled=false",
        "--set",
        "collector.enabled=false",
        "--set",
        "query.enabled=false",
    ],
)

local_resource(
    "todo-api-bin",
    "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/todo-api ./cmd/app",
    deps = ["./cmd", "./internal"],
)

# ko_build('todo-api', './cmd/app', deps=['./cmd', './internal'])
docker_build_with_restart(
    "todo-api-image",
    ".",
    entrypoint = ["/app/bin/todo-api"],
    only = [
        "./bin",
    ],
    live_update = [
        sync("./bin", "/app/bin"),
    ],
)
docker_build(
    "todo-app-image",
    context = "./todo-app",
    ignore = ["./todo-app/dist/", "./todo-app/node_modules/"],
    live_update = [
        fall_back_on("./todo-app/vite.config.js"),
        sync("./todo-app/", "/app/"),
        run(
            "npm ci",
            trigger = ["./todo-app/package.json", "./todo-app/package-lock.json"],
        ),
    ],
)

k8s_yaml(kustomize("manifests"))
k8s_resource("todo-api", resource_deps = ["todo-api-bin"])
