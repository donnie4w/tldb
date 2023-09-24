// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package tc

var alterEnText = `<html>

<head>
    <title>tldb</title>
</head>

<body class="body">
    <h3>Alter table structure</h3>
    <hr>
    <table>
        <tr>
            <form id="alterform" action="/alter" method="post">
                <input name="type" value="1" hidden />
                <td>table name：</td>
                <td><input type="text" name="tableName" placeholder="table name" value="" /></td>
            </form>
        </tr>
        <tr></tr>
        <tr>
            <td></td>
            <td><button onclick='javascript:document.getElementById("alterform").submit();'>Check out the table structure</button></td>
        </tr>
    </table>
    <hr>
    {{if ne .TableName ""}}
    <table>
        <tbody id="ctable">
            <tr>
                <td>table name：</td>
                <td><input type="text" id="tablen" placeholder="table name" value="{{ .TableName }}" /></td>
            </tr>
            {{range $k,$v := .Columns }}
            <tr>
                <td>field name：</td>
                <td><span name="colums"><input type="text" placeholder="field name" value="{{ $k }}" readonly />
                    <select>
                        <option value="{{$v.Type}}" selected>{{$v.Tname}}</option>
                        <option value="0">String</option>
                        <option value="1">INT64(64-bit)</option>
                        <option value="2">INT32(32-bit)</option>
                        <option value="3">INT16(16-bit)</option>
                        <option value="4">INT8(8-bit)</option>
                        <option value="5">FLOAT64(64-bit)</option>
                        <option value="6">FLOAT32(32-bit)</option>
                        <option value="7">BINARY(byte array)</option>
                        <option value="8">Byte</option>
                        <option value="9">Unsigned INT64</option>
                        <option value="10">Unsigned INT32</option>
                        <option value="11">Unsigned INT16</option>
                        <option value="12">Unsigned INT8</option>
                    </select>
                    create field index
                        {{if $v}}
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
            <td><button onclick="add();">add field</button></td>
        </tr>
    </table>
    <hr>
    <button style="background-color: #7fbbff;width: 100px;height: 30px;font-size: large;"
        onclick="javascipt:if (confirm('confirm Alter？')){submit();};">submit</button>
    <script>
        function add() {
            var tr = document.createElement("tr");
            tr.innerHTML = '<td>field name：</td><td><span name="colums"><input type="text" placeholder="field name" value="" />'
                + ' <select name="fieldtype"><option value="0" selected>String</option><option value="1">INT64(64-bit)</option><option value="2">INT32(32-bit)</option><option value="3">INT16(16-bit)</option><option value="4">INT8(8-bit)</option><option value="5">FLOAT64(64-bit)</option><option value="6">FLOAT32(32-bit)</option><option value="7">BINARY(byte array)</option><option value="8">Byte</option><option value="9">Unsigned INT64</option><option value="10">Unsigned INT32</option><option value="11">Unsigned INT16</option><option value="12">Unsigned INT8</option></select>'
                + ' create field index <input type="checkbox" /></span></td><button onclick="del(this);">cancel</button>';
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

var deleteEnText = `<html>

<head>
    <title>tldb</title>
</head>

<body class="body">
    <h3>Delete table data</h3>
    <hr>
    <span><b>Query data by id</b></span>
    <form id="deleteform" action="/delete" method="post">
        <input name="type" value="1" hidden />
        table name<input name="tableName" placeholder="table name" value="{{ .TableName }}" />
        table ID<input name="tableId" placeholder="table ID"  value="" />
        <input type="submit" value="check out data" />
    </form>
    <hr>
    {{if ne .TableName ""}}
    <table border="0">
        <tbody id="ctable">
            <form id="deleteform2" action="/delete" method="post">
                <input name="type" value="2" hidden />
                <tr>
                    <th>table name：</th>
                    <td><input type="text" name="tableName" placeholder="表名" value="{{ .TableName }}" style="border: none;" /></td>
                </tr>
                <tr>
                    <th>ID：</th>
                    <td><input type="text" name="tableId" placeholder="ID" value="{{ .ID }}"  readonly style="border: none;"/></td>
                </tr>
                {{range $k,$v := .ColumnValue }}
                <tr>
                    <th>field name：</th>
                    <td><input type="text" placeholder="field name" value="{{ $k }}"
                                readonly style="border: none;" />
                        </td>
                    <td><textarea readonly>{{ $v }}</textarea></td>
                </tr>
                {{end}}
            </form>
        </tbody>
        <tr>
            <td></td>
            <td><button onclick="javascipt:if (confirm('confirm delete?')){document.getElementById('deleteform2').submit();};">Delete</button></td>
        </tr>
    </table>
    {{end}}

</body>

</html>`
var createEnText = `<html>

<head>
    <title>tldb</title>
</head>

