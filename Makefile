##Tail
prebuild:
	# compile proto.
	rm -rf ./pb/*
	smn_protocpl -i protos -gm github.com/ProtossGenius/hacknet -lang go -o /tmp
	mv /tmp/github.com/ProtossGenius/hacknet/pb/* ./pb
debug:

qrun:

test:

install:

clean:

