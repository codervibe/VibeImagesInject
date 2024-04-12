package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/blackhat-go/bhg/ch-13/imgInject/models"
	"github.com/blackhat-go/bhg/ch-13/imgInject/pnglib"
	"github.com/blackhat-go/bhg/ch-13/imgInject/utils"
	"github.com/spf13/pflag"
)

var (
	flags = pflag.FlagSet{SortFlags: false}
	opts  models.CmdLineOpts
	png   pnglib.MetaChunk
)

func init() {
	flags.StringVarP(&opts.Input, "input", "i", "", "原始映像文件的路径")
	flags.StringVarP(&opts.Output, "output", "o", "", "输出新映像文件的路径")
	flags.BoolVarP(&opts.Meta, "meta", "m", false, "显示实际的图像元细节")
	flags.BoolVarP(&opts.Suppress, "suppress", "s", false, "抑制块十六进制数据(可以很大)")
	flags.StringVar(&opts.Offset, "offset", "", "初始化数据注入的偏移位置")
	flags.BoolVar(&opts.Inject, "inject", false, "启用此选项以在指定的偏移位置注入数据")
	flags.StringVar(&opts.Payload, "payload", "", "有效载荷是将作为字节流读取的数据")
	flags.StringVar(&opts.Type, "type", "rNDm", "Type是要注入的Chunk头的名称")
	flags.StringVar(&opts.Key, "key", "", "有效负载的加密密钥")
	flags.BoolVar(&opts.Encode, "encode", false, "异或编码有效负载")
	flags.BoolVar(&opts.Decode, "decode", false, "异或解码有效载荷")
	flags.Lookup("type").NoOptDefVal = "rNDm"
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	if flags.NFlag() == 0 {
		flags.PrintDefaults()
		os.Exit(1)
	}
	if opts.Input == "" {
		log.Fatal("致命的 ——input标志是必需的")
	}
	if opts.Offset != "" {
		byteOffset, _ := strconv.ParseInt(opts.Offset, 0, 64)
		opts.Offset = strconv.FormatInt(byteOffset, 10)
	}
	if opts.Suppress && (opts.Meta == false) {
		log.Fatal("Fatal: The --meta flag is required when using --suppress")
	}
	if opts.Meta && (opts.Offset != "") {
		log.Fatal("Fatal: The --meta flag is mutually exclusive with --offset")
	}
	if opts.Inject && (opts.Offset == "") {
		log.Fatal("Fatal: The --offset flag is required when using --inject")
	}
	if opts.Inject && (opts.Payload == "") {
		log.Fatal("Fatal: The --payload flag is required when using --inject")
	}
	if opts.Inject && opts.Key == "" {
		fmt.Println("Warning: No key provided. Payload will not be encrypted")
	}
	if opts.Encode && opts.Key == "" {
		log.Fatal("Fatal: The --encode flag requires a --key value")
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Example Usage: %s -i in.png -o out.png --inject --offset 0x85258 --payload 1234\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example Encode Usage: %s -i in.png -o encode.png --inject --offset 0x85258 --payload 1234 --encode --key secret\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example Decode Usage: %s -i encode.png -o decode.png --offset 0x85258 --decode --key secret\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags: %s {OPTION]...\n", os.Args[0])
	flags.PrintDefaults()
	os.Exit(0)
}

func main() {
	dat, err := os.Open(opts.Input)

	defer dat.Close()
	bReader, err := utils.PreProcessImage(dat)
	if err != nil {
		log.Fatal(err)
	}
	png.ProcessImage(bReader, &opts)
}