<body class="body">
    <h3>Create table</h3>
    <hr>
    <table>
        <tbody id="ctable">
            <tr>
                <td>table name：</td>
                <td><input type="text" id="tablen" placeholder="table name" value="" /></td>
            </tr>
            <tr>
                <td>field name：</td>
                <td><span name="colums"><input type="text" placeholder="field name" value="" /> 
                    <select>
                        <option value="0" selected>String</option>
                        <option value="1">INT64(64-bit)</option>
                        <option value="2">INT32(32-bit)</option>
                        <option value="3">INT16(16-bit)</option>
                        <option value="4">INT8(8-bit)</option>
                        <option value="5">FLOAT64(64-bit)</option>
                        <option value="6">FLOAT32(32-bit)</option>
                        <option value="7">BINARY(byte array)</option>
                        <option value="8">Byte</option>
                        <option value="9">Unsigned INT64</option>
                        <option value="10">Unsigned INT32</option>
                        <option value="11">Unsigned INT16</option>
                        <option value="12">Unsigned INT8</option>
                    </select>
                    create filed index  <input type="checkbox" /></span></td>
            </tr>
        </tbody>

        <div id="createDiv">
        </div>
        <form id="createform" action="/create" method="post">
        </form>
        <tr></tr>
        <tr>
            <td></td>
            <td><button onclick="add();">add field</button></td>
        </tr>
    </table>

    <hr>
    <button style="background-color: #7fbbff;width: 100px;height: 30px;font-size: large;"
        onclick="javascipt:if (confirm('confirm create？')){submit();};">submit</button>
    <script>
        function add() {
            var tr = document.createElement("tr");
            tr.innerHTML = '<td>field name：</td><td><span name="colums"><input type="text" placeholder="field name" value="" />'
                + ' <select name="fieldtype"><option value="0" selected>String</option><option value="1">INT64(64-bit)</option><option value="2">INT32(32-bit)</option><option value="3">INT16(16-bit)</option><option value="4">INT8(8-bit)</option><option value="5">FLOAT64(64-bit)</option><option value="6">FLOAT32(32-bit)</option><option value="7">BINARY(byte array)</option><option value="8">Byte</option><option value="9">Unsigned INT64</option><option value="10">Unsigned INT32</option><option value="11">Unsigned INT16</option><option value="12">Unsigned INT8</option></select>'
                + ' create field index <input type="checkbox" /></span></td><button onclick="del(this);">cancel</button>';
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
var dataEnText = `<html>

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
        <h3 style="display: inline;">tldb OperationPlatform</h3>
    </span>
    &nbsp;&gt;&gt;
    <h4 class="important" style="display: inline;width:100%;">cluster state：</h4>
    {{if .Stat }}
    <span style="display:inline-block; background-color: aquamarine;width: 200px;">running</span>
    {{else}}
    <span style="display:inline-block; background-color: rgb(255, 0, 0); width: 200px;">initialize... &#9200;</span>
    {{end}}
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=zh">[中文]</a></h6>
    </span>
    <hr>
    <a href='/init'>users</a>
    <a href='/sysvar'>cluster env</a>
    <a href='/sys'>node params</a>
    <a href='/data' style="font-weight: bold;">data manipulation</a>
    <a href='/mq'>MQ DATA</a>
    <a href='/log'>sys log</a>
    <a href='/monitor'>monitor</a>
    <a href='/login'>login</a>
    <hr>
    <div>
        <div style="font-size: large; font-weight: bold;">data table structure</div>
        <div style="overflow:scroll;max-height: 300px;">
            <table border="1" class="important">
                <tr>
                    <th>table name</th>
                    <th>index field</th>
                    <th>field name</th>
                    <th>current id</th>
                    <th>export table data</th>
                </tr>
                {{range $k,$v := .Tb }}
                <tr>
                    <td>{{ $v.Name }}</td>
                    <td>{{ $v.Idxs }}</td>
                    <td>{{ $v.Columns }}</td>
                    <td>{{ $v.Seq }}</td>
                    <td><button onclick="javascipt:if (confirm('Export table data with caution because it may occupy a large amount of server memory. Confirm to export table data?')){exportdata(this);};">export data</button></td>
                </tr>
                {{end}}
            </table>
        </div>
        <hr>
        <button
            onclick="openPage('/create')">Create table</button>&nbsp;<button onclick="openPage('/alter')">Alter table</button>&nbsp;<button onclick="openPage('/drop')">Drop table</button>&nbsp;<button onclick="openPage('/insert')">Insert</button>&nbsp;<button onclick="openPage('/update')">Update</button>&nbsp;<button onclick="openPage('/delete')">Delete</button>
        <hr>
        <span><b>Query data by ID</b></span>
        <form id="dataform" action="/data" method="post">
            <input name="type" value="1" hidden />
            table name<input name="tableName" placeholder="table name" value="{{ .Sb.Name }}" />
            table ID<input name="tableId" placeholder="table ID" value="{{ .Sb.Id }}" />
            <input type="submit" value="Query" />
        </form>
        <form id="exportform" action="/export" method="post">
            <input name="exportName" id="exportName" value="" hidden>
        </form>
        <span><b>Query multiple data by ID</b></span>
        <form id="dataform" action="/data" method="post">
            <input name="type" value="3" hidden />
            table name<input name="tableName" placeholder="table name" value="{{ .Sb.Name }}" />
            table ID<input name="start" placeholder="start ID" value="{{ .Sb.StartId }}" />
            count<input name="limit" placeholder="number of queries" value="{{ .Sb.Limit }}" />
            <input type="submit" value="Query" />
        </form>

        <span><b>Query data by index</b></span>
        <form id="dataform" action="/data" method="post">
            <input name="type" value="2" hidden />
            table name<input name="tableName" placeholder="table name" value="{{ .Sb.Name }}" />
            field name<input name="cloName" placeholder="field name" value="{{ .Sb.ColumnName }}" />
            field value<input name="cloValue" placeholder="field value" value="{{ .Sb.ColumnValue }}" />
            start<input name="start" placeholder="start sequence number" value="{{ .Sb.StartId }}" />
            count<input name="limit" placeholder="number of queries" value="{{ .Sb.Limit }}" />
            <input type="submit" value="Query" />
        </form>
        <hr>
        {{if ne .Sb.Name ""}}
        <h4>result：</h4>
        <h6 class="text-danger">TABLENAME：{{ .Sb.Name }}</h6>
        {{end}}
        <table border="1" style="width: 100%;font-size: small;">
            {{range $k,$v := .Tds }}
                {{if eq $k 0}}
                <tr>
                <th>ID</th>
                {{range $k1,$v1 := $v.Columns }}
                <th>{{ $k1 }}</th>
                {{end}}
                </tr>
                {{end}}
            <tr>
                <td>{{ $v.Id }}</td>
                {{range $k1,$v1 := $v.Columns }}
                <td><textarea style="width: 100%;" readonly>{{ $v1 }}</textarea></td>
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
var initEnText = `<html>

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
        <h3 style="display: inline;">tldb OperationPlatform</h3>
    </span>
    &nbsp;&gt;&gt;
    <h4 class="important" style="display: inline;width:100%;">cluster state：</h4>
    {{if .Stat }}
    <span style="display:inline-block; background-color: aquamarine;width: 200px;">running</span>
    {{else}}
    <span style="display:inline-block; background-color: rgb(255, 0, 0); width: 200px;">initialize... &#9200;</span>
    {{end}}
    {{else if .Init}}
    <h3 style="display: inline;">tldb OperationPlatform </h3>
    {{end}}
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=zh">[中文]</a></h6>
    </span>
    <hr>
    <a href='/init' style="font-weight: bold;">users</a>
    <a href='/sysvar'>cluster env</a>
    <a href='/sys'>node params</a>
    <a href='/data'>data manipulation</a>
    <a href='/mq'>MQ DATA</a>
    <a href='/log'>sys log</a>
    <a href='/monitor'>monitor</a>
    <a href='/login'>login</a>
    {{if .ShowCreate }}
    <hr>
    <div>
        <h3>user account management</h3>
        <hr>
        <h4>create administrator <h5 class="important">{{ .Show }}</h5>
        </h4>
        <form id="createAdminform" action="/init?type=1" method="post">
            <input name="adminName" placeholder="username" />
            <input name="adminPwd" placeholder="password" type="password" />
            administrator<input name="adminType" type="radio" value="1" checked />&nbsp;&nbsp;
            {{if not .Init}}
            data administrator<input name="adminType" type="radio" value="2" />
            {{end}}
            <input type="submit" value="create administrator" />
        </form>

        {{if not .Init}}
        <hr>
        <h4>create MQ client</h4>
        <form id="createMQform" action="/init?type=1" method="post">
            <input name="mqName" placeholder="MQ username" />
            <input name="mqPwd" placeholder="password" type="password" />
            <input type="submit" value="create MQ client" />
        </form>
        <hr>
        <h4>create database client</h4>
        <form id="createCliform" action="/init?type=1" method="post">
            <input name="cliName" placeholder="client username" />
            <input name="cliPwd" placeholder="password" type="password" />
            <input type="submit" value="create database client" />
        </form>
        {{end}}
    </div>
    <hr>
    {{end}}
    {{if not .Init}}
    <hr>
    <div class="important" style="font-size: small;">
        <h4>administrator</h4>
        {{range $k,$v := .AdminUser}}
        <form id="adminform" action="/init?type=2" method="post">
            <input name="adminName" value='{{ $k }}' readonly /> access right:{{ $v }}
            <input type="button" value="delete user" onclick="javascipt:if (confirm('confirm delete?')){this.parentNode.submit();};" />
        </form>
        {{end}}
        <hr>
        <h4>MQ client</h4>
        {{range $k,$v := .MqUser }}
        <form id="mqform" action="/init?type=2" method="post">
            <input name="mqName" value="{{ $v }}" readonly />
            <input type="button" value="delete user" onclick="javascipt:if (confirm('confirm delete?')){this.parentNode.submit();};" />
        </form>
        {{end}}
        <hr>
        <h4>db client</h4>
        {{range $k,$v := .CliUser }}
        <form id="cliform" action="/init?type=2" method="post">
            <input name="cliName" value="{{ $v }}" readonly />
            <input type="button" value="delete user" onclick="javascipt:if (confirm('confirm delete?')){this.parentNode.submit();};" />
        </form>
        {{end}}
    </div>
    <hr>
    {{end}}

</html>`
var insertEnText = `<html>

<head>
    <title>tldb</title>
</head>

<body class="body">
    <h3>Insert table data</h3>
    <hr>
    <table>
        <tr>
            <form id="insertform" action="/insert" method="post">
                <input name="type" value="1" hidden />
                <td>table name：</td>
                <td><input type="text" name="tableName" placeholder="table name" value="" /></td>
            </form>
        </tr>
        <tr></tr>
        <tr>
            <td></td>
            <td><button onclick='javascript:document.getElementById("insertform").submit();'>Check out the table
                    structure</button></td>
        </tr>
    </table>
    <hr>
    {{if ne .TableName ""}}
    <table>
        <tbody id="ctable">
            <form id="insertform2" action="/insert" method="post">
                <input name="type" value="2" hidden />
                <tr>
                    <th>table name：</th>
                    <td><input type="text" name="tableName" placeholder="table name" value="{{ .TableName }}" /></td>
                </tr>
                {{range $k,$v := .Columns }}
                <tr>
                    <th>field name：</th>
                    <td><span><input type="text" name="colums" placeholder="field name" value="{{ $k }}" readonly />
                            <input type="text" name="values" placeholder="field value" value="" style="width: 400px;" />
                        </span>
                    </td>
                </tr>
                {{end}}
            </form>
        </tbody>
        <tr>
            <td></td>
            <td><button onclick="javascipt:if (confirm('confirm insert?')){document.getElementById('insertform2').submit();};">submit</button></td>
        </tr>
    </table>
    {{end}}
</body>

</html>`
var loadEnText = `<html>

<body style="text-align:center;">
    <h2>Importing data</h2>
    <div>
        <h3><span id='s'></span> pieces of data have been imported</h3>
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
            document.getElementById('e').innerHTML = "Data import complete";
        } else {
            document.getElementById('s').innerHTML = evt.data;
        }
    }
    ws.onclose = function (evt) {
        document.getElementById("e2").innerHTML = '<hr><h4>请<a href="javascript:window.history.go(-1)">click here</a>go back。<h4>'
    };
    ws.onopen = function (evt) {
    };
    ws.onerror = function (evt, e) {
    };
