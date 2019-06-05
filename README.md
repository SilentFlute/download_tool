1. 安装`go`
2. 配置环境变量:
    用户目录下使用`vim .bash_profile`,然后粘贴`export GOPATH=/xxx/xxx/gopath`进去,路径可以依据情况修改,项目路径应该放在`/xxx/xxx/gopath/src`中
3. 接着执行一下修改`. .bash_profile`
4. 打开终端`go build`之后运行`build`出来的可执行文件即可
5. 在`localhost`相应端口下输入播放地址请求`.m3u8`的地址,并填入名字
6. 打开`localhost`相应端口下的`play.html`并输入名字即可播放