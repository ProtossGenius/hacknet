include('write.js')
include('str.js')
function parse() {
    it = readFile("./gkritf/command.go")
    lines = it.split("\n")
    for (i in lines) {
        line = lines[i]
        if (line.startWith("type") && line.indexOf("struct") != -1) {
            name = line.substr(4, line.indexOf("struct") - 4).trim();
            parseCmd(name);
        }
    }
}

function parseCmd(name) {
    writeln("")
    writeln("// New" + name + " GeekerMsg to " + name + ".")
    writeln("func New" + name + "(params GeekerMsg) (" + name + ", error){")
    writeln("\tcmd := " + name + "{}")
    writeln("\terr := smn_data.GetDataFromStr(string(params), &cmd)")
	writeln('\treturn cmd, err')
    writeln("}\n")
    writeln("// Message " + name + " to geeker message.")
    writeln("func (c "+name+") Message() GeekerMsg {")
    writeln("\tstr, _ := smn_data.ValToJson(c)")
    writeln('\treturn GeekerMsg("' + name + '#" + str)')
    writeln("}\n")
}
