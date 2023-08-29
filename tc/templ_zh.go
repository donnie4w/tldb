// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package tc

var alterText = `<html>

<head>
    <title>tldb</title>
</head>

<body class="body">
    <h3>修改表结构</h3>
    <hr>
    <table>
        <tr>
            <form id="alterform" action="/alter" method="post">
                <input name="type" value="1" hidden />
                <td>表名：</td>
                <td><input type="text" name="tableName" placeholder="表名" value="" /></td>
            </form>
        </tr>
        <tr></tr>
        <tr>
            <td></td>
            <td><button onclick='javascript:document.getElementById("alterform").submit();'>检出表结构</button></td>
        </tr>
    </table>
    <hr>
    {{if ne .TableName ""}}
    <table>
        <tbody id="ctable">
            <tr>
                <td>表名：</td>
                <td><input type="text" id="tablen" placeholder="表名" value="{{ .TableName }}" /></td>
            </tr>
            {{range $k,$v := .Columns }}
            <tr>
                <td>字段名：</td>
                <td><span name="colums"><input type="text" placeholder="字段名" value="{{ $k }}" readonly />
                        <select>
                            <option value="{{$v.Type}}" selected>{{$v.Tname}}</option>
                            <option value="0">String(字符串)</option>
                            <option value="1">INT64(64位整型)</option>
                            <option value="2">INT32(32位整型)</option>
                            <option value="3">INT16(16位整型)</option>
                            <option value="4">INT8(8位整型)</option>
                            <option value="5">FLOAT64(64位浮点型)</option>
                            <option value="6">FLOAT32(32位浮点型)</option>
                            <option value="7">BINARY(字节数组)</option>
                            <option value="8">Byte(字节)</option>
                            <option value="9">Unsigned INT64</option>
                            <option value="10">Unsigned INT32</option>
                            <option value="11">Unsigned INT16</option>
                            <option value="12">Unsigned INT8</option>
                        </select>
                        建字段索引
                        {{if $v.Idx}}
                        <input type="checkbox" checked onclick="return false" />
                        {{else}}
                        <input type="checkbox" onclick="return false" />
                        {{end}}
                    </span>
                </td>
            </tr>
            {{end}}
        </tbody>

        <div id="createDiv">
        </div>
        <form id="createform" action="/create" method="post">
        </form>
        <tr></tr>
        <tr>
            <td></td>
            <td><button onclick="add();">增加字段</button></td>
        </tr>
    </table>
    <hr>
    <button style="background-color: #7fbbff;width: 100px;height: 30px;font-size: large;"
        onclick="javascipt:if (confirm('确定修改表结构？')){submit();};">确定提交</button>
    <script>
        function add() {
            var tr = document.createElement("tr");
            tr.innerHTML = '<td>字段名：</td><td><span name="colums"><input type="text" placeholder="字段名" value="" />'
                + ' <select name="fieldtype"><option value="0" selected>String(字符串)</option><option value="1">INT64(64位整型)</option><option value="2">INT32(32位整型)</option><option value="3">INT16(16位整型)</option><option value="4">INT8(8位整型)</option><option value="5">FLOAT64(64位浮点型)</option><option value="6">FLOAT32(32位浮点型)</option><option value="7">BINARY(字节数组)</option><option value="8">Byte(字节)</option><option value="9">Unsigned INT64</option><option value="10">Unsigned INT32</option><option value="11">Unsigned INT16</option><option value="12">Unsigned INT8</option></select>'
                + ' 建字段索引 <input type="checkbox" /></span></td><button onclick="del(this);">删除</button>';
            document.getElementById('ctable').appendChild(tr);
        }

        function del(obj) {
            obj.parentNode.parentNode.removeChild(obj.parentNode);
        }

        function submit() {
            var vs = document.getElementsByName('colums').values();
            var ss = [];
            for (var cn of vs) {
                var tv = cn.getElementsByTagName("input")[0].value;
                var ft = cn.getElementsByTagName("select")[0].value;
                var iv = cn.getElementsByTagName("input")[1].checked;
                ss.push('<input hidden name="colum" value="' + tv + '" /><input hidden name="ftype" value="' + ft + '" /><input hidden name="index" value="' + iv + '" />')
            }
            ss.push('<input name="type" value="2" hidden />');
            ss.push('<input name="tableName" hidden  value="' + document.getElementById("tablen").value + '" />');
            var s = ss.join('');
            document.getElementById('createDiv').innerHTML = document.getElementById('createDiv').innerHTML + s;
            document.getElementById("createform").innerHTML = document.getElementById('createDiv').innerHTML;
            document.getElementById('createform').submit();
        }
    </script>
    {{end}}
</body>

</html>`
var createText = `<html>

<head>
    <title>tldb</title>
</head>

<body class="body">
    <h3>新建表结构</h3>
    <hr>
    <table>
        <tbody id="ctable">
            <tr>
                <td>表名：</td>
                <td><input type="text" id="tablen" placeholder="表名" value="" /></td>
            </tr>
            <tr>
                <td>字段名：</td>
                <td><span name="colums"><input type="text" placeholder="字段名" value="" />
                        <select>
                            <option value="0" selected>String(字符串)</option>
                            <option value="1">INT64(64位整型)</option>
                            <option value="2">INT32(32位整型)</option>
                            <option value="3">INT16(16位整型)</option>
                            <option value="4">INT8(8位整型)</option>
                            <option value="5">FLOAT64(64位浮点型)</option>
                            <option value="6">FLOAT32(32位浮点型)</option>
                            <option value="7">BINARY(字节数组)</option>
                            <option value="8">Byte(字节)</option>
                            <option value="9">Unsigned INT64</option>
                            <option value="10">Unsigned INT32</option>
                            <option value="11">Unsigned INT16</option>
                            <option value="12">Unsigned INT8</option>
                        </select>
                        建字段索引<input type="checkbox" /></span></td>
            </tr>
        </tbody>

        <div id="createDiv">
        </div>
        <form id="createform" action="/create" method="post">
        </form>
        <tr></tr>
        <tr>
            <td></td>
            <td><button onclick="add();">增加字段</button></td>
        </tr>
    </table>

    <hr>
    <button style="background-color: #7fbbff;width: 100px;height: 30px;font-size: large;"
        onclick="javascipt:if (confirm('确定创建表？')){submit();};">确定提交</button>
    <script>
        function add() {
            var tr = document.createElement("tr");
            tr.innerHTML = '<td>字段名：</td><td><span name="colums"><input type="text" placeholder="字段名" value="" />' 
                +' <select name="fieldtype"><option value="0" selected>String(字符串)</option><option value="1">INT64(64位整型)</option><option value="2">INT32(32位整型)</option><option value="3">INT16(16位整型)</option><option value="4">INT8(8位整型)</option><option value="5">FLOAT64(64位浮点型)</option><option value="6">FLOAT32(32位浮点型)</option><option value="7">BINARY(字节数组)</option><option value="8">Byte(字节)</option><option value="9">Unsigned INT64</option><option value="10">Unsigned INT32</option><option value="11">Unsigned INT16</option><option value="12">Unsigned INT8</option></select>'
                +' 建字段索引<input type="checkbox" /></span></td><button onclick="del(this);">删除</button>';
            document.getElementById('ctable').appendChild(tr);
        }

        function del(obj) {
            obj.parentNode.parentNode.removeChild(obj.parentNode);
        }

        function submit() {
            var vs = document.getElementsByName('colums').values();
            var ss = [];
            for (var cn of vs) {
                var tv = cn.getElementsByTagName("input")[0].value;
                var ft = cn.getElementsByTagName("select")[0].value;
                var iv = cn.getElementsByTagName("input")[1].checked;
                ss.push('<input hidden name="colum" value="' + tv + '" /><input hidden name="ftype" value="' + ft + '" /><input hidden name="index" value="' + iv + '" />')
            }
            ss.push('<input name="type" value="1" hidden />');
            ss.push('<input name="tableName" hidden  value="' + document.getElementById("tablen").value + '" />');
            var s = ss.join('');
            document.getElementById('createDiv').innerHTML = document.getElementById('createDiv').innerHTML + s;
            document.getElementById("createform").innerHTML = document.getElementById('createDiv').innerHTML;
            document.getElementById('createform').submit();
        }
    </script>

</body>

</html>`
var dataText = `<html>

<head>
    <title>tldb</title>
    <style>
        .important {
            color: rgb(200, 20, 20);
            font-weight: bold;
        }

        .body {
            background-color: rgb(254, 254, 254);
        }
    </style>
</head>

<body class="body">
    <span>
        <h3 style="display: inline;">tldb 操作平台</h3>
    </span>
    &nbsp;&gt;&gt;
    <h4 class="important" style="display: inline;width:100%;">集群状态：</h4>
    {{if .Stat }}
    <span style="display:inline-block; background-color: aquamarine;width: 200px;">运行</span>
    {{else}}
    <span style="display:inline-block; background-color: rgb(255, 0, 0); width: 200px;">初始化... &#9200;</span>
    {{end}}
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=en">[EN]</a></h6>
    </span>
    <hr>
    <a href='/init'>用户管理</a>
    <a href='/sysvar'>集群环境</a>
    <a href='/sys'>节点参数</a>
    <a href='/data' style="font-weight: bold;">数据操作</a>
    <a href='/mq'>MQ数据</a>
    <a href='/log'>系统日志</a>
    <a href='/monitor'>监控</a>
    <a href='/login'>登录</a>
    <hr>
    <div>
        <div style="font-size: large; font-weight: bold;">数据表结构</div>
        <div style="overflow:scroll;max-height: 300px;">
            <table border="1" class="important">
                <tr>
                    <th>表名</th>
                    <th>索引字段</th>
                    <th>字段名</th>
                    <th>当前ID</th>
                    <th>导出表数据</th>
                </tr>
                {{range $k,$v := .Tb }}
                <tr>
                    <td>{{ $v.Name }}</td>
                    <td>{{ $v.Idxs }}</td>
                    <td>{{ $v.Columns }}</td>
                    <td>{{ $v.Seq }}</td>
                    <td><button onclick="javascipt:if (confirm('表数据量大时，可能占用服务器大量内存，谨慎导出.是否确定导出表数据？')){exportdata(this);};">导出</button></td>
                </tr>
                {{end}}
            </table>
        </div>
        <hr>
        <button onclick="openPage('/create')">新建表</button>&nbsp;<button
            onclick="openPage('/alter')">Alter表</button>&nbsp;<button
            onclick="openPage('/drop')">删除表</button>&nbsp;<button
            onclick="openPage('/insert')">插入数据</button>&nbsp;<button
            onclick="openPage('/update')">更新数据</button>&nbsp;<button onclick="openPage('/delete')">删除数据</button>
        <hr>
        <form id="exportform" action="/export" method="post">
            <input name="exportName" id="exportName" value="" hidden>
        </form>
        <span><b>根据ID查询数据</b></span>
        <form id="dataform" action="/data" method="post">
            <input name="type" value="1" hidden />
            表名<input name="tableName" placeholder="表名" value="{{ .Sb.Name }}" />
            表ID<input name="tableId" placeholder="表ID" value="{{ .Sb.Id }}" />
            <input type="submit" value="查询" />
        </form>
        <span><b>根据ID查询多条数据</b></span>
        <form id="dataform" action="/data" method="post">
            <input name="type" value="3" hidden />
            表名<input name="tableName" placeholder="表名" value="{{ .Sb.Name }}" />
            表ID<input name="start" placeholder="起始ID" value="{{ .Sb.StartId }}" />
            查询条数<input name="limit" placeholder="查询条数" value="{{ .Sb.Limit }}" />
            <input type="submit" value="查询" />
        </form>

        <span><b>根据索引查询</b></span>
        <form id="dataform" action="/data" method="post">
            <input name="type" value="2" hidden />
            表名<input name="tableName" placeholder="表名" value="{{ .Sb.Name }}" />
            字段名<input name="cloName" placeholder="字段名" value="{{ .Sb.ColumnName }}" />
            字段值<input name="cloValue" placeholder="字段值" value="{{ .Sb.ColumnValue }}" />
            起始<input name="start" placeholder="起始" value="{{ .Sb.StartId }}" />
            查询条数<input name="limit" placeholder="查询条数" value="{{ .Sb.Limit }}" />
            <input type="submit" value="查询" />
        </form>
        <hr>
        <h4>查询结果：</h4>
        <h3 class="important">{{ .Sb.Name }}</h3>
        <table border="1">
            <tr class="important" style="font-size: small;">
                <th>表Id</th>
                {{range $k,$v := .ColName }}
                <th>字段名</th>
                <th>字段值</th>
                {{end}}
            </tr>
            {{range $k,$v := .Tds }}
            <tr>
                <td>{{ $v.Id }}</td>
                {{range $k1,$v1 := $v.Columns }}
                <td>{{ $k1 }}</td>
                <td><textarea readonly>{{ $v1 }}</textarea></td>
                {{end}}
            </tr>
            {{end}}
        </table>
    </div>
    <script>
        function openPage(o) {
            window.open(o, "TLDB", "height=500, width=800, top=50, left=100,menubar=0,status=0,titlebar=0");
        }

        function exportdata(o) {
            document.getElementById("exportName").value = o.parentNode.parentNode.cells[0].innerText;
            document.getElementById("exportform").submit();
        }
    </script>
</body>

</html>`
var deleteText = `<html>

<head>
    <title>tldb</title>
</head>

<body class="body">
    <h3>删除表数据</h3>
    <hr>
    <span><b>根据ID查询数据</b></span>
    <form id="deleteform" action="/delete" method="post">
        <input name="type" value="1" hidden />
        表名<input name="tableName" placeholder="表名" value="{{ .TableName }}" />
        表ID<input name="tableId" placeholder="表ID"  value="" />
        <input type="submit" value="检出数据" />
    </form>
    <hr>
    {{if ne .TableName ""}}
    <table border="0">
        <tbody id="ctable">
            <form id="deleteform2" action="/delete" method="post">
                <input name="type" value="2" hidden />
                <tr>
                    <th>表名：</th>
                    <td><input type="text" name="tableName" placeholder="表名" value="{{ .TableName }}" style="border: none;" /></td>
                </tr>
                <tr>
                    <th>ID：</th>
                    <td><input type="text" name="tableId" placeholder="ID" value="{{ .ID }}"  readonly style="border: none;"/></td>
                </tr>
                {{range $k,$v := .ColumnValue }}
                <tr>
                    <th>字段名：</th>
                    <td><input type="text" placeholder="字段" value="{{ $k }}"
                                readonly style="border: none;" />
                        </td>
                    <td><textarea readonly>{{ $v }}</textarea></td>
                </tr>
                {{end}}
            </form>
        </tbody>
        <tr>
            <td></td>
            <td><button onclick="javascipt:if (confirm('确定删除数据?')){document.getElementById('deleteform2').submit();};">删除</button></td>
        </tr>
    </table>
    {{end}}

</body>

</html>`
var initText = `<html>

<head>
    <title>tldb</title>
    <style>
        .important {
            color: rgb(200, 20, 20);
            font-weight: bold;
        }

        .body {
            background-color: rgb(254, 254, 254);
        }
    </style>
</head>

<body class="body">
    {{if not .Init}}
    <span>
        <h3 style="display: inline;">tldb 操作平台</h3>
    </span>
    &nbsp;&gt;&gt;
    <h4 class="important" style="display: inline;width:100%;">集群状态：</h4>
    {{if .Stat }}
    <span style="display:inline-block; background-color: aquamarine;width: 200px;">运行</span>
    {{else}}
    <span style="display:inline-block; background-color: rgb(255, 0, 0); width: 200px;">初始化... &#9200;</span>
    {{end}}
    {{else if .Init}}
    <h3 style="display: inline;">tldb 操作平台 </h3>
    {{end}}
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=en">[EN]</a></h6>
    </span>
    <hr>
    <a href='/init' style="font-weight: bold;">用户管理</a>
    <a href='/sysvar'>集群环境</a>
    <a href='/sys'>节点参数</a>
    <a href='/data'>数据操作</a>
    <a href='/mq'>MQ数据</a>
    <a href='/log'>系统日志</a>
    <a href='/monitor'>监控</a>
    <a href='/login'>登录</a>
    {{if .ShowCreate }}
    <hr>
    <div>
        <h3>用户账户管理</h3>
        <hr>
        <h4>新建管理员 <h5 class="important">{{ .Show }}</h5>
        </h4>
        <form id="createAdminform" action="/init?type=1" method="post">
            <input name="adminName" placeholder="用户名" />
            <input name="adminPwd" placeholder="密码" type="password" />
            管理员<input name="adminType" type="radio" value="1" checked />&nbsp;&nbsp;
            {{if not .Init}}
            数据管理员<input name="adminType" type="radio" value="2" />
            {{end}}
            <input type="submit" value="新建管理员" />
        </form>

        {{if not .Init}}
        <hr>
        <h4>新建MQ客户端</h4>
        <form id="createMQform" action="/init?type=1" method="post">
            <input name="mqName" placeholder="MQ用户名" />
            <input name="mqPwd" placeholder="密码" type="password" />
            <input type="submit" value="新建MQ客户端" />
        </form>
        <hr>
        <h4>新建数据库客户端</h4>
        <form id="createCliform" action="/init?type=1" method="post">
            <input name="cliName" placeholder="客户端用户名" />
            <input name="cliPwd" placeholder="密码" type="password" />
            <input type="submit" value="新建数据库客户端" />
        </form>
        {{end}}
    </div>
    <hr>
    {{end}}
    {{if not .Init}}
    <hr>
    <div class="important" style="font-size: small;">
        <h4>后台管理员</h4>
        {{range $k,$v := .AdminUser}}
        <form id="adminform" action="/init?type=2" method="post">
            <input name="adminName" value='{{ $k }}' readonly /> 权限:{{ $v }}
            <input type="button" value="删除用户" onclick="javascipt:if (confirm('确定删除?')){this.parentNode.submit();};" />
        </form>
        {{end}}
        <hr>
        <h4>MQ客户端</h4>
        {{range $k,$v := .MqUser }}
        <form id="mqform" action="/init?type=2" method="post">
            <input name="mqName" value="{{ $v }}" readonly />
            <input type="button" value="删除用户" onclick="javascipt:if (confirm('确定删除?')){this.parentNode.submit();};" />
        </form>
        {{end}}
        <hr>
        <h4>数据库客户端</h4>
        {{range $k,$v := .CliUser }}
        <form id="cliform" action="/init?type=2" method="post">
            <input name="cliName" value="{{ $v }}" readonly />
            <input type="button" value="删除用户" onclick="javascipt:if (confirm('确定删除?')){this.parentNode.submit();};" />
        </form>
        {{end}}
    </div>
    <hr>
    {{end}}

</html>`
var insertText = `<html>

<head>
    <title>tldb</title>
</head>

<body class="body">
    <h3>新增表数据</h3>
    <hr>
    <table>
        <tr>
            <form id="insertform" action="/insert" method="post">
                <input name="type" value="1" hidden />
                <td>表名：</td>
                <td><input type="text" name="tableName" placeholder="表名" value="" /></td>
            </form>
        </tr>
        <tr></tr>
        <tr>
            <td></td>
            <td><button onclick='javascript:document.getElementById("insertform").submit();'>检出表结构</button></td>
        </tr>
    </table>
    <hr>
    {{if ne .TableName ""}}
    <table>
        <tbody id="ctable">
            <form id="insertform2" action="/insert" method="post">
                <input name="type" value="2" hidden />
                <tr>
                    <th>表名：</th>
                    <td><input type="text" name="tableName" placeholder="表名" value="{{ .TableName }}" /></td>
                </tr>
                {{range $k,$v := .Columns }}
                <tr>
                    <th>字段名：</th>
                    <td><span><input type="text" name="colums" placeholder="字段名" value="{{ $k }}" readonly />
                            <input type="text" name="values" placeholder="字段值" value="" style="width: 400px;" />
                        </span>
                    </td>
                </tr>
                {{end}}
            </form>
        </tbody>
        <tr>
            <td></td>
            <td><button onclick="javascipt:if (confirm('confirm insert?')){document.getElementById('insertform2').submit();};">提交</button></td>
        </tr>
    </table>
    {{end}}
</body>

</html>`
var loadText = `<html>

<body style="text-align:center;">
    <h2>导入数据中</h2>
    <div>
        <h3>已导入数据：<span id='s'></span>条</h3>
    </div>
    <div id="e">
    </div>
    <div id="e2">
    </div>
</body>
<script type="text/javascript">
    var pro = window.location.protocol;
    var wspro = "ws:";
    if (pro === "https:") {
        wspro = "wss:";
    }
    var ws = new WebSocket(wspro + "//" + window.location.host + "/loadData");
    ws.onmessage = function (evt) {
        if (evt.data == "") {
            document.getElementById('e').innerHTML = "数据导入完成";
        } else {
            document.getElementById('s').innerHTML = evt.data;
        }
    }
    ws.onclose = function (evt) {
        document.getElementById("e2").innerHTML = '<hr><h4>请<a href="javascript:window.history.go(-1)">点击此处</a>返回。<h4>'
    };
    ws.onopen = function (evt) {
    };
    ws.onerror = function (evt, e) {
    };
</script>

</html>`
var loginText = `<html>
<head>
    <title>tldb</title>
    <style>
        .important{
            color: rgb(200, 20, 20);
            font-weight: bold;
        }
        .body{
            background-color: rgb(254, 254, 254);
        }
    </style>
</head>
<body class="body">
    <h3 style="display: inline;">tldb 操作平台</h3>
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=en">[EN]</a></h6>
    </span>
    <hr>
    <div id="login">
        <h3>登录</h3>
        <form id="loginform" action="/login" method="post">
            <input name="type" value="1" hidden />
            <input name="name" placeholder="用户名" />
            <input name="pwd" placeholder="密码" type="password" />
            <input type="submit" value="登录" />
        </form>
    </div>
    <hr>
</html>`
var mqText = `<html>

<head>
    <title>tldb</title>
    <style>
        .important{
            color: rgb(200, 20, 20);
            font-weight: bold;
        }
        .body{
            background-color: rgb(254, 254, 254);
        }
    </style>
</head>
<body class="body">
    <span>
        <h3 style="display: inline;">tldb 操作平台</h3>
    </span>
    &nbsp;&gt;&gt;
    <h4 class="important" style="display: inline;width:100%;">集群状态：</h4>
    {{if .Stat }}
    <span style="display:inline-block; background-color: aquamarine;width: 200px;">运行</span>
    {{else}}
    <span style="display:inline-block; background-color: rgb(255, 0, 0); width: 200px;">初始化... &#9200;</span>
    {{end}}
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=en">[EN]</a></h6>
    </span>
    <hr>
    <a href='/init'>用户管理</a>
    <a href='/sysvar'>集群环境</a>
    <a href='/sys'>节点参数</a>
    <a href='/data'>数据操作</a>
    <a href='/mq' style="font-weight: bold;">MQ数据</a>
    <a href='/log'>系统日志</a>
    <a href='/monitor'>监控</a>
    <a href='/login'>登录</a>
    <hr>
    <div>
        <div style="font-size: large; font-weight: bold;">数据操作</div>
        <hr>
        <div style="overflow:scroll;max-height: 300px;">
            <table border="1" class="important" style="font-size: x-small;">
                <tr>
                    <th>发布字段(不含MEM)</th>
                    <th>当前ID</th>
                    <th>当前订阅数</th>
                    <th>删除</th>
                </tr>
                {{range $k,$v := .Tb }}
                <tr>
                    <td>{{ $v.Name }}</td>
                    <td>{{ $v.Seq }}</td>
                    <td>{{ $v.Sub }}</td>
                    <td>
                        <form action="/mq" method="post">
                            <input name = "atype" value="1" hidden />
                            <input name = "tableName" value="{{ $v.Name }}" hidden />
                            <input type="button" value="删除" onclick="javascipt:if (confirm('确定删除?')){this.parentNode.submit();};" />
                        </form>
                    </td>
                </tr>
                {{end}}
            </table>
        </div>
        <hr>
        <span><b>根据ID查询MQ数据</b></span>
        <form id="dataform" action="/mq" method="post">
            <input name = "type" value="2" hidden />
            Topic<input name="tableName" placeholder="topic" value="{{ .Sb.Name }}" />
            ID<input name="tableId" placeholder="id" value="{{ .Sb.Id }}" />
            <input type="submit" value="查询" />
        </form>
        <hr>
        <span><b>根据ID查询MQ数据</b></span>
        <form id="dataform" action="/mq" method="post">
            <input name = "type" value="3" hidden />
            Topic<input name="tableName" placeholder="topic" value="{{ .Sb.Name }}" />
            起始ID<input name="start" placeholder="起始id" value="{{ .Sb.StartId }}" />
            查询条数<input name="limit" placeholder="查询条数" value="{{ .Sb.Limit }}" />
            <input type="submit" value="查询" />
        </form>
        <hr>
        <h4>查询结果：</h4>
        <h3 class="important">{{ .Sb.Name }}</h3>
        <table border="1" style="width: 100%;">
            <tr class="important" style="font-size: small;">
                <th style="width: 45px;">消息ID</th>
                <th>数据</th>
            </tr>
            {{range $k,$v := .Tds }}
            <tr>
                <td>{{ $v.Id }}</td>
                {{range $k1,$v1 := $v.Columns }}
                <td><textarea readonly style="width: 100%;">{{ $v1 }}</textarea></td>
                {{end}}
            </tr>
            {{end}}
        </table>
    </div>
</body>

</html>`
var sysText = `<html>

<head>
    <title>tldb</title>
    <style>
        .important {
            color: rgb(200, 20, 20);
            font-weight: bold;
        }

        .body {
            background-color: rgb(254, 254, 254);
        }
    </style>
</head>

<body class="body">
    <span>
        <h3 style="display: inline;">tldb 操作平台</h3>
    </span>
    &nbsp;&gt;&gt;
    <h4 class="important" style="display: inline;width:100%;">集群状态：</h4>
    {{if .Stat }}
    <span style="display:inline-block; background-color: aquamarine;width: 200px;">运行</span>
    {{else}}
    <span style="display:inline-block; background-color: rgb(255, 0, 0); width: 200px;">初始化... &#9200;</span>
    {{end}}
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=en">[EN]</a></h6>
    </span>
    <hr>
    <a href='/init'>用户管理</a>
    <a href='/sysvar'>集群环境</a>
    <a href='/sys' style="font-weight: bold;">节点参数</a>
    <a href='/data'>数据操作</a>
    <a href='/mq'>MQ数据</a>
    <a href='/log'>系统日志</a>
    <a href='/monitor'>监控</a>
    <a href='/login'>登录</a>
    <hr>
    <div>
        <table border="1" style="font-size: 15px;">
            <tr>
                <th>名称</th>
                <th>当前值</th>
                <th>启动设置参数</th>
                <th>说明</th>
            </tr>   
            <tr>
                <td>本地数据文件：</td>
                <td class="important">{{ .SYS.DBFILEDIR }}</td>
                <td> -dir</td>
                <td>数据文件地址</td>
            </tr>
            <tr>
                <td>BINLOG日志文件大小：</td>
                <td>{{ .SYS.BINLOGSIZE }}(MB)</td>
                <td> -binsize</td>
                <td>binlog数据文件按每{{ .SYS.BINLOGSIZE }}M进行压缩备份</td>
            </tr>
            <tr>
                <td>MQ是否使用tls</td>
                <td>{{ .SYS.MQTLS }}</td>
                <td> -clitls</td>
                <td>wss协议 访问MQ服务</td>
            </tr>
            <tr>
                <td>web admin是否使用tls</td>
                <td>{{ .SYS.ADMINTLS }} </td>
                <td> -admintls</td>
                <td>https协议 访问管理后台 </td>
            </tr>
            <tr>
                <td>客户端传输是否使用tls</td>
                <td>{{ .SYS.CLITLS }} </td>
                <td> -mqtls</td>
                <td>sslsocket数据库客户端访问服务器</td>
            </tr>
            <tr>
                <td>客户端 crt文件地址</td>
                <td>{{ .SYS.CLICRT }}</td>
                <td> -clicrt</td>
                <td>客户端安全访问协议的SSL crt证书文件地址</td>
            </tr>
            <tr>
                <td>客户端 key文件地址</td>
                <td>{{ .SYS.CLIKEY }} </td>
                <td> -clikey</td>
                <td>客户端安全访问协议的SSL key证书文件地址</td>
            </tr>
            <tr>
                <td>mq crt文件地址</td>
                <td>{{ .SYS.MQCRT }} </td>
                <td> -mqcrt</td>
                <td>MQ安全访问协议的SSL crt证书文件地址</td>
            </tr>
            <tr>
                <td>mq key文件地址</td>
                <td>{{ .SYS.MQKEY }} </td>
                <td> -mqkey</td>
                <td>MQ安全访问协议的SSL key证书文件地址</td>
            </tr>
            <tr>
                <td>web admin crt文件地址</td>
                <td>{{ .SYS.ADMINCRT }} </td>
                <td> -admincrt</td>
                <td>管理后台安全访问协议的SSL crt证书文件地址</td>
            </tr>
            <tr>
                <td>web admin key文件地址</td>
                <td>{{ .SYS.ADMINKEY }}</td>
                <td> -adminkey</td>
                <td>管理后台安全访问协议的SSL key证书文件地址</td>
            </tr>
            <tr>
                <td>增删改上限并发数</td>
                <td>{{ .SYS.COCURRENT_PUT }} </td>
                <td> -put</td>
                <td>客户端链接增删改并发数,超过则排队等待</td>
            </tr>
            <tr>
                <td>查询上限并发数</td>
                <td>{{ .SYS.COCURRENT_GET }} </td>
                <td> -get</td>
                <td>客户端链接查询并发数,超过则排队等待</td>
            </tr>
            <tr>
                <td>集群命名空间</td>
                <td>{{ .SYS.NAMESPACE }}</td>
                <td> -ns</td>
                <td>集群中节点命名空间必须相同，否则不能连接</td>
            </tr>
            <tr>
                <td>节点集群链接密码</td>
                <td>{{ .SYS.PWD }}</td>
                <td> -pwd</td>
                <td>集群节点之间链接密码</td>
            </tr>
            <tr>
                <td>节点集群链接SSL加密验证 公钥地址</td>
                <td>{{ .SYS.PUBLICKEY }}</td>
                <td> -publickey</td>
                <td>默认使用tldb程序中公钥;可另指定公钥地址</td>
            </tr>
            <tr>
                <td>节点集群链接SSL加密验证 私钥地址</td>
                <td>{{ .SYS.PRIVATEKEY }}</td>
                <td> -privatekey</td>
                <td>默认使用tldb程序中私钥;可另指定私钥地址</td>
            </tr>
            <tr>
                <td>节点集群链接地址</td>
                <td>{{ .SYS.ADDR }}</td>
                <td> -cs</td>
                <td>节点之间集群服务链接地址</td>
            </tr>
            <tr>
                <td>MQ服务地址</td>
                <td>{{ .SYS.MQADDR }}</td>
                <td> -mq</td>
                <td>MQ服务地址</td>
            </tr>
            <tr>
                <td>数据库客户端服务地址</td>
                <td>{{ .SYS.CLIADDR }}</td>
                <td> -cli</td>
                <td>数据库客户端服务器地址.</td>
            </tr>
            <tr>
                <td>管理后台服务地址</td>
                <td>{{ .SYS.WEBADMINADDR }}</td>
                <td> -admin</td>
                <td>web管理后台服务地址.</td>
            </tr>
            <tr>
                <td>集群节点数下限</td>
                <td>{{ .SYS.CLUSTER_NUM }}</td>
                <td> -clus</td>
                <td>默认系统自动分配;值为0时,节点单点运行,否则集群运行.</td>
            </tr>
            <tr>
                <td>集群节点数下限固定</td>
                <td>{{ .SYS.CLUSTER_NUM_FINAL }}(默认false:系统分配)</td>
                <td> -clus_final</td>
                <td>默认系统自动分配大小.值true时,-clus非零参数值生效</td>
            </tr>
            <tr>
                <td>程序版本</td>
                <td>v{{ .SYS.VERSION }}</td>
                <td></td>
                <td>当前程序的开发版本</td>
            </tr>
        </table>
        <hr>
        <hr>

        <span style="font-size:large;font-weight: bold;">节点导入Bin.Log 压缩包数据[以数据追加的方式]</span>
        <span style="font-size: xx-small;">导入文件为tldb生成的压缩binlog文件</span>
        <form id="loadForm1" action="/sys" method="post" enctype="multipart/form-data">
            <input name="atype" value="1" hidden />
            <input type="file" id="loadfile1" name="loadfile1" />
            <button
                onclick="javascipt:if (confirm('导入数据,可能导致本节点数据与其他节点数据不一致。确定导入数据?')){this.parentNode.submit();};">导入数据</button>
        </form>

        <hr>
        <span style="font-size:large;font-weight: bold;">节点导入Bin.Log 压缩包数据[以数据覆盖的方式]</span>
        <span style="font-size: xx-small;">导入文件为tldb生成的压缩binlog文件</span>
        <form id="loadForm2" action="/sys" method="post" enctype="multipart/form-data">
            <input name="atype" value="2" hidden />
            <input type="file" id="loadfile2" name="loadfile2" />
            <button
                onclick="javascipt:if (confirm('导入数据,可能导致本节点数据与其他节点数据不一致。确定导入数据?')){this.parentNode.submit();};">导入数据</button>
        </form>
        <hr>
        <hr>
        <form id="sysForm" action="/sys" method="post" enctype="multipart/form-data">
            <input name="atype" value="3" hidden />
            <button
                onclick="javascipt:if (confirm('确定关闭本节点所有服务？')){document.getElementById('sysForm').submit();};">关闭本节点服务</button>
        </form>
    </div>
</body>

</html>`
var sysvarText = `<html>

<head>
    <title>tldb</title>
    <meta http-equiv="refresh" content="30">
    <style>
        .important {
            color: rgb(200, 20, 20);
            font-weight: bold;
        }

        .body {
            background-color: rgb(254, 254, 254);
        }
    </style>
</head>

<body class="body">
    <span>
        <h3 style="display: inline;">tldb 操作平台</h3>
    </span>
    &nbsp;&gt;&gt;
    <h4 class="important" style="display: inline;width:100%;">集群状态：</h4>
    {{if .Stat }}
    <span style="display:inline-block; background-color: aquamarine;width: 200px;">运行</span>
    {{else}}
    <span style="display:inline-block; background-color: rgb(255, 0, 0); width: 200px;">初始化... &#9200;</span>
    {{end}}
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=en">[EN]</a></h6>
    </span>
    <hr>
    <a href='/init'>用户管理</a>
    <a href='/sysvar' style="font-weight: bold;">集群环境</a>
    <a href='/sys'>节点参数</a>
    <a href='/data'>数据操作</a>
    <a href='/mq'>MQ数据</a>
    <a href='/log'>系统日志</a>
    <a href='/monitor'>监控</a>
    <a href='/login'>登录</a>
    <hr>
    <div>
        <table border="1" style="font-size:15px;">
            <tr>
                <td>节点启动时间(本地时间)</td>
                <td colspan="2">{{ .SYS.StartTime }}</td>
            </tr>
            <tr>
                <td>节点服务器时间</td>
                <td colspan="2">{{ .SYS.LocalTime }}</td>
            </tr>
            <tr>
                <td>集群系统修正时间</td>
                <td colspan="2">{{ .SYS.Time }}</td>
            </tr>
            <tr>
                <td>节点UUID</td>
                <td class="important" colspan="2">{{ .SYS.UUID }}</td>
            </tr>
            <tr>
                <td>当前状态为RUN集群节点</td>
                <td class="important" colspan="2">{{ .SYS.RUNUUIDS }}</td>
            </tr>
            <tr>
                <td class="important">节点运行状态</td>
                <form id="statForm" action="/sysvar" method="post">
                    <input name="atype" value="3" hidden />
                </form>
                <form id="statForm2" action="/sysvar" method="post">
                    <input name="atype" value="4" hidden />
                </form>
                <td colspan="2">
                    {{if eq .SYS.STAT "0"}}
                    <span style="display:inline-block; background-color: rgb(255, 0, 0);width: 70px;">就绪 &#9200;</span>
                    {{else if eq .SYS.STAT "1"}}
                    <span style="display:inline-block; background-color: rgb(255, 255, 0);width: 70px;">代理
                        &#128274;</span>
                    {{else if eq .SYS.STAT "2"}}
                    <span style="display:inline-block; background-color:aquamarine;width: 70px;">运行 &#9989;</span>
                    {{end}}
                    <span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>
                    <button
                        onclick="javascipt:if (confirm('状态修改可能导致部分客户端当前操作失败。确定重置为代理？(谨慎)')){document.getElementById('statForm').submit();};">重置为代理</button>
                    <button
                        onclick="javascipt:if (confirm('状态修改可能导致部分客户端当前操作失败。确定重置为就绪？(谨慎)')){document.getElementById('statForm2').submit();};" />重置为就绪</button>
                    {{if eq .SYS.STAT "0"}}
                    <span style="font-size: small;">目前同步{{.SYS.SyncCount}}条数据</span>
                    {{end}}
                </td>

            </tr>
            <tr>
                <td>集群时间与本地时间偏差</td>
                <form id="timeForm" action="/sysvar" method="post">
                    <input name="atype" value="2" hidden />
                    <td colspan="2">
                        <input name="time_deviation" value="{{ .SYS.TIME_DEVIATION }}"
                            style="border: none;width: 100px;" />(纳秒)
                        <input type="button" value="修改时间偏差(谨慎)"
                            onclick="javascipt:if (confirm('修改节点时间可能导致日志数据出错。确定修改节点偏差？')){document.getElementById('timeForm').submit();};" />
                    </td>
                </form>
            </tr>
            <tr>
                <td>当前节点集群服务访问地址</td>
                <td colspan="2">{{ .SYS.ADDR }}</td>
            </tr>
            <tr>
                <td>后台地址</td>
                <td colspan="2">{{ .SYS.ADMINADDR }}</td>
            </tr>
            <tr>
                <td>MQ连接地址</td>
                <td colspan="2">{{ .SYS.MQADDR }}/mq</td>
            </tr>
            <tr>
                <td>数据库客户端连接地址</td>
                <td colspan="2">{{ .SYS.CLIADDR }}</td>
            </tr>
            <tr>
                <td>集群最小节点数</td>
                <td colspan="2">
                    <div class="important">
                        {{if eq .SYS.CLUSTER_NUM "0"}}
                        单点运行
                        {{else}}
                        {{ .SYS.CLUSTER_NUM }}
                        {{end}}
                    </div>
                </td>
            </tr>
            <tr>
                <td>存放数据节点个数</td>
                <form id="storeNumForm" action="/sysvar" method="post">
                    <input name="atype" value="5" hidden />
                    <td colspan="2">
                        <input name="storeNum" value="{{ .SYS.STORENODENUM }}"
                            style="border: none;width: 50px;" />[默认0,表示全部节点]
                        <input type="button" value="修改数据节点数"
                            onclick="javascipt:if (confirm('修改数据节点数可能部分客户端当前操作失败(修改成功后，数据将同步到其他节点)。确定修改？')){document.getElementById('storeNumForm').submit();};" />
                    </td>
                </form>
            </tr>
            <tr>
                <td>当前节点并发增删改数据统计</td>
                <td colspan="2">{{ .SYS.CCPUT }}</td>
            </tr>
            <tr>
                <td>当前节点并发查询数据统计</td>
                <td colspan="2">{{ .SYS.CCGET }}</td>
            </tr>
            <tr>
                <td>当前节点启动至今增删改数据统计</td>
                <td colspan="2">{{ .SYS.COUNTPUT }}</td>
            </tr>
            <tr>
                <td>当前节点启动至今查询数据统计</td>
                <td colspan="2">{{ .SYS.COUNTGET }}</td>
            </tr>
            {{range $k,$v := .RN }}
            <tr>
                <th rowspan="8">远程节点: {{ $v.UUID }}</th>
            </tr>
            <tr>
                <td style="width: 100px;">UUID</td>
                <td>{{ $v.UUID }}</td>
            </tr>
            <tr>
                <td class="important">状态</td>
                <td>{{ $v.StatDesc }}</td>
            </tr>
            <tr>
                <td>远程地址</td>
                <td>{{ $v.Addr }}</td>
            </tr>
            <tr>
                <td>远程IP</td>
                <td>{{ $v.Host }}</td>
            </tr>
            <tr>
                <td>后台地址</td>
                <td>{{ $v.AdminAddr }}</td>
            </tr>
            <tr>
                <td>MQ地址</td>
                <td>{{ $v.MQAddr }}/mq</td>
            </tr>
            <tr>
                <td>客户端地址</td>
                <td>{{ $v.CliAddr }}</td>
            </tr>
            {{end}}
        </table>
        <hr>
        <h5 class="important">{{ .Show }}</h5>
        <h3>集群操作</h3>
        <form id="" action="/sysvar" method="post">
            <input name="atype" value="1" hidden />
            <table border="1">
                <tr>
                    <th>增加集群节点并连接</th>
                </tr>
                <tr>
                    <th>目标节点地址</th>
                    <td><input type="text" id="addr" name="addr" value="" placeholder="目标节点地址" /></td>
                    <td><input type="submit" value="确定" /></td>
                </tr>
            </table>
        </form>
        <hr>
    </div>
</body>

</html>`
var dropText = `<html>

<head>
    <title>tldb</title>
</head>

<body class="body">
    <h3>删除表及表数据</h3>
    <hr>
    <div style="overflow:scroll;max-height: 300px;">
        <table border="1" class="important">
            <tr>
                <th>表名</th>
                <th>索引字段</th>
                <th>字段名</th>
                <th>当前ID</th>
                <th>删除表及表数据</th>
            </tr>
            {{range $k,$v := .Tb }}
            <tr>
                <td>{{ $v.Name }}</td>
                <td>{{ $v.Idxs }}</td>
                <td>{{ $v.Columns }}</td>
                <td>{{ $v.Seq }}</td>
                <td><button onclick="javascipt:if (confirm('确定删除表及表数据?')){del(this);};">删除</button></td>
            </tr>
            {{end}}
        </table>
        <form id="dropform" action="/drop" method="post">
            <input name="type" value="1" hidden />
            <input id="tableName" name="tableName" value="" hidden />
        </form>
    </div>
    <script>
        function del(o) {
            var n = o.parentNode.parentNode.getElementsByTagName("td")[0].innerText;
            document.getElementById("tableName").value = n;
            document.getElementById("dropform").submit();
        }
    </script>
</body>

</html>`
var updateText = `<html>

<head>
    <title>tldb</title>
    <style>
        .important {
            color: rgb(200, 20, 20);
            font-weight: bold;
        }

        .body {
            background-color: rgb(254, 254, 254);
        }
    </style>
</head>

<body class="body">
    <h3>修改表数据</h3>
    <hr>
    <span><b>根据ID查询数据</b></span>
    <form id="updateform" action="/update" method="post">
        <input name="type" value="1" hidden />
        表名<input name="tableName" placeholder="表名" value="{{ .TableName }}" />
        表ID<input name="tableId" placeholder="表ID" value="" />
        <input type="submit" value="检出数据" />
    </form>
    <hr>
    {{if ne .TableName ""}}
    <table border="0">
        <tbody id="ctable">
            <form id="updateform2" action="/update" method="post">
                <input name="type" value="2" hidden />
                <tr>
                    <th>表名：</th>
                    <td><input type="text" name="tableName" placeholder="表名" value="{{ .TableName }}"
                            style="border: none;" /></td>
                </tr>
                <tr>
                    <th>ID：</th>
                    <td><input type="text" name="tableId" placeholder="ID" value="{{ .ID }}" readonly
                            style="border: none;" /></td>
                </tr>
                {{range $k,$v := .ColumnValue }}
                <tr>
                    <th>字段名：</th>
                    <td><input type="text" name="colums" placeholder="字段" value="{{ $k }}" readonly
                            style="border: none;" />
                    </td>
                    <td><textarea name="values">{{ $v }}</textarea></td>
                </tr>
                {{end}}
            </form>
        </tbody>
        <tr>
            <td></td>
            <td><button
                    onclick="javascipt:if (confirm('确定更新数据?')){document.getElementById('updateform2').submit();};">更新</button>
            </td>
        </tr>
    </table>
    {{end}}

</body>

</html>`

