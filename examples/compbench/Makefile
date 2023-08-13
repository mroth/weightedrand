.PHONY: analysis clean
FILE=results.txt
COUNT=10

analysis: $(FILE)
	benchstat -filter=".name:Multiple" -col="/lib" -row="/size" -table="/concurrency" $<

clean:
	rm -f $(FILE)

$(FILE):
	go test -bench=. -count=${COUNT} | tee $@
