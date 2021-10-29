## deleteRedisKeys

删除redisKey

### Synopsis

批量删除redisKey,支持模糊匹配

```
deleteRedisKeys  -a server:6379 -p password -k xxxKey*   [flags]
```

### Options

```
  -a, --address string    redis 服务器地址 如:127.0.0.1:6379
  -d, --db int            redis 数据库 0
  -h, --help              help for deleteRedisKeys
  -k, --key string        需要删除的redisKey 如：xxxx:xxx*
  -l, --loglevel string   设置日志级别. 支持: debug, info, warn, error, fatal (default "info")
  -p, --password string   redis 密码
```

## 添加Icon
### 安装mod
``` 
go get github.com/akavel/rsrc
```

### Build 添加图标
deleteRedisKeys.exe.exe.manifest
```
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
<assemblyIdentity
    version="1.0.0.0"
    processorArchitecture="*"
    name="controls"
    type="win32"
></assemblyIdentity>
<dependency>
    <dependentAssembly>
        <assemblyIdentity
            type="win32"
            name="Microsoft.Windows.Common-Controls"
            version="6.0.0.0"
            processorArchitecture="*"
            publicKeyToken="6595b64144ccf1df"
            language="*"
        ></assemblyIdentity>
    </dependentAssembly>
</dependency>
</assembly>
```

### 重新编译
``` 
rsrc -manifest deleteRedisKeys.exe.exe.manifest -ico favicon.ico -o deleteRedisKeys.exe.syso # 执行生成syso文件
go build -o deleteRedisKeys.exe.exe .  
```
