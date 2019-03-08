# <center> Go-SDK </center>

<center> Go-SDK is an SDK for creating it-chain SmartContract. </center>

 

## How to create smart contract using Go-SDK

### 1. Installation Go-SDK

```bash
go get -u github.com/DE-labtory/sdk
```



### 2. Create your own Icode project

create new go project for your icode.

your icode need to implement or invoke folloing:

- Implementation go-sdk handler for handle invoke  or query
- Parse -p flag for port in main ( The port is used to create the ibox. )
- Create Ibox using sdk NewIBox
- Set handler in IBox
- Call ibox.On function at the end of main
- Contain all library using in your icode to vendor folder



sample icode : https://github.com/junbeomlee/learn-icode



### 3. Implement requires

- <a name="implementHandler"></a>Implementation go-sdk handler for handle invoke  or query

 sdk-go handler is like below

```go
type RequestHandler interface {
	Name() string
	Versions() []string
	Handle(request *pb.Request, cell *Cell) *pb.Response
}
```

your icode need to implement this interface.

`Name()`  function need to be return icode name.

`Versions()` function need to be return versions that your icode can process.

`Handle(request *pb.Request, cell *Cell) *pb.Response` function need to handle request that your icode want to handle.



- <a name="parsePort"></a>Parse -p flag for port in main ( The port is used to create the ibox. )

ICode need port for interact with it-chain engine. [tesseract](https://github.com/DE-labtory/tesseract) will give port using -p flag.

So your icode need to parse -p flag for port.

sample code is :

```go
import (
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Port int `short:"p" long:"port" description:"set port"`
}

func main() {
	logger.EnableFileLogger(true, "./icode.log")
    port := opts.Port
    ...
    ...
    ...
}
```



- Create Ibox using sdk NewIBox

IBox handle transaction excute request and interact with engine through tesseract.

You can create IBox just call `sdk.NewIBox([port])`. Port args must put [parsed data](#parsePort) 



- Set handler in IBox

After you [implement handler](#implementHandler) , you must set handler to ibox.

You can set handler using `[ibox instance].SetHandler([handler instance])`



- Call ibox.On function at the end of main

In the end of main your icode, You must call ibox On function.

IBox on function will create rpc server and start interact with engine.

You can write code like `[ibox instance].On([timeout sec])`



- Contain all library using in your icode to vendor folder

if you use any library in your icode, you must make [vendor](https://golang.org/s/go15vendor) folder.



## Sample Icode

You can refer sample icode in https://github.com/DE-labtory/sdk/blob/master/example



### Author

@[hea9549](hea9549@gmail.com)