</script>

</html>`
var loginEnText = `<html>
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
    <h3 style="display: inline;">tldb OperationPlatform</h3>
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=zh">[中文]</a></h6>
    </span>
    <hr>
    <div id="login">
        <h3>login</h3>
        <form id="loginform" action="/login" method="post">
            <input name="type" value="1" hidden />
            <input name="name" placeholder="username" />
            <input name="pwd" placeholder="password" type="password" />
            <input type="submit" value="login" />
        </form>
    </div>
    <hr>
</html>`
var mqEnText = `<html>

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
        <h3 style="display: inline;">tldb OperationPlatform</h3>
    </span>
    &nbsp;&gt;&gt;
    <h4 class="important" style="display: inline;width:100%;">cluster state：</h4>
    {{if .Stat }}
    <span style="display:inline-block; background-color: aquamarine;width: 200px;">running</span>
    {{else}}
    <span style="display:inline-block; background-color: rgb(255, 0, 0); width: 200px;">initialize... &#9200;</span>
    {{end}}
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=zh">[中文]</a></h6>
    </span>
    <hr>
    <a href='/init'>users</a>
    <a href='/sysvar'>cluster env</a>
    <a href='/sys'>node params</a>
    <a href='/data'>data manipulation</a>
    <a href='/mq' style="font-weight: bold;">MQ DATA</a>
    <a href='/log'>sys log</a>
    <a href='/monitor'>monitor</a>
    <a href='/login'>login</a>
    <hr>
    <div>
        <div style="font-size: large; font-weight: bold;">data manipulation</div>
        <hr>
        <div style="overflow:scroll;max-height: 300px;">
            <table border="1" class="important" style="font-size: x-small;">
                <tr>
                    <th>pub topic(not contain MEM)</th>
                    <th>current ID</th>
                    <th>number of sub</th>
                    <th>Truncate</th>
                </tr>
                {{range $k,$v := .Tb }}
                <tr>
                    <td>{{ $v.Name }}</td>
                    <td>{{ $v.Seq }}</td>
                    <td>{{ $v.Sub }}</td>
                    <td>
                        <form action="/mq" method="post">
                            <input name="atype" value="1" hidden />
                            <input name="tableName" value="{{ $v.Name }}" hidden />
                            <input type="button" value="Truncate"
                                onclick="javascipt:if (confirm('confirm truncate?')){this.parentNode.submit();};" />
                        </form>
                    </td>
                </tr>
                {{end}}
            </table>
        </div>
        <hr>
        <span><b>Clears the specified range of MQ data</b></span>
        <form id="dataform" action="/mq" method="post">
            <input name="atype" value="2" hidden />
            <input class="btn btn-danger" type="button" value="delete"
                onclick="javascipt:if (confirm('sure delete?')){this.parentNode.submit();};" />
            <input name="tableName" placeholder="Topic" value="" />
            <input name="fromId" placeholder="input ID" value="" />
            <input name="limit" placeholder="input count" value="" />
        </form>
        <hr>
        <span><b>Query MQ data by Id</b></span>
        <form id="dataform" action="/mq" method="post">
            <input name="type" value="2" hidden />
            Topic<input name="tableName" placeholder="topic" value="{{ .Sb.Name }}" />
            ID<input name="tableId" placeholder="ID" value="{{ .Sb.Id }}" />
            <input type="submit" value="Query" />
        </form>
        <hr>
        <span><b>Query multiple MQ data by ID</b></span>
        <form id="dataform" action="/mq" method="post">
            <input name="type" value="3" hidden />
            Topic<input name="tableName" placeholder="topic" value="{{ .Sb.Name }}" />
            start ID<input name="start" placeholder="start id" value="{{ .Sb.StartId }}" />
            count<input name="limit" placeholder="number of queries" value="{{ .Sb.Limit }}" />
            <input type="submit" value="Query" />
        </form>
        <hr>
        {{if ne .Sb.Name ""}}
        <h4>result：</h4>
        <h6 class="text-danger">Topic：{{ .Sb.Name }}</h6>
        {{end}}
        <table border="1" style="width: 100%;font-size: small;">
            {{range $k,$v := .Tds }}
                {{if eq $k 0}}
                <tr>
                <th>ID</th>
                {{range $k1,$v1 := $v.Columns }}
                <th>{{ $k1 }}</th>
                {{end}}
                </tr>
                {{end}}
            <tr>
                <td>{{ $v.Id }}</td>
                {{range $k1,$v1 := $v.Columns }}
                <td><textarea style="width: 100%;" readonly>{{ $v1 }}</textarea></td>
                {{end}}
            </tr>
            {{end}}
        </table>
    </div>
