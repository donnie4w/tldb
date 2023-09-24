// Code generated by Thrift Compiler (0.18.1). DO NOT EDIT.

package tlcli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"
	// thrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/donnie4w/gothrift/thrift"
	"strings"
	"regexp"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = errors.New
var _ = context.Background
var _ = time.Now
var _ = bytes.Equal
// (needed by validator.)
var _ = strings.Contains
var _ = regexp.MatchString

const ERR_UNDEFINED = 1300
const ERR_AUTH_NOPASS = 1301
const ERR_TIMEOUT = 506
const ERR_NO_MATCH_PARAM = 401
const ERR_TABLE_EXIST = 409
const ERR_DATA_NOEXIST = 410
const ERR_TABLE_NOEXIST = 411
const ERR_COLUMN_NOEXIST = 412
const ERR_IDX_NOEXIST = 413

func init() {
}

