// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package tc

import (
	htmlTpl "html/template"
	textTpl "text/template"

	"github.com/donnie4w/tlnet"
)

type TXTYPE int
type LANG int

const (
	ZH LANG = 0
	EN LANG = 1
)

const (
	_ TXTYPE = iota
	LOGIN
	INIT
	SYSVAR
	DATA
	MQ
	SYS
	CREATE
	ALTER
	INSERT
	DELETE
	UPDATE
	DROP
	LOAD
	LOADDATA
	MONITOR
)

var mod = 1 //0debug，1release

func tplToHtml(lang LANG, flag TXTYPE, v any, hc *tlnet.HttpContext) {
	switch flag {
	case LOGIN:
		tpl(lang, "./tc/html/login.html", loginText, "./tc/html/loginEn.html", loginEnText, v, hc)
	case INIT:
		tpl(lang, "./tc/html/init.html", initText, "./tc/html/initEn.html", initEnText, v, hc)
	case SYSVAR:
		tpl(lang, "./tc/html/sysvar.html", sysvarText, "./tc/html/sysvarEn.html", sysvarEnText, v, hc)
	case DATA:
		tpl(lang, "./tc/html/data.html", dataText, "./tc/html/dataEn.html", dataEnText, v, hc)
	case MQ:
		tpl(lang, "./tc/html/mq.html", mqText, "./tc/html/mqEn.html", mqEnText, v, hc)
	case SYS:
		tpl(lang, "./tc/html/sys.html", sysText, "./tc/html/sysEn.html", sysEnText, v, hc)
	case CREATE:
		tpl(lang, "./tc/html/create.html", createText, "./tc/html/createEn.html", createEnText, v, hc)
	case ALTER:
		tpl(lang, "./tc/html/alter.html", alterText, "./tc/html/alterEn.html", alterEnText, v, hc)
	case INSERT:
		tpl(lang, "./tc/html/insert.html", insertText, "./tc/html/insertEn.html", insertEnText, v, hc)
	case DELETE:
		tpl(lang, "./tc/html/delete.html", deleteText, "./tc/html/deleteEn.html", deleteEnText, v, hc)
	case UPDATE:
		tpl(lang, "./tc/html/update.html", updateText, "./tc/html/updateEn.html", updateEnText, v, hc)
	case DROP:
		tpl(lang, "./tc/html/drop.html", dropText, "./tc/html/dropEn.html", dropEnText, v, hc)
	case LOAD:
		tpl(lang, "./tc/html/load.html", loadText, "./tc/html/loadEn.html", loadEnText, v, hc)
	case MONITOR:
		tpl(lang, "./tc/html/monitor.html", monitorText, "./tc/html/monitorEn.html", monitorEnText, v, hc)
	}
}

func tpl(lang LANG, tplZHPath, tplZHText, tplENPath, tplENText string, v any, hc *tlnet.HttpContext) {
	if lang == ZH {
		if mod == 0 {
			textTplByPath(tplZHPath, v, hc)
		} else if mod == 1 {
			textTplByText(tplZHText, v, hc)
		}
	} else if lang == EN {
		if mod == 0 {
			textTplByPath(tplENPath, v, hc)
		} else if mod == 1 {
			textTplByText(tplENText, v, hc)
		}
	}
}

func textTplByPath(path string, data any, hc *tlnet.HttpContext) {
	if tp, err := textTpl.ParseFiles(path); err == nil {
		tp.Execute(hc.Writer(), data)
	} else {
		logger.Error(err)
	}
}

func textTplByText(text string, data any, hc *tlnet.HttpContext) {

	tl := textTpl.New("tldb")
	if _, err := tl.Parse(text); err == nil {
		tl.Execute(hc.Writer(), data)
	} else {
		logger.Error(err)
	}
}

func htmlTplByPath(path string, data any, hc *tlnet.HttpContext) {
	if tp, err := htmlTpl.ParseFiles(path); err == nil {
		tp.Execute(hc.Writer(), data)
	} else {
		logger.Error(err)
	}
}
