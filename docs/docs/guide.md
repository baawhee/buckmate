## Installation

Build from source:
    
```
git clone https://github.com/afternun/buckmate.git && go build -C buckmate && mv buckmate/buckmate /usr/local/bin
``` 

or download latest release from GitHub:


=== "MacOS ARM"

    ```
    curl https://github.com/afternun/buckmate/releases/latest/download/buckmate_Darwin_arm64.tar.gz -LO && tar -xzvf buckmate_Darwin_arm64.tar.gz && mv buckmate /usr/local/bin
    ```
    
    === "MacOS x86_64"

    ```
    curl https://github.com/afternun/buckmate/releases/latest/download/buckmate_Darwin_x86_64.tar.gz -LO && tar -xzvf buckmate_Darwin_x86_64.tar.gz && mv buckmate /usr/local/bin
    ```

=== "Linux ARM"

    ```
    curl https://github.com/afternun/buckmate/releases/latest/download/buckmate_Linux_arm64.tar.gz -LO && tar -xzvf buckmate_Linux_arm64.tar.gz && mv buckmate /usr/local/bin
    ```
    
    === "Linux i386"

    ```
    curl https://github.com/afternun/buckmate/releases/latest/download/buckmate_Linux_i386.tar.gz -LO && tar -xzvf buckmate_Linux_i386.tar.gz && mv buckmate /usr/local/bin
    ```
    
    === "Linux x86_64"

    ```
    curl https://github.com/afternun/buckmate/releases/latest/download/buckmate_Linux_x86_64.tar.gz -LO && tar -xzvf buckmate_Linux_x86_64.tar.gz && mv buckmate /usr/local/bin
    ```

=== "Windows ARM"

    ```
    curl.exe https://github.com/afternun/buckmate/releases/latest/download/buckmate_Windows_arm64.tar.gz -LO
    ```
    
    === "Windows i386"

    ```
    curl.exe https://github.com/afternun/buckmate/releases/latest/download/buckmate_Windows_i386.tar.gz -LO
    ```
    
    === "Windows x86_64"

    ```
    curl.exe https://github.com/afternun/buckmate/releases/latest/download/buckmate_Windows_x86_64.tar.gz -LO
    ```

    Append directory containing above binary to PATH environment variable


[Browse releases here](https://github.com/afternun/buckmate/releases)

## Initial setup

!!! note "Configure AWS credentials"

    If you are deploying to or from AWS S3 bucket configure AWS credentials according to their instructions.

!!! note "Examples"

    Take a look at `e2e/tests` directory in the code repository

   In the directory create:

* `Deployment.yaml` - here you define common configuration for your deployment, one that is shared across any environment that you work on
    
    ```
    source:
      address: location from which files should be copied 
      (use `s3://` prefix for s3 buckets,
       absolute path for files on disk,
       or path relative to location of this file)
    target:
      address: location to which files should be copied
      (use `s3://` prefix for s3 buckets,
       absolute path for files on disk,
       or path relative to location of this file)
    configBoundary: string that acts as prefix and suffix for config map values (Default %%%)
    configMap:
      string key: string value
    ```

!!! note "Config Map"

    **buckmate** will go over files downloaded from `source` and files defined in `files` directory and look for strings that are wrapped in `configBoundary`. If such string is found, it will be replaced with corresponding value from `configMap`. 
    
    Example: If a file would contain string `%%%header%%%` and `configMap` an entry `header: My Awesome Header`, string `%%%header%%%` would be replaced with `My Awesome Header`. 

    **Environment specific configuration takes precedence over common configuration**

* (Optional): `files` directory

    This can hold any files that will be copied alongside files downloaded from `source`

* (Optional): directory with name of your choosing with another `Deployment.yaml` and `files` directory

    This can hold environment specific configuration. To use it, run **buckmate** with `--env` flag

## Run

```
  buckmate apply
```

!!! note "Versioning"

    **buckmate** will add metadata to S3 objects `buckmate-version` with UUID string as value.
    This is used to differentiate between previous and new deployment. Files that do not match new version will be removed on the deployment.