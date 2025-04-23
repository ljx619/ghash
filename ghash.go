package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
)

var (
	algorithm = flag.String("a", "sha256", "Hash algorithm (md5, sha1, sha256, sha512)")
	filePath  = flag.String("f", "", "Input file path (optional)")
)

var hashes = map[string]func() hash.Hash{
	"md5":    md5.New,
	"sha1":   sha1.New,
	"sha256": sha256.New,
	"sha512": sha512.New,
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  hasher -a md5 -f file.txt")
		fmt.Fprintln(os.Stderr, "  echo 'hello' | hasher -a sha1")
	}
	flag.Parse()

	// 验证算法是否支持
	newHash, ok := hashes[*algorithm]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unsupported algorithm: %s\n", *algorithm)
		fmt.Fprintln(os.Stderr, "Supported algorithms:")
		for k := range hashes {
			fmt.Fprintln(os.Stderr, "  -", k)
		}
		os.Exit(1)
	}

	// 选择输入源
	var input io.Reader = os.Stdin
	if *filePath != "" {
		f, err := os.Open(*filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		input = f
	}

	// 计算哈希
	h := newHash()
	if _, err := io.Copy(h, input); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// 输出结果
	fmt.Printf("%x\n", h.Sum(nil))
}
