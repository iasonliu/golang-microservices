# Product Images

## Uploading

Note: need to use `--data-binary` to ensure file is not converted to text

```
curl -vv localhost:9090/1/go.mod -X PUT --data-binary @test.png
```

## with Gzip
```
curl -v localhost:9091/images/1/aa.png --compressed -o test1.png
```