# OSS服务高并发场景下测试失败根因分析报告

-----

### 一. 问题描述

#### 名词说明：

- `OSS`：阿里云对象存储`OSS（Object Storage Service）`是阿里云提供的海量、安全、低成本、高持久的云存储服务
- `Bucket`：存储空间是您用于存储对象`Object`的容器，所有的对象都必须隶属于某个存储空间
- `Object`：对象是`OSS`存储数据的基本单元，也被称为`OSS`的文件
- `Endpoint`：表示`OSS`对外服务的访问域名
- `AccessKey`：简称AK，指的是访问身份验证中用到的`AccessKey ID`和`AccessKey Secret`

#### 现象总结：

​	`OSS`在高并发`压力测试`场景下获取和上传`oss`文件内容接口经常出现 `dial tcp x.x.x.x: connect: cannot assign requested address` 和`dial tcp x.x.x.x i/o timeout` （上述`x.x.x.x`是一个`IP`地址）现象，而在单个接口请求时基本无法捕获到，造成并发场景下和混合场景下报错率居高不下。

##### 日志截图如下：

![assign-error](/home/administrator/Downloads/images/assign-error.png)

​                                                                                                               图 1.1 **assign** 错误日志

![错误信息](/home/administrator/Downloads/images/iotimeout-log.png)

​                                                                                                               图 1.2 **io timeout** 错误日志

#### 项目环境声明：

- 系统环境：`hub.deepin.com/library/debian:latest;` 
- `golang`环境：`golang:hub.deepin.com/library/golang:latest`
- 框架环境：`gin:v1.6.3`
- `utcloud`版本环境：`master：555ff752`
- 测试部署环境：`k8s集群`

#### 流程图：

```flow
st=>start: 请求入口
buckOp=>operation: 获取Bucket
genKeyOp=>operation: 根据业务生成key
getOssOp=>operation: 根据获取OSS信息
typeCon=>condition: 是否是下载接口
getOss=>operation: 上传oss
endOp=>end: 结束框

st->buckOp->getOssOp->getOssOp->typeCon

typeCon(yes)->endOp
typeCon(no)->getOss->endOp



```

​                                                                                                              图 1.2 `api`请求流程图

#### 代码追溯：

业务代码 `ossC.go`

```go
// 获取bucket
func GetBucket(bucketName string) (*oss.Bucket, error) {
	...根据配置文件初始化bucket相关的配置，没有http请求
	client, err := oss.New(_aoi.AliyunEndpoint, _aoi.AliyunAccessKeyID, _aoi.AliyunAccessKeySecret)
	... 错误处理
	return client.Bucket(bucketName)
}                
```

`SDK`中的 `client.go`

```go
func (client Client) Bucket(bucketName string) (*Bucket, error) {
	... 初始化bucket
	return &Bucket{
		client,
		bucketName,
	}, nil
}
```

`ossC.go`中获取远程文件信息

```go
// 获取OSS 远程文件信息,为了代码更清晰忽略了错误处理
func GetMetaDetail(uid, appid, dkey, bucketName string) (*MetaDetail, error) {
	bucket, _ := GetBucket(bucketName)
	
    ... 此处是判断获取bucket是否正常

	... 此处是根据业务数据生成oss的rkey
   
	head, err := bucket.GetObjectDetailedMeta(rkey)
	if err != nil {
		zlog.Errorf("oss get object meta detail error %v", err) // 可以看到报错的位置在此处
		return nil, err
	}
	... 此处是OSS返回的对象转化为业务对象逻辑

	return &MetaDetail{
		Size:   size,
		Md5Sum: hex.EncodeToString(bmd5),
	}, nil
}
```

sdk中代码追溯

