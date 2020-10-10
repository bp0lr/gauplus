
# GauPlus

A modified version of (http://wwww.github.com/lc/gau) for personal usage.
Support workers, proxies and some extra things.

## Usage:
Examples:

```bash
$ echo "example.com" | gauplus
$ cat domains.txt | gauplus
$ gauplus example.com
$ gauplus -o example-urls.txt example.com
$ echo "example.com" | gauplus -p "http://proxy.packetstream.io:31112" --random-agent -o result.txt -t 25
```

To display the help for the tool use the `-h` flag:

```bash
$ gauplus -h

 -json
        write output as json
  -o string
        filename to write results to
  -p string
        use proxy
  -providers string
        providers to fetch urls for (default "wayback,otx,commoncrawl")
  -random-agent
        use random user-agent
  -retries uint
        amount of retries for http client (default 5)
  -subs
        include subdomains of target domain
  -t int
        amount of parallel workers (default 5)
  -v    enable verbose mode
  -version
        show gauplus version

```

### comparison
```
[root@DarkStar]─[/opt/bp0/lovan/gau] wc -l targets.txt
31 targets.txt

[root@DarkStar]─[/opt/bp0/lovan/gau] time cat targets.txt | gau
real    7m17.529s
user    0m0.360s
sys     0m0.345s

[root@DarkStar]─[/opt/bp0/lovan/gauplus] time cat targets.txt | gauplus -p "http://proxy.packetstream.io:31112" --random-agent -t 25
real    0m49.899s
user    0m0.380s
sys     0m0.408s
```

## Installation:
### From source:
```
$ GO111MODULE=on go get -u -v github.com/bp0lr/gauplus
```

## Useful?, buy lc some coffe!

<a href="http://buymeacoff.ee/cdl" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>