</body>

</html>`
var sysEnText = `<html>

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
        <h3 style="display: inline;">tldb OperationPlatform</h3>
    </span>
    &nbsp;&gt;&gt;
    <h4 class="important" style="display: inline;width:100%;">cluster state：</h4>
    {{if .Stat }}
    <span style="display:inline-block; background-color: aquamarine;width: 200px;">running</span>
    {{else}}
    <span style="display:inline-block; background-color: rgb(255, 0, 0); width: 200px;">initialize... &#9200;</span>
    {{end}}
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=zh">[中文]</a></h6>
    </span>
    <hr>
    <a href='/init'>users</a>
    <a href='/sysvar'>cluster env</a>
    <a href='/sys' style="font-weight: bold;">node params</a>
    <a href='/data'>data manipulation</a>
    <a href='/mq'>MQ DATA</a>
    <a href='/log'>sys log</a>
    <a href='/monitor'>monitor</a>
    <a href='/login'>login</a>
    <hr>
    <div>
        <table border="1" style="font-size: 15px;">
            <tr>
                <th>name</th>
                <th>value</th>
                <th>startup parameters</th>
                <th>description</th>
            </tr>
            <tr>
                <td>Local data file(compressed file)：</td>
                <td class="important">{{ .SYS.DBFILEDIR }}</td>
                <td> -dir</td>
                <td>data file address</td>
            </tr>
            <tr>
                <td>BINLOG file size：</td>
                <td>{{ .SYS.BINLOGSIZE }}(MB)</td>
                <td> -binsize</td>
                <td>binlog data files are compressed every {{ .SYS.BINLOGSIZE }}M</td>
            </tr>
            <tr>
                <td>Whether MQ uses tls</td>
                <td>{{ .SYS.MQTLS }}</td>
                <td> -clitls</td>
                <td>access MQ using 'wss://'</td>
            </tr>
            <tr>
                <td>Whether web admin uses tls</td>
                <td>{{ .SYS.ADMINTLS }} </td>
                <td> -admintls</td>
                <td>access management using 'https://'</td>
            </tr>
            <tr>
                <td>Whether database client uses tls</td>
                <td>{{ .SYS.CLITLS }} </td>
                <td> -mqtls</td>
                <td>database client uses ssl to access the server</td>
            </tr>
            <tr>
                <td>Crt file address of tls on the client service</td>
                <td>{{ .SYS.CLICRT }}</td>
                <td> -clicrt</td>
                <td>crt certificate file address of the SSL certificate of the client service</td>
            </tr>
            <tr>
                <td>Key file address of tls on the client service</td>
                <td>{{ .SYS.CLIKEY }} </td>
                <td> -clikey</td>
                <td>key certificate file address of the SSL certificate of the client service</td>
            </tr>
            <tr>
                <td>Crt file address of tls on mq service</td>
                <td>{{ .SYS.MQCRT }} </td>
                <td> -mqcrt</td>
                <td>crt certificate file address of the SSL certificate of mq service</td>
            </tr>
            <tr>
                <td>Key file address of tls on mq service</td>
                <td>{{ .SYS.MQKEY }} </td>
                <td> -mqkey</td>
                <td>key certificate file address of the SSL certificate of mq service</td>
            </tr>
            <tr>
                <td>Crt file address of tls on management platform</td>
                <td>{{ .SYS.ADMINCRT }} </td>
                <td> -admincrt</td>
                <td>crt certificate file address of the SSL certificate of management platform</td>
            </tr>
            <tr>
                <td>Key file address of tls on management platform</td>
                <td>{{ .SYS.ADMINKEY }}</td>
                <td> -adminkey</td>
                <td>key certificate file address of the SSL certificate of management platform</td>
            </tr>
            <tr>
                <td>Number of cocurrent on insert/delete/update data</td>
                <td>{{ .SYS.COCURRENT_PUT }} </td>
                <td> -put</td>
                <td>number of cocurrent on insert/delete/update data.If more, wait in line</td>
            </tr>
            <tr>
                <td>Number of cocurrent on select data</td>
                <td>{{ .SYS.COCURRENT_GET }} </td>
                <td> -get</td>
                <td>number of cocurrent on select data.If more, wait in line</td>
            </tr>
            <tr>
                <td>Cluster namespace</td>
                <td>{{ .SYS.NAMESPACE }}</td>
                <td> -ns</td>
                <td>the namespaces of nodes in the cluster must be the same; otherwise, the nodes cannot be connected
                </td>
            </tr>
            <tr>
                <td>Password for connecting between nodes</td>
                <td>{{ .SYS.PWD }}</td>
                <td> -pwd</td>
                <td>password for connecting nodes in cluster</td>
            </tr>
            <tr>
                <td>Authenticated using TLS between nodes,Public key file address</td>
                <td>{{ .SYS.PUBLICKEY }}</td>
                <td> -publickey</td>
                <td>The public key in the tldb program is used by default.   
                    You can specify another public key address
                </td>
            </tr>
            <tr>
                <td>Authenticated using TLS between nodes,private key file address</td>
                <td>{{ .SYS.PRIVATEKEY }}</td>
                <td> -privatekey</td>
                <td>The private key in the tldb program is used by default. 
                    You can specify another private key address
                </td>
            </tr>
            <tr>
                <td>Cluster service address</td>
                <td>{{ .SYS.ADDR }}</td>
                <td> -cs</td>
                <td>address of the cluster service between nodes</td>
            </tr>
            <tr>
                <td>MQ service address</td>
                <td>{{ .SYS.MQADDR }}</td>
                <td> -mq</td>
                <td>MQ service address</td>
            </tr>
            <tr>
                <td>Client service address</td>
                <td>{{ .SYS.CLIADDR }}</td>
                <td> -cli</td>
                <td>client service address</td>
            </tr>
            <tr>
                <td>Management platform service address</td>
                <td>{{ .SYS.WEBADMINADDR }}</td>
                <td> -admin</td>
                <td>management platform service address.</td>
            </tr>
            <tr>
                <td>Minimum number of cluster nodes</td>
                <td>{{ .SYS.CLUSTER_NUM }}</td>
                <td> -clus</td>
                <td>automatic allocation by default;If the value is 0,single-node service</td>
            </tr>
            <tr>
                <td>Minimum number of cluster nodes fixed</td>
                <td>{{ .SYS.CLUSTER_NUM_FINAL }}(default false:means system allocation)</td>
                <td> -clus_final</td>
                <td>by default, the system automatically allocates the size.When the value is true, the -clus non-zero parameter value takes effect</td>
            </tr>
            <tr>
                <td>Program Version</td>
                <td>v{{ .SYS.VERSION }}</td>
                <td></td>
                <td>The development version of the current program</td>
            </tr>
        </table>
        <hr>
        <hr>

        <span style="font-size:large;font-weight: bold;">Import the compressed package data of bin.log[using data append]</span>
        <span style="font-size: xx-small;">the imported file is a compressed binlog file generated by tldb</span>
        <form id="loadForm1" action="/sys" method="post" enctype="multipart/form-data">
            <input name="atype" value="1" hidden />
            <input type="file" id="loadfile1" name="loadfile1" />
            <button
                onclick="javascipt:if (confirm('import data to this node may be cause inconsistent with data on other nodes。comfirm import?')){this.parentNode.submit();};">import data</button>
        </form>

        <hr>
        <span style="font-size:large;font-weight: bold;">Import the compressed package data of bin.log[using data coverage]</span>
        <span style="font-size: xx-small;">the imported file is a compressed binlog file generated by tldb</span>
        <form id="loadForm2" action="/sys" method="post" enctype="multipart/form-data">
            <input name="atype" value="2" hidden />
            <input type="file" id="loadfile2" name="loadfile2" />
            <button
                onclick="javascipt:if (confirm('import data to this node may be cause inconsistent with data on other nodes。comfirm import?')){this.parentNode.submit();};">import data</button>
        </form>
        <hr>
        <hr>
        <form id="sysForm" action="/sys" method="post" enctype="multipart/form-data">
            <input name="atype" value="3" hidden />
            <button
                onclick="javascipt:if (confirm('Close all services on the node？')){document.getElementById('sysForm').submit();};">close node service</button>
        </form>
    </div>
