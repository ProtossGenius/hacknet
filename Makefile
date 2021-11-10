##Tail
prebuild:
	# compile proto.
	rm -rf ./pb
	mkdir pb
	smnrpc-autocode  -cfg ./datas/ric/client_ric.json
#	smn_protocpl -i protos -gm github.com/ProtossGenius/hacknet -lang go -o /tmp
	mv /tmp/github.com/ProtossGenius/hacknet/pb/* ./pb
	smist -jsPath='./meta_data/js'
debug:

qrun: prebuild
test:

install:

clean:
