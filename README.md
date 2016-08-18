# Springbeat

Welcome to Springbeat.

**Important Notes:** 
 1. For now, only two endpoints are supported, namely `/metrics` and `/health`. We'll add [more endpoints](http://docs.spring.io/spring-boot/docs/1.4.0.RELEASE/reference/htmlsingle/#production-ready-endpoints) as we go
 2. This plugin will only work if your Spring Boot application has the `spring-boot-starter-actuator` dependency

```
		<dependency>
			<groupId>org.springframework.boot</groupId>
			<artifactId>spring-boot-starter-actuator</artifactId>
		</dependency>
```

Ensure that this folder is at the following location:
`${GOPATH}/github.com/consulthys`

## Getting Started with Springbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.6
* [Glide](https://github.com/Masterminds/glide) >= 0.10.0

### Init Project
To get running with Springbeat, run the following command:

```
make init
```

To commit the first version before you modify it, run:

```
make commit
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Springbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/consulthys/springbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Springbeat run the command below. This will generate a binary
in the same directory with the name springbeat.

```
make
```


### Run

To run Springbeat with debugging output enabled, run:

```
./springbeat -c springbeat.yml -e -d "*"
```


### Test

To test Springbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`


### Package

To be able to package Springbeat the requirements are as follows:

 * [Docker Environment](https://docs.docker.com/engine/installation/) >= 1.10
 * $GOPATH/bin must be part of $PATH: `export PATH=${PATH}:${GOPATH}/bin`

To cross-compile and package Springbeat for all supported platforms, run the following commands:

```
cd dev-tools/packer
make deps
make images
make
```

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/springbeat.template.json and etc/springbeat.asciidoc

```
make update
```


### Cleanup

To clean  Springbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Springbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/consulthys
cd ${GOPATH}/github.com/consulthys
git clone https://github.com/consulthys/springbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).