</body>

</html>`
var sysvarEnText = `<html>

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
        <h3 style="display: inline;">tldb OperationPlatform</h3>
    </span>
    &nbsp;&gt;&gt;
    <h4 class="important" style="display: inline;width:100%;">cluster state：</h4>
    {{if .Stat }}
    <span style="display:inline-block; background-color: aquamarine;width: 200px;">running</span>
    {{else}}
    <span style="display:inline-block; background-color: rgb(255, 0, 0); width: 200px;">initialize... &#9200;</span>
    {{end}}
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=zh">[中文]</a></h6>
    </span>
    <hr>
    <a href='/init'>users</a>
    <a href='/sysvar' style="font-weight: bold;">cluster env</a>
    <a href='/sys'>node params</a>
    <a href='/data'>data manipulation</a>
    <a href='/mq'>MQ DATA</a>
    <a href='/log'>sys log</a>
    <a href='/monitor'>monitor</a>
    <a href='/login'>login</a>
    <hr>
    <div>
        <table border="1" style="font-size:15px;">
            <tr>
                <td>Node startup time (local time)</td>
                <td colspan="2">{{ .SYS.StartTime }}</td>
            </tr>
            <tr>
                <td>Node server time</td>
                <td colspan="2">{{ .SYS.LocalTime }}</td>
            </tr>
            <tr>
                <td>Cluster system correction time</td>
                <td colspan="2">{{ .SYS.Time }}</td>
            </tr>
            <tr>
                <td>Node UUID</td>
                <td class="important" colspan="2">{{ .SYS.UUID }}</td>
            </tr>
            <tr>
                <td>Running Cluster Nodes</td>
                <td class="important" colspan="2">{{ .SYS.RUNUUIDS }}</td>
            </tr>
            <tr>
                <td class="important">Node state</td>
                <form id="statForm" action="/sysvar" method="post">
                    <input name="atype" value="3" hidden />
                </form>
                <form id="statForm2" action="/sysvar" method="post">
                    <input name="atype" value="4" hidden />
                </form>
                <td colspan="2">
                    {{if eq .SYS.STAT "0"}}
                    <span style="display:inline-block; background-color: rgb(255, 0, 0);width: 70px;">Ready &#9200;</span>
                    {{else if eq .SYS.STAT "1"}}
                    <span style="display:inline-block; background-color: rgb(255, 255, 0);width: 70px;">Proxy
                        &#128274;</span>
                    {{else if eq .SYS.STAT "2"}}
                    <span style="display:inline-block; background-color:aquamarine;width: 70px;">Running&#9989;</span>
                    {{end}}
                    <span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>
                    <button
                        onclick="javascipt:if (confirm('current operations on some clients may fail if the status is changed。Confirm reset to Proxy?(CAUTION)')){document.getElementById('statForm').submit();};">reset to Proxy</button>
                    <button
                        onclick="javascipt:if (confirm('current operations on some clients may fail if the status is changed。Confirm reset to Ready?(CAUTION)')){document.getElementById('statForm2').submit();};" />reset to Ready</button>
                    {{if eq .SYS.STAT "0"}}
                    <span style="font-size: small;">Currently,{{.SYS.SyncCount}} data is synchronized</span>
                    {{end}}
                </td>

            </tr>
            <tr>
                <td>Deviation between cluster time and local time</td>
                <form id="timeForm" action="/sysvar" method="post">
                    <input name="atype" value="2" hidden />
                    <td colspan="2">
                        <input name="time_deviation" value="{{ .SYS.TIME_DEVIATION }}"
                            style="border: none;width: 100px;" />(nanosecond)
                        <input type="button" value="Modified time deviation(CAUTION)"
                            onclick="javascipt:if (confirm('Changing the node time may cause errors in log data。Confirm modifie time deviation？')){document.getElementById('timeForm').submit();};" />
                    </td>
                </form>
            </tr>
            <tr>
                <td>Access address of the current node in the cluster service</td>
                <td colspan="2">{{ .SYS.ADDR }}</td>
            </tr>
            <tr>
                <td>Access address of the management platform</td>
                <td colspan="2">{{ .SYS.ADMINADDR }}</td>
            </tr>
            <tr>
                <td>MQ service address</td>
                <td colspan="2">{{ .SYS.MQADDR }}/mq</td>
            </tr>
            <tr>
                <td>Database client service address</td>
                <td colspan="2">{{ .SYS.CLIADDR }}</td>
            </tr>
            <tr>
                <td>Minimum number of nodes in cluster</td>
                <td colspan="2">
                    <div class="important">
                        {{if eq .SYS.CLUSTER_NUM "0"}}
                        stand-alone node
                        {{else}}
                        {{ .SYS.CLUSTER_NUM }}
                        {{end}}
                    </div>
                </td>
            </tr>
            <tr>
                <td>Number of storage nodes</td>
                <form id="storeNumForm" action="/sysvar" method="post">
                    <input name="atype" value="5" hidden />
                    <td colspan="2">
                        <input name="storeNum" value="{{ .SYS.STORENODENUM }}"
                            style="border: none;width: 50px;" />[default:0,means all nodes]
                        <input type="button" value="Number of storage nodes"
                            onclick="javascipt:if (confirm('current operations on some clients may fail if the number is changed(After the modification, data will be synchronized to other nodes)。Confirm modify？')){document.getElementById('storeNumForm').submit();};" />
                    </td>
                </form>
            </tr>
            <tr>
                <td>data statistics of cocurrent operations on insert/delete/update data</td>
                <td colspan="2">{{ .SYS.CCPUT }}</td>
            </tr>
            <tr>
                <td>data statistics of cocurrent operations on select data</td>
                <td colspan="2">{{ .SYS.CCGET }}</td>
            </tr>
            <tr>
                <td>data statistics of insert/delete/update data</td>
                <td colspan="2">{{ .SYS.COUNTPUT }}</td>
            </tr>
            <tr>
                <td>data statistics of select data</td>
                <td colspan="2">{{ .SYS.COUNTGET }}</td>
            </tr>
            {{range $k,$v := .RN }}
            <tr>
                <th rowspan="8">remode node: {{ $v.UUID }}</th>
            </tr>
            <tr>
                <td style="width: 100px;">UUID</td>
                <td>{{ $v.UUID }}</td>
            </tr>
            <tr>
                <td class="important">status</td>
                <td>{{ $v.StatDesc }}</td>
            </tr>
            <tr>
                <td>cluster service address</td>
                <td>{{ $v.Addr }}</td>
            </tr>
            <tr>
                <td>service IP address</td>
                <td>{{ $v.Host }}</td>
            </tr>
            <tr>
                <td>management platform service address</td>
                <td>{{ $v.AdminAddr }}</td>
            </tr>
            <tr>
                <td>MQ service address</td>
                <td>{{ $v.MQAddr }}/mq</td>
            </tr>
            <tr>
                <td>db client service address</td>
                <td>{{ $v.CliAddr }}</td>
            </tr>
            {{end}}
        </table>
        <hr>
        <h5 class="important">{{ .Show }}</h5>
        <h3>Cluster operation</h3>
        <form id="" action="/sysvar" method="post">
            <input name="atype" value="1" hidden />
            <table border="1">
                <tr>
                    <th>Add cluster nodes</th>
                </tr>
                <tr>
                    <th>Destination node address</th>
                    <td><input type="text" id="addr" name="addr" value="" placeholder="Destination node address" /></td>
                    <td><input type="submit" value="submit" /></td>
                </tr>
            </table>
        </form>
        <hr>
    </div>
