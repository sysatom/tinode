package setting

import "github.com/tinode/chat/server/extra/types"

type Rule []Row

type Row struct {
	Key    string
	Type   types.FormFieldType
	Title  string
	Detail string
}
