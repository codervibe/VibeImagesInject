# VibeImagesInject
* 图片 隐写术 
* 可以将 一定的内容写入到 图片中
~~~cmd
go run ./main.go
  -i, --input string           原始映像文件的路径
  -o, --output string          输出新映像文件的路径
  -m, --meta                   显示实际的图像元细节
  -s, --suppress               抑制块十六进制数据(可以很大)
      --offset string          初始化数据注入的偏移位置
      --inject                 启用此选项以在指定的偏移位置注入数据
      --payload string         有效载荷是将作为字节流读取的数据
      --type string[="rNDm"]   Type是要注入的Chunk头的名称 (default "rNDm")
      --key string             有效负载的加密密钥
      --encode                 异或编码有效负载
      --decode                 异或解码有效载荷
~~~

## 执行示例
~~~shell
go run ./main.go -i .\battlecat.png -o .\battlecatpayload.png --inject --offset 0x85258 --payload helloworld
~~~