var monitorText = `<html>

<head>
    <title>tldb</title>
    <style>
        .important {
            color: rgb(200, 20, 20);
            font-weight: bold;
        }

        .body {
            background-color: rgb(254, 254, 254);
        }
    </style>
</head>

<body class="body">
    <span>
        <h3 style="display: inline;">tldb 操作平台</h3>
    </span>
    &nbsp;&gt;&gt;
    <h4 class="important" style="display: inline;width:100%;">集群状态：</h4>
    {{if . }}
    <span style="display:inline-block; background-color: aquamarine;width: 200px;">运行</span>
    {{else}}
    <span style="display:inline-block; background-color: rgb(255, 0, 0); width: 200px;">初始化... &#9200;</span>
    {{end}}
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=en">[EN]</a></h6>
    </span>
    <hr>
    <a href='/init'>用户管理</a>
    <a href='/sysvar'>集群环境</a>
    <a href='/sys'>节点参数</a>
    <a href='/data'>数据操作</a>
    <a href='/mq'>MQ数据</a>
    <a href='/log'>系统日志</a>
    <a href='/monitor' style="font-weight: bold;">监控</a>
    <a href='/login'>登录</a>
    <hr>
    <div>
        <h3>系统数据监控</h3>
        <hr>
        监控时间间隔(单位:秒)<input id="stime" placeholder="输入时间" value="3" />
        <button onclick="monitorLoad();">开始</button>
        <button onclick="stop();">停止</button>
        <button onclick="clearData();">清除数据</button>
    </div>
    <hr>
    <div>
        <table border="1" style="font-size: smaller;">
            <tr>
                <th></th>
                <th>内存分配(MB)</th>
                <th>内存分配总数(MB)</th>
                <th>内存回收次数</th>
                <th>协程数</th>
                <th>CPU核数</th>
                <th>并发增改数</th>
                <th>增改数</th>
                <th>并发查询数</th>
                <th>查询数</th>
                <th>内存使用率</th>
                <th>磁盘剩余空间(GB)</th>
                <th>CPU使用率</th>
            </tr>
            <tbody id="monitorBody">
            </tbody>
        </table>
    </div>
</body>
<script type="text/javascript">
    var pro = window.location.protocol;
    var wspro = "ws:";
    if (pro === "https:") {
        wspro = "wss:";
    }
    var wsmnt = null;
    var id = 1;
    function WS() {
        this.ws = null;
    }

    WS.prototype.monitor = function () {
        let obj = this;
        this.ws = new WebSocket(wspro + "//" + window.location.host + "/monitorData");
        this.ws.onopen = function (evt) {
            obj.ws.send(document.getElementById("stime").value);
        }
        this.ws.onmessage = function (evt) {
            if (evt.data != "") {
                var json = JSON.parse(evt.data);
                var tr = document.createElement('tr');
                var d = '<td style="font-weight: bold;">' + id++ + '</td>'
                    + '<td>' + Math.round(json.Alloc / (1 << 20)) + '</td>'
                    + '<td>' + Math.round(json.TotalAlloc / (1 << 20)) + '</td>'
                    + '<td>' + json.NumGC + '</td>'
                    + '<td>' + json.NumGoroutine + '</td>'
                    + '<td>' + json.NumCPU + '</td>'
                    + '<td>' + json.CcPut + '</td>'
                    + '<td>' + json.CountPut + '</td>'
                    + '<td>' + json.CcGet + '</td>'
                    + '<td>' + json.CountGet + '</td>'
                    + '<td>' + Math.round(json.RamUsage * 10000) / 100 + '%</td>'
                    + '<td>' + json.DiskFree + '</td>'
                    + '<td>' + Math.round(json.CpuUsage * 100) / 100 + '%</td>';
                tr.innerHTML = d;
                document.getElementById("monitorBody").appendChild(tr);
            }
        }
    }

    WS.prototype.close = function () {
        this.ws.close();
    }

    function monitorLoad() {
        if (typeof wsmnt != "undefined" && wsmnt != null && wsmnt != "") {
            wsmnt.close();
        }
        wsmnt = new (WS);
        wsmnt.monitor();
    }

    function stop() {
        if (typeof wsmnt != "undefined" && wsmnt != null && wsmnt != "") {
            wsmnt.close();
        }
    }

    function clearData() {
        document.getElementById("monitorBody").innerHTML = "";
    }

</script>

</html>`
