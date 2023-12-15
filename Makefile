run:
	go run main.go

cpu-test:
	CPU_TEST=true go run main.go > res.log
	pushd ./nestest && go run nestest_diff.go > ../diff.log && popd

clean:
	rm -f *.log