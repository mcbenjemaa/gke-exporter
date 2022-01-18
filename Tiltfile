# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')


IMG='exporter:latest'

DOCKERFILE = '''FROM golang:alpine
    WORKDIR /
    COPY ./bin/exporter /
    CMD ["/exporter"]
    '''

def binary():
    return 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o bin/exporter main.go'


DIRNAME = os.path.basename(os. getcwd())
if os.path.exists('go.mod') == False:
    local("go mod init %s" % DIRNAME)

deps = ['collector', 'main.go']
deps.append('pkg')

local_resource('Watch&Compile', binary(), deps=deps)

k8s_yaml('kubernetes/manifest.yaml')

docker_build_with_restart(IMG, '.', 
     dockerfile_contents=DOCKERFILE,
     entrypoint='/exporter --project-id=training-300214',
     only=['./bin/exporter'],
     live_update=[
           sync('./bin/exporter', '/exporter'),
       ]
)




# k8s_resource('example-go', port_forwards=8000,
#              resource_deps=['deploy', 'example-go-compile'])