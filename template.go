package heldiamgo

//使用js的artTemplate做模版解析引擎

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/astaxie/beego/utils"
	"github.com/robertkrimen/otto"
)

//templateJS artTemplate模版引擎
var templateJS = `!(function(){var template=function(filename,content){return typeof content==="string"?compile(content,{filename:filename}):renderFile(filename,content)};template.version="3.0.0";template.config=function(name,value){defaults[name]=value};var defaults=template.defaults={openTag:"<%",closeTag:"%>",escape:true,cache:true,compress:false,parser:null};var cacheStore=template.cache={};template.render=function(source,options){return compile(source,options)};var renderFile=template.renderFile=function(filename,data){var fn=template.get(filename)||showDebugInfo({filename:filename,name:"Render Error",message:"Template not found"});return data?fn(data):fn};template.get=function(filename){var cache;if(cacheStore[filename]){cache=cacheStore[filename]}else{if(typeof document==="object"){var elem=document.getElementById(filename);if(elem){var source=(elem.value||elem.innerHTML).replace(/^\s*|\s*$/g,"");cache=compile(source,{filename:filename})}}}return cache};var toString=function(value,type){if(typeof value!=="string"){type=typeof value;if(type==="number"){value+=""}else{if(type==="function"){value=toString(value.call(value))}else{value=""}}}return value};var escapeMap={"<":"&#60;",">":"&#62;",'"':"&#34;","'":"&#39;","&":"&#38;"};var escapeFn=function(s){return escapeMap[s]};var escapeHTML=function(content){return toString(content).replace(/&( ?![\w#]+;)|[<>"']/g,escapeFn)};var isArray=Array.isArray||function(obj){return({}).toString.call(obj)==="[object Array]"};var each=function(data,callback){var i,len;if(isArray(data)){for(i=0,len=data.length;i<len;i++){callback.call(data,data[i],i,data)}}else{for(i in data){callback.call(data,data[i],i)}}};var utils=template.utils={$helpers:{},$include:renderFile,$string:toString,$escape:escapeHTML,$each:each};template.helper=function(name,helper){helpers[name]=helper};var helpers=template.helpers=utils.$helpers;template.onerror=function(e){var message="Template Error\n\n";for(var name in e){message+="<"+name+">\n"+e[name]+"\n\n"}if(typeof console==="object"){console.error(message)}};var showDebugInfo=function(e){template.onerror(e);return function(){return"{Template Error}"}};var compile=template.compile=function(source,options){options=options||{};for(var name in defaults){if(options[name]===undefined){options[name]=defaults[name]}}var filename=options.filename;try{var Render=compiler(source,options)}catch(e){e.filename=filename||"anonymous";e.name="Syntax Error";return showDebugInfo(e)}function render(data){try{return new Render(data,filename)+""}catch(e){if(!options.debug){options.debug=true;return compile(source,options)(data)}return showDebugInfo(e)()}}render.prototype=Render.prototype;render.toString=function(){return Render.toString()};if(filename&&options.cache){cacheStore[filename]=render}return render};var forEach=utils.$each;var KEYWORDS="break,case,catch,continue,debugger,default,delete,do,else,false"+",finally,for,function,if,in,instanceof,new,null,return,switch,this"+",throw,true,try,typeof,var,void,while,with"+",abstract,boolean,byte,char,class,const,double,enum,export,extends"+",final,float,goto,implements,import,int,interface,long,native"+",package,private,protected,public,short,static,super,synchronized"+",throws,transient,volatile"+",arguments,let,yield"+",undefined";var REMOVE_RE=/\/\*[\w\W]*?\*\/|\/\/[^\n]*\n|\/\/[^\n]*$|"(?:[^"\\]|\\[\w\W])*"|'(?:[^'\\]|\\[\w\W])*'|\s*\.\s*[$\w\.]+/g;var SPLIT_RE=/[^\w$]+/g;var KEYWORDS_RE=new RegExp(["\\b"+KEYWORDS.replace(/,/g,"\\b|\\b")+"\\b"].join("|"),"g");var NUMBER_RE=/^\d[^,]*|,\d[^,]*/g;var BOUNDARY_RE=/^,+|,+$/g;var SPLIT2_RE=/^$|,+/;function getVariable(code){return code.replace(REMOVE_RE,"").replace(SPLIT_RE,",").replace(KEYWORDS_RE,"").replace(NUMBER_RE,"").replace(BOUNDARY_RE,"").split(SPLIT2_RE)}function stringify(code){return"'"+code.replace(/('|\\)/g,"\\$1").replace(/\r/g,"\\r").replace(/\n/g,"\\n")+"'"}function compiler(source,options){var debug=options.debug;var openTag=options.openTag;var closeTag=options.closeTag;var parser=options.parser;var compress=options.compress;var escape=options.escape;var line=1;var uniq={$data:1,$filename:1,$utils:1,$helpers:1,$out:1,$line:1};var isNewEngine="".trim;var replaces=isNewEngine?["$out='';","$out+=",";","$out"]:["$out=[];","$out.push(",");","$out.join('')"];var concat=isNewEngine?"$out+=text;return $out;":"$out.push(text);";var print="function(){"+"var text=''.concat.apply('',arguments);"+concat+"}";var include="function(filename,data){"+"data=data||$data;"+"var text=$utils.$include(filename,data,$filename);"+concat+"}";var headerCode="'use strict';"+"var $utils=this,$helpers=$utils.$helpers,"+(debug?"$line=0,":"");var mainCode=replaces[0];var footerCode="return new String("+replaces[3]+");";forEach(source.split(openTag),function(code){code=code.split(closeTag);var $0=code[0];var $1=code[1];if(code.length===1){mainCode+=html($0)}else{mainCode+=logic($0);if($1){mainCode+=html($1)}}});var code=headerCode+mainCode+footerCode;if(debug){code="try{"+code+"}catch(e){"+"throw {"+"filename:$filename,"+"name:'Render Error',"+"message:e.message,"+"line:$line,"+"source:"+stringify(source)+".split(/\\n/)[$line-1].replace(/^\\s+/,'')"+"};"+"}"}try{var Render=new Function("$data","$filename",code);Render.prototype=utils;return Render}catch(e){e.temp="function anonymous($data,$filename) {"+code+"}";throw e}function html(code){line+=code.split(/\n/).length-1;if(compress){code=code.replace(/\s+/g," ").replace(/<!--[\w\W]*?-->/g,"")}if(code){code=replaces[1]+stringify(code)+replaces[2]+"\n"}return code}function logic(code){var thisLine=line;if(parser){code=parser(code,options)}else{if(debug){code=code.replace(/\n/g,function(){line++;return"$line="+line+";"})}}if(code.indexOf("=")===0){var escapeSyntax=escape&&!/^=[=#]/.test(code);code=code.replace(/^=[=#]?|[\s;]*$/g,"");if(escapeSyntax){var name=code.replace(/\s*\([^\)]+\)/,"");if(!utils[name]&&!/^(include|print)$/.test(name)){code="$escape("+code+")"}}else{code="$string("+code+")"}code=replaces[1]+code+replaces[2]}if(debug){code="$line="+thisLine+";"+code}forEach(getVariable(code),function(name){if(!name||uniq[name]){return}var value;if(name==="print"){value=print}else{if(name==="include"){value=include}else{if(utils[name]){value="$utils."+name}else{if(helpers[name]){value="$helpers."+name}else{value="$data."+name}}}}headerCode+=name+"="+value+",";uniq[name]=true});return code+"\n"}}defaults.openTag="{{";defaults.closeTag="}}";var filtered=function(js,filter){var parts=filter.split(":");var name=parts.shift();var args=parts.join(":")||"";if(args){args=", "+args}return"$helpers."+name+"("+js+args+")"};defaults.parser=function(code,options){code=code.replace(/^\s/,"");var arr=code.split(" ");var split=[];for(k in arr){split[k]=arr[k]}var key=split.shift();var args=split.join(" ");switch(key){case"if":code="if("+args+"){";break;case"else":if(split.shift()==="if"){split=" if("+split.join(" ")+")"}else{split=""}code="}else"+split+"{";break;case"/if":code="}";break;case"each":var object=split[0]||"$data";var as=split[1]||"as";var value=split[2]||"$value";var index=split[3]||"$index";var param=value+","+index;if(as!=="as"){object="[]"}code="$each("+object+",function("+param+"){";break;case"/each":code="});";break;case"echo":code="print("+args+");";break;case"print":case"include":code=key+"("+split.join(",")+");";break;default:if(/^\s*\|\s*[\w\$]/.test(args)){var escape=true;if(code.indexOf("#")===0){code=code.substr(1);escape=false}var i=0;var array=code.split("|");var len=array.length;var val=array[i++];for(;i<len;i++){val=filtered(val,array[i])}code=(escape?"=":"=#")+val}else{if(template.helpers[key]){code="=#"+key+"("+split.join(",")+");"}else{code="="+code}}break}return code};if(typeof define==="function"){define(function(){return template})}else{if(typeof exports!=="undefined"){module.exports=template}else{this.template=template}}})();template.utils.$include=includefunc;template;`