</body>

</html>`
var dropEnText = `<html>

<head>
    <title>tldb</title>
</head>

<body class="body">
    <h3>Drop table</h3>
    <hr>
    <div style="overflow:scroll;max-height: 300px;">
        <table border="1" class="important">
            <tr>
                <th>table name</th>
                <th>index field</th>
                <th>field name</th>
                <th>current ID</th>
                <th>drop</th>
            </tr>
            {{range $k,$v := .Tb }}
            <tr>
                <td>{{ $v.Name }}</td>
                <td>{{ $v.Idxs }}</td>
                <td>{{ $v.Columns }}</td>
                <td>{{ $v.Seq }}</td>
                <td><button onclick="javascipt:if (confirm('confirm drop tables and table data?')){del(this);};">drop</button></td>
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
var updateEnText = `<html>

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
    <h3>Update table data</h3>
    <hr>
    <span><b>Query data by ID</b></span>
    <form id="updateform" action="/update" method="post">
        <input name="type" value="1" hidden />
        table name<input name="tableName" placeholder="table name" value="{{ .TableName }}" />
        table ID<input name="tableId" placeholder="table ID" value="" />
        <input type="submit" value="Check out data" />
    </form>
    <hr>
    {{if ne .TableName ""}}
    <table border="0">
        <tbody id="ctable">
            <form id="updateform2" action="/update" method="post">
                <input name="type" value="2" hidden />
                <tr>
                    <th>table name：</th>
                    <td><input type="text" name="tableName" placeholder="table name" value="{{ .TableName }}"
                            style="border: none;" /></td>
                </tr>
                <tr>
                    <th>ID：</th>
                    <td><input type="text" name="tableId" placeholder="ID" value="{{ .ID }}" readonly
                            style="border: none;" /></td>
                </tr>
                {{range $k,$v := .ColumnValue }}
                <tr>
                    <th>field name：</th>
                    <td><input type="text" name="colums" placeholder="field" value="{{ $k }}" readonly
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
                    onclick="javascipt:if (confirm('confirm update?')){document.getElementById('updateform2').submit();};">update</button>
            </td>
        </tr>
    </table>
    {{end}}

</body>

</html>`

