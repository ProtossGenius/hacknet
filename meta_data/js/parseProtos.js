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

		writeln(msg.name + "(email string, hackerAddr *net.UDPAddr, msg *" + pkg + "." + msg.name + ") (*hmsg.Message, map[string]interface{}, error)")
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
	writeln("var _resp *hmsg.Message\n")
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
		writeln("writeMsg(binder, hackerAddr, _resp)")
		tabs--;

	}
	writeln("default:")
	tabs++;
	writeln('return "unknow Enum", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(errors.New(ErrUnexceptEnum))')
	tabs--;
	writeln("}")
	writeTabs(tabs)
}

function packMsgs(file){
	pinfo = parse(file);
	pkg = pinfo.pkg;
	msgs = pinfo.msgs;
	tabs = 0;
	writeln = function(str) {
		writeTabs(tabs);
		write(str + '\n');
	}

	for (i in msgs) {
		msg = msgs[i].name;
		pmsg = pkg + "_" + msg;
		pfName = "pack_" +  pmsg 
		writeln('// ' + pfName + ' pack message ' + pmsg + ".")
		writeln("func " + pfName + '(email string, msg *' + pkg + "." + msg + ') (resp *hmsg.Message, err error) {')
		++tabs;
		writeln("var data []byte\n")
		writeln("if data, err = proto.Marshal(msg); err != nil {")
		++tabs;
		writeln("return nil, err")
		--tabs;
		writeln("}\n")
		writeln('return &hmsg.Message{Email: email, Enum: int32(smn_dict.EDict_' + pkg + '_' + msg + '), Msg : string(data)}, nil')
		--tabs;
		writeln("}\n")
	}
}