//JSTemplate JS模版引擎
type JSTemplate struct {
	engine         *otto.Object      //js解释器对象
	templateDIR    string            //模版文件路径
	ExtendType     string            //扩展文件名
	Cache          bool              //是否缓存
	templateCaches map[string][]byte //模版解析缓存
	err            error             //错误
}

//NewJSTemplate 深层JS模版引擎对象
func NewJSTemplate(TemplateDIR string, ExtendType ...string) *JSTemplate {
	result := &JSTemplate{
		Cache: true,
	}
	result.SetTemplateDir(TemplateDIR)
	if ExtendType != nil && len(ExtendType) > 0 {
		result.ExtendType = ExtendType[0]
	}
	result.init()
	return result
}

func (je *JSTemplate) init() {
	ot := otto.New()
	ot.Set("includefunc", je.includeTemplate)
	je.engine, je.err = ot.Object(templateJS)
	return
}

//SetTemplateDir 设置模版目录
func (je *JSTemplate) SetTemplateDir(TemplateDir string) {
	je.templateDIR = TemplateDir
}

//Error 错误记录
func (je *JSTemplate) Error() error {
	return je.err
}

//Template 解析模版
//
//@Param body 模版内容
//
//@Param params 模版参数
func (je *JSTemplate) Template(body []byte, params ...interface{}) ([]byte, error) {
	if je.engine == nil {
		je.init()
	}
	if je.err != nil {
		return nil, je.err
	}
	value, err := je.processTemplate(body, params...)
	je.err = err
	if err != nil {
		return nil, err
	}
	return []byte(value.String()), nil
}

func (je *JSTemplate) processTemplate(body []byte, params ...interface{}) (value otto.Value, err error) {
	value, err = je.engine.Call("compile", string(body))
	if err != nil {
		return
	}
	if params != nil && len(params) > 0 {
		value, err = value.Call(value, params[0])
	} else {
		value, err = value.Call(value, map[string]interface{}{})
	}
	return
}

//重写模板`include``语句实现方法，转换模板为绝对路径
func (je *JSTemplate) includeTemplate(call otto.FunctionCall) otto.Value {
	file := call.Argument(0).String() + je.ExtendType
	fileAbsPath := filepath.Join(je.templateDIR, file)
	if e := utils.FileExists(fileAbsPath); !e {
		je.err = errors.New("can't find template file:" + file)
		return otto.FalseValue()
	}
	var data []byte
	if _, ok := je.templateCaches[file]; ok {
		data = je.templateCaches[file]
	} else {
		data, je.err = ioutil.ReadFile(fileAbsPath)
		if je.Cache && je.err != nil {
			je.templateCaches[file] = data
		}
	}
	if je.err != nil {
		return otto.FalseValue()
	}
	value, err := je.processTemplate(data, call.Argument(1))
	if err != nil {
		return otto.FalseValue()
	}
	return value
}
