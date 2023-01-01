go build -o gopoly -v ../../main.go
chmod +x ./gopoly
./gopoly -c testdata/.gopoly.yaml
