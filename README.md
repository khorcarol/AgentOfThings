# AgentOfThings

Repository for the Agent of Things Part IB CST project.

## Building

Builds are stored inside the `/build` folder.

Run this from the root dir:

```sh
make build
```

You can also format all `.go` files in the project by running:

```sh
make format
```

Clear out the `/build` folder by running:

```sh
make clean
```

## Building android

You need to first install the [fyne command line tool](https://docs.fyne.io/started/packaging):

```sh
go install fyne.io/fyne/v2/cmd/fyne@latest
```

You must also set the following environment variables to your Android SDK and NDK:

```sh
export ANDROID_HOME=/path/to/Android/SDK
export ANDROID_NDK_HOME=/path/to/Android/NDK
```

Then build using:

```sh
make build-android
```
