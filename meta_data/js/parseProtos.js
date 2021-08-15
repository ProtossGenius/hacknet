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

		writeln(msg.name + "(email string, hackerAddr *net.UDPAddr, msg *" + pkg + "." + msg.name + ") (string, map[string]interface{}, error)")
	}

	write("}\n\n")
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
	writeln("const (")
	++tabs;
	writeln('PackResult = "packResult"')
	writeln('UnmarshalMsgMsg = "Unmarshal msg.Msg"')
	--tabs;
	writeln(")\n")
	writeln("var _resp string\n")
	writeln('var detail details\n')
	writeln('var _result *hmsg.Message\n')
	writeln('switch msg.Enum {');
	for (i in msgs) {
		msg = msgs[i].name
		writeln("case int32(smn_dict.EDict_" + pkg + "_" + msg + "):");
		tabs++;
		writeln("_subMsg" + " := new(" +  pkg + "." + msg +")")
		writeln("if err = msg.Msg.UnmarshalTo(_subMsg); err != nil {")
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
		writeln('hnlog.Info("dealPackage", details{"method": "s.' + msg + '", "_resp": _resp, "details": detail, "err": err})\n')
		
		if (msg != "ForwardMsg"){
			writeln("if _result, err = Pack_" + pkg + "_Result(msg.Email, &" + pkg + ".Result{")
			tabs++;
			writeln("Enums: " + "int32(smn_dict.EDict_" + pkg + "_" + msg + "), ")
			writeln("Info: _resp,")
			tabs--;
			writeln("}); err != nil {")
			tabs++;
			writeln('return PackResult,  details{"email": msg.Email, "_resp": _resp, "error": err}, wrapError(err)')
			tabs--;
			writeln("}\n")
			writeln("writeMsg(binder, hackerAddr, _result)")
		}
		tabs--;

	}
	writeln("default:")
	tabs++;
	writeln('return "unknow Enum", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(ErrUnexceptEnum)')
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
		pfName = "Pack_" +  pmsg 
		writeln('// ' + pfName + ' pack message ' + pmsg + ".")
		writeln("func " + pfName + '(email string, msg *' + pkg + "." + msg + ') (resp *hmsg.Message, err error) {')
		++tabs;
		writeln('var any *anypb.Any')
		writeln('if any, err = ptypes.MarshalAny(msg); err != nil {')
		++tabs;
		writeln("return nil, err")
		--tabs;
		writeln("}\n")
		writeln('return &hmsg.Message{Email: email, Enum: int32(smn_dict.EDict_' + pkg + '_' + msg + '), Msg : any}, nil')
		--tabs;
		writeln("}\n")
	}
}

