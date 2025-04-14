for f in $(ls ../keys/c/ciscodevnet/*.asc); do
	f=$(basename $f)
	GH_TOKEN=FIXME go run ./cmd/verify-gpg-key/ -key-file=../keys/c/ciscodevnet/$f -org ciscodevnet -provider-data ../providers/ -username=unknown 2>&1 | tee ${f}.log
	cat ${f}.log | grep "Key is valid for provider version" > ${f}.valid
done

