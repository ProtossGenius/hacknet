include('str.js');
function parse(file){ // AAAAAAAA
	it = readFile(file)
	list = it.split('\n')
	pkg = ''
	msgs = []

	desc = ''
	for (i in list){
		line = list[i]
		if (line.startWith("package")){
			pkg = line.substr(8, line.indexOf(';') - 8).trim();
			continue;
		}
		if (line.startWith('//')) {
			desc = line;
			continue;
		}

		if (line.startWith("message")) {
			msg = {}
			msg.name = line.substr(8, line.indexOf("{") - 8).trim();
			if (line.indexOf('//') != -1){
				msg.desc = line.substr(line.indexOf("//"))
			} else {
				msg.desc = desc;
			}

			desc = '';

			msgs.push(msg);
		}
	}

	return {
		pkg: pkg, 
		msgs:msgs,
	}
}

function writeTabs(tabs) {
	for (i = 0; i < tabs; i++) {
		write('\t');
	}
}

function proto2GoItf(file, itfName, itfDesc) {
	setIgnoreInput(true);
	write("// " + itfName + " " + itfDesc + ".\n")
	write("type " +  itfName + " interface {\n")
	writeln = function(str) {
		write('\t' + str + '\n');
	}
	pinfo = parse(file);
	pkg = pinfo.pkg
	msgs = pinfo.msgs

	for (i in msgs) {
		msg = msgs[i];
		if (!msg.desc.startWith('// ' + msg.name)){
			writeln("// " + msg.name + " " + msg.desc + ".")
		} else {
			writeln(msg.desc)
		}

		writeln(msg.name + "(email string, hackerAddr *net.UDPAddr, msg *" + pkg + "." + msg.name + ") (*proto.Message, map[string]interface{}, error)")
	}

	write("}\n")
}

function proto2GoSwitch(file, tabs) {
	setIgnoreInput(true);
	if (tabs == null) {
		tabs = 0;
	}
	writeln = function(str) {
		writeTabs(tabs);
		write(str + '\n');
	}
	pinfo = parse(file);
	pkg = pinfo.pkg
	msgs = pinfo.msgs
	writeln('const UnmarshalMsgMsg = "Unmarshal msg.Msg"\n')
	writeln("var _resp *proto.Message\n")
	writeln('var detail details\n')
	writeln('switch msg.Enum {');
	for (i in msgs) {
		msg = msgs[i].name
		writeln("case int32(smn_dict.EDict_" + pkg + "_" + msg + "):");
		tabs++;
		writeln("_subMsg" + " := new(" +  pkg + "." + msg +")")
		writeln("if err = proto.Unmarshal([]byte(msg.Msg), _subMsg); err != nil {")
		tabs++;
		writeln('return UnmarshalMsgMsg, details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)')
		tabs--;
		writeln("}")
		write("\n")
		writeln("if _resp, detail, err = s." + msg + "(msg.Email, hackerAddr, _subMsg); err != nil {")
		tabs++;
		writeln('return "s.' + msg + '", detail, wrapError(err)')
		tabs--;
		writeln("}\n")
		writeln("s.write(msg.Email, hackerAddr, _resp)")
		tabs--;

	}
	writeln("default:")
	tabs++;
	writeln('return "unknow Enum", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(errors.New(ErrUnexceptEnum))')
	tabs--;
	writeln("}")
	writeTabs(tabs)
}

//protoPath = "../../protos/cs.proto"
//write("\n/*");
//proto2GoSwitch(protoPath, 1)
//proto2GoItf(protoPath)
//write("*/")

//:SMIST include("parseProtos.js"); setIgnoreInput(true);
/*	switch msg.Enum {
	case int32(smn_dict.Deict_cs_Register):
		_subMsg := new(cs.Register)
		if err := proto.Unmarshal([]byte(msg.Msg), subMsg); err != nil {
			return "unmarshal msg.Msg", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err := s.Register(_subMsg); err != nil {
			return "s.Register", detail, wrapError(err)
		} else {
			s.write(hackerAddr, _resp)
		}
	case int32(smn_dict.Deict_cs_CheckEmail):
		_subMsg := new(cs.CheckEmail)
		if err := proto.Unmarshal([]byte(msg.Msg), subMsg); err != nil {
			return "unmarshal msg.Msg", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err := s.CheckEmail(_subMsg); err != nil {
			return "s.CheckEmail", detail, wrapError(err)
		} else {
			s.write(hackerAddr, _resp)
		}
	case int32(smn_dict.Deict_cs_AskHack):
		_subMsg := new(cs.AskHack)
		if err := proto.Unmarshal([]byte(msg.Msg), subMsg); err != nil {
			return "unmarshal msg.Msg", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err := s.AskHack(_subMsg); err != nil {
			return "s.AskHack", detail, wrapError(err)
		} else {
			s.write(hackerAddr, _resp)
		}
	case int32(smn_dict.Deict_cs_HeartJump):
		_subMsg := new(cs.HeartJump)
		if err := proto.Unmarshal([]byte(msg.Msg), subMsg); err != nil {
			return "unmarshal msg.Msg", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err := s.HeartJump(_subMsg); err != nil {
			return "s.HeartJump", detail, wrapError(err)
		} else {
			s.write(hackerAddr, _resp)
		}
	}
	// Register register this client to server.
	Register(*cs.Register)
	// CheckEmail check if email belong to register.
	CheckEmail(*cs.CheckEmail)
	// AskHack ask connect another client.
	AskHack(*cs.AskHack)
	// HeartJump heart jump just for keep alive.
	HeartJump(*cs.HeartJump)
*/