```go
// 以下已是OSS的SDK内部的方法
func (bucket Bucket) GetObjectDetailedMeta(objectKey string, options ...Option) (http.Header, error) {
	...参数处理
    resp, err := bucket.do("HEAD", objectKey, params, options, nil, nil)
    ...出参处理
}

func (bucket Bucket) do(method, objectName string, params map[string]interface{}, options []Option,
	... 参数处理
	resp, err := bucket.Client.Conn.Do(method, bucket.BucketName, objectName,
		params, headers, data, 0, listener)
	...出参处理
}

func (conn Conn) Do(method, bucketName, objectName string, params map[string]interface{}, headers map[string]string,
	... 参数处理
	return conn.doRequest(method, uri, resource, headers, data, initCRC, listener)
}

func (conn Conn) doRequest(method string, uri *url.URL, canonicalizedResource string, headers map[string]string,
	data io.Reader, initCRC uint64, listener ProgressListener) (*Response, error) {
	...参数处理
    resp, err := conn.client.Do(req)
    ... 出参处理
}
                    
// 以下已是go标准库代码 net/http
func (c *Client) Do(req *Request) (*Response, error) {
	return c.do(req)
}
```

-----

### 二. 问题分析

#### 复现测试模拟脚本：

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// 起步为150并发，每隔5分钟增加50并发，直至1500并发
func main() {
	startBench := 150
	bench := startBench
	endBench := 1500
	step := 50
	for {
		if bench == endBench {
			return
		}
		for i := 1; i < bench; i++ {
			go func() {
				doRequest()
			}()
		}
		time.Sleep(5 * time.Minute)
		bench += step
	}
}

func doRequest() {
	client := &http.Client{}
	var req *http.Request

	host := "http://127.0.0.1:9090"
	route := "/api/v0/app/meta"

	data := make(map[string]string)
	data["bin_path"] = "/data/home/xxx/go/src/deepin-sync-daemon/test/dbus/helper/bench_test"
	data["key"] = "/home/uos/Desktop/language_info.json"
	dataJson, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("生成json出错%v\n", err)
		return
	}

	req, err = http.NewRequest("PUT", host+route, bytes.NewBuffer(dataJson))
	if err != nil {
		fmt.Printf("组建http请求体报错！%v\n", err)
		return
	}
	cookie := &http.Cookie{Name: "token", Value: "40c207f994c52234a790828659e5260b"}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发起http请求报错！%v\n", err)
		return
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("http响应报错！%v\n", err)
		return
	}
	fmt.Println(string(body))
}
```



#### 问题1 :

查询TCP流程

![tcp-flow](/home/administrator/Downloads/images/tcp-flow.gif)

​                                                                                                             图 2.1 `tcp`流程图 

主动关闭连接的一方，连接会处在TIME-WAIT的状态下，需要等2MSL时间后，系统才会回收这条连接，端口才可以继续被使用。

状态：描述

```mark
CLOSED：无连接是活动的或正在进行
LISTEN：服务器在等待进入呼叫
SYN_RECV：一个连接请求已经到达，等待确认
SYN_SENT：应用已经开始，打开一个连接
ESTABLISHED：正常数据传输状态
FIN_WAIT1：应用说它已经完成
FIN_WAIT2：另一边已同意释放
ITMED_WAIT：等待所有分组死掉
CLOSING：两边同时尝试关闭
TIME_WAIT：另一边已初始化一个释放
LAST_ACK：等待所有分组死掉
```

根据代码追溯可以看到直接进入到`SDK`中去了，而且基本没有业务流程，基本上根据官方demo实现的流程。

linux机器最多有65535个端口，0~1023是系统保留的，应用程序可以使用的端口是剩下的6w多个。如果client频繁连接Server，就会导致client本地端口已经耗尽，没有端口可用，就会报出该错。

##### 猜想1：

`服务器端口链接端口数设置过低并发场景下耗尽`

##### 验证：

执行命令修改如下2个内核参数 （需要root权限）

` sysctl -w net.ipv4.tcp_timestamps=1  开启对于TCP时间戳的支持,若该项设置为0，则下面一项设置不起作用` 

`sysctl -w net.ipv4.tcp_tw_recycle=1  表示开启TCP连接中TIME-WAIT sockets的快速回收`

` sysctl -w net.ipv4.tcp_tw_reuse = 1  表示开启重用。允许将TIME-WAIT sockets重新用于新的TCP连接，默认为0，表示关闭
` 

调整配置文件 `vi /etc/sysctl.conf`

`net.ipv4.ip_local_port_range = 10000   65000   意味着10000~65000端口可用`

经测试脚本运行发现**该报错依然复现**，猜想1被否决

##### 猜想2：

`内存泄露，最可能的是httpTransport没有及时被垃圾回收，造成文件句柄不够用`

##### 验证：

复用 `Transport`, 对`GetBucket`逻辑进行如下改造

```go
var client *oss.Client

