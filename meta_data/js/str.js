String.prototype.startWith=function(pre){
	return this.indexOf(pre) == 0
}

String.prototype.trim=function(){
　　return this.replace(/(^\s*)|(\s*$)/g, "");
}

String.prototype.ltrim=function(){
　　return this.replace(/(^\s*)/g,"");
}

String.prototype.rtrim=function(){
　　return this.replace(/(\s*$)/g,"");

}