var monitorEnText = `<html>

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
        <h3 style="display: inline;">tldb OperationPlatform</h3>
    </span>
    &nbsp;&gt;&gt;
    <h4 class="important" style="display: inline;width:100%;">cluster state：</h4>
    {{if . }}
    <span style="display:inline-block; background-color: aquamarine;width: 200px;">running</span>
    {{else}}
    <span style="display:inline-block; background-color: rgb(255, 0, 0); width: 200px;">initialize... &#9200;</span>
    {{end}}
    <span style="text-align:right">
        <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=zh">[中文]</a></h6>
    </span>
    <hr>
    <a href='/init'>users</a>
    <a href='/sysvar'>cluster env</a>
    <a href='/sys'>node params</a>
    <a href='/data'>data manipulation</a>
    <a href='/mq'>MQ DATA</a>
    <a href='/log'>sys log</a>
    <a href='/monitor'  style="font-weight: bold;">monitor</a>
    <a href='/login'>login</a>
    <hr>
    <div>
        <h3>System data monitoring</h3>
        <hr>
        Monitoring interval (unit: second) <input id="stime" placeholder="input time" value="3" />
        <button onclick="monitorLoad();">Start</button>
        <button onclick="stop();">Stop</button>
        <button onclick="clearData();">ClearData</button>
    </div>
    <hr>
    <div>
        <table border="1" style="font-size: smaller;">
            <tr>
                <th></th>
                <th>Alloc(MB)</th>
                <th>TotalAlloc(MB)</th>
                <th>NumGC</th>
                <th>NumGoroutine</th>
                <th>NumCPU</th>
                <th>ConcurrencyCUD</th>
                <th>NumCUD</th>
                <th>ConcurrencyR</th>
                <th>NumR</th>
                <th>RamUsage</th>
                <th>DiskFree(GB)</th>
                <th>CpuUsage</th>
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