func GetClientInstance() *oss.Client {
	once.Do(func() {
		osClient, err := oss.New(_aoi.AliyunEndpoint, _aoi.AliyunAccessKeyID, _aoi.AliyunAccessKeySecret)
		... 其他参数设置
		client = osClient
	})
	return client
}

func GetBucket(bucketName string) (*oss.Bucket, error) {
	if _aoi == nil {
		zlog.Panicf("aliyun oss init config")
		return nil, errors.New("aliyun oss not init config")
	}

	bucket, err := client.Bucket(bucketName)
	return bucket, err
}
```

继续用压力测试脚本进行测试，发现该错误已经消失，问题1已验证并修复。

-----

#### 问题2：

分析错误抛出的提示:`"dial tcp x.x.x.x i/o timeout"` ，可能存在如下几种可能：

1. 该`IP`地址对应的`DNS`请求响应超时
2. 阿里云服务器防火墙或者其他安全措施拒绝高频响应
3. `SDK`调用时超时设置出错

针对上述假设的可能性进行逐一排查：

##### 猜想1：

```go
该`IP`地址对应的`DNS`请求响应超时
```

##### 验证：

由于上面的代码追溯定位可以确定，`"dial tcp x.x.x.x i/o timeout"`就是请求阿里oss触发的，故可以在并发压力测试同步进行ping操作

```go
ping oss-cn-xxx.aliyuncs.com 
```

![ping记录](/home/administrator/Downloads/images/dns-log.png)

​                                                                                                              图 2.1 `ping`请求流程图 

可以确定该`host`对应的`IP`就是日志中记录的`IP`，多次测试后结果发现在测试期间`DNS`解析正常，没有任何异常记录，证明在压力测试下`DNS`没有任何波动，猜想1可以忽略。

#### 猜想2：

```go
阿里云服务器防火墙或者其他安全措施拒绝高频响应
```

##### 验证：

如果阿里云拒绝服务，那么拿到`oss`的请求的入参和~~出参~~(不一定有)去找技术客服反馈情况，请求支持。

跟进代码内，阅读

```go
func GetMetaDetail(uid, appid, dkey, bucketName string) (*MetaDetail, error) {
	...此处是业务处理代码
	head, err := bucket.GetObjectDetailedMeta(rkey)
	if err != nil {
        ... 添加请求头X-Oss-Request-Id日志
		zlog.Errorf("oss get object meta detail error %#v request-id is %s\n", err, head.Get("X-Oss-Request-Id")) // 此处添加日志
		return nil, err
	}
	... 出参处理
}
func (bucket Bucket) GetObjectDetailedMeta(objectKey string, options ...Option) (http.Header, error) {
	...
}
```

经定位日志后拿不到`SDK`生成的请求头`X-Oss-Request-Id`，可以发现报错发生于请求前，故阿里云拒绝服务的猜测不成立。

#### 猜测3：

```markdown
SDK调用时超时设置出错
```

一.因为服务端打印的日志其实只是服务端应用层打印的日志。但客户端应用层发出数据后，中间还经过客户端的传输层，网络层，数据链路层和物理层，再经过服务端的物理层，数据链路层，网络层，传输层到服务端的应用层。假设服务端应用层处耗时100ms，再原路返回。那剩下的3s-100ms可能是耗在了整个流程里的各个层上。比如网络不好的情况下，传输层TCP使劲丢包重传之类的原因。

二.网络没问题，客户端到服务端链路整个收发流程大概耗时就是100ms左右。客户端处理逻辑问题导致超时

##### 验证：

开启`wireshark`抓包

![图片](/home/administrator/Downloads/images/wireshark.png)

分析下，从刚开始三次握手（**画了红框的地方**）。

到最后出现超时报错 `i/o timeout` （**画了蓝框的地方**）。

从`time`那一列从`7`到`10`，确实间隔`3s`。而且看**右下角**的蓝框，是`51169`端口发到`80`端口的一次`Reset`连接。

`80`端口是服务端的端口。换句话说就是客户端`3s`超时**主动**断开链接的。

但是再仔细看下**第一行**三次握手到**最后**客户端超时主动断开连接的中间，其实有**非常多次HTTP请求**。

##### 超时原因

​		假设第一次请求要`100ms`，每次请求完`http://baidu.com` 后都**放入**连接池中，下次继续复用，重复`29`次，耗时`2900ms`。

