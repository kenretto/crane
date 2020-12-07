## Crane 

[![Build Status](https://travis-ci.org/kenretto/crane.svg?branch=master&status=created)](https://github.com/kenretto/crane)

Configurable integration tool


#### Description
Crane It is not a framework, but a layer of three-party libraries that may be needed for daily use, so that some libraries can support automatic reloading by modifying the configuration, making web development more convenient, instead of integrating various libraries before development Tools, bother

#### Installation
* go get
```
go get -u github.com/kenretto/crane
```

* git submodule
  * git submodule add
    ```
    git submodule add https://github.com/kenretto/crane.git 
    ```
  * go.mod replace (open go.mod)
    ```
    replace (
    	github.com/kenretto/crane version hash => your submodule crane path
    )
    ```
    
#### Examples
* see [https://github.com/kenretto/crane/tree/master/example](https://github.com/kenretto/crane/tree/master/example)

#### Usage
If you are using the method build in the example, launch the `Crane.Run` method start, 
Then, using `buildout-binary start` will make the service separate from the parent process that started it, and run in the mode of daemon,
If you run with `buildout-binary start --daemon=false`, the service will remain in the current session(this will facilitate debugging at development time, such as using GoLand)