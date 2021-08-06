function writeln(str) {
	write(str + "\n")
}

function createConsts(constName, constDesc, constStatus) {
	writeln("// " + constName + " " + constDesc + ".");
	writeln("type " + constName + " int\n");
	writeln("// " + constName + "\'s members.");
	first = true
	writeln("const (")
	for (key in constStatus) {
		writeln("\t// " + constStatus[key] + ".");
		write("\t" + key);
		if (first) {
			write(" = iota");
			first = false;
		}
		write("\n");
	}
	writeln(")\n")

	writeln("// Get" + constName + "Desc get "+constName+"'s desc.")
	writeln("func Get" + constName + "Desc(constStatus " + constName + ") string {")
	writeln("\tswitch constStatus {")
	for (key in constStatus) {
		writeln("\tcase " + key + ":")
		writeln("\t\treturn \"" + constStatus[key] + "\"")
	}
	writeln("\tdefault:")
	writeln("\t\treturn \"?UNKNOW?\"")
	writeln("\t}")
	writeln("}\n")

	writeln("// Get" + constName + "Name get "+constName+"'s desc.")
	writeln("func Get" + constName + "Name(constStatus " + constName + ") string {")
	writeln("\tswitch constStatus {")
	for (key in constStatus) {
		writeln("\tcase " + key + ":")
		writeln("\t\treturn \"" + key + "\"")
	}
	writeln("\tdefault:")
	writeln("\t\treturn \"?UNKNOW?\"")
	writeln("\t}")
	writeln("}")
}