第`30`次请求的时候，连接从建立开始到服务返回前就已经用了`3000ms`，刚好到设置的**3s**超时阈值，那么此时客户端就会报超时 `i/o timeout` 。

虽然这时候服务端其实才花了`100ms`，但耐不住前面`29次`加起来的耗时已经很长。

也就是说只要通过 `http.Transport` 设置了 `err = conn.SetDeadline(time.Now().Add(time.Second * 3))`，并且你用了**长连接**，哪怕服务端处理再快，客户端设置的超时再长，总有一刻，你的程序会报超时错误。

原本预期是给每次调用设置一个超时，而不是给整个连接设置超时。

另外，上面出现问题的原因是给**长连接设置了超时，且长连接会复用**

- `http.Transport`(复用链接请求)里的建立连接时的一些超时设置干掉了。
- 在发起`http`请求的时候会场景`http.Client`（每次调用客户端），此时加入超时设置，这里的超时就可以理解为单次请求的超时了。同样可以看下注释

`SetDeadline sets the read and write deadlines associated with the connection.`

##### 代码调整

```go
func GetClientInstance() *oss.Client {
	once.Do(func() {
		osClient, err := oss.New(_aoi.AliyunEndpoint, _aoi.AliyunAccessKeyID, _aoi.AliyunAccessKeySecret, func(c *oss.Client) {
			t := http.DefaultTransport.(*http.Transport)
			t.MaxConnsPerHost = 1500 //测试配置
			t.MaxIdleConns = 100 //测试配置
			t.MaxIdleConnsPerHost = 100 //测试配置
			c.HTTPClient = http.DefaultClient
		})
		...其他逻辑
		client = osClient
	})
	return client
}
```

验证：

​	压力测试下，问题2已不复现，问题验证并修复。

-----

### 三.  解决办法

1. 复用请求链接（问题1）

2. 设置单个请求超时时间（问题2）

-----

### 四、总结

1. `http`请求要复用`http.Transport`，`SDK`案例不可全信
2. 超时设置在`http.Client`里，为每个请求设置超时，不要在 `http.Transport`中设置超时，那是连接的超时，不是请求的超时。否则可能会出现莫名 `io timeout`报错。


### 参考文档：

 - 阿里云上传`demo`：[https://help.aliyun.com/document_detail/88620.html](https://help.aliyun.com/document_detail/88620.html)
 - Cannot assign requested address问题总结 [https://www.jianshu.com/p/51a953b789a4](https://www.jianshu.com/p/51a953b789a4)
 - `go`与`TCP`常见的问题: [https://tonybai.com/2015/11/17/tcp-programming-in-golang/](https://tonybai.com/2015/11/17/tcp-programming-in-golang/)

