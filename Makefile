run:
	go run main.go

cpu-test:
	CPU_TEST=true go run main.go nestest/nestest.nes > res.log || true
	pushd ./nestest && go run nestest_diff.go > ../diff.log && popd

clean:
	rm -f *.log